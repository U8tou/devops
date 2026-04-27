// Package apicrypt 前后端约定 AES-256-CBC + SHA256 密钥、信封 {"c":base64(iv||cipher)}
package apicrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// HeaderName 请求/响应携带此头时表示 body 使用信封加密
const HeaderName = "X-Shh-Encrypted"

// HeaderValue 头取值，为 "1" 时启用加解密流程
const HeaderValue = "1"

// LocalsEncResponse Fiber Locals 键：本响应需加密输出
const LocalsEncResponse = "shh_enc_response"

type envelope struct {
	C string `json:"c"`
}

// DeriveKey SHA256(UTF8(passphrase)) 作为 AES-256 密钥
func DeriveKey(passphrase string) []byte {
	s := sha256.Sum256([]byte(passphrase))
	return s[:]
}

// EncryptToEnvelopeJSON 明文 UTF-8 JSON 字节 -> {"c":"base64(iv||ciphertext)"} JSON
func EncryptToEnvelopeJSON(plain []byte, passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, errors.New("apicrypt: empty passphrase")
	}
	key := DeriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padLen := aes.BlockSize - (len(plain) % aes.BlockSize)
	padded := make([]byte, len(plain)+padLen)
	copy(padded, plain)
	for i := len(plain); i < len(padded); i++ {
		padded[i] = byte(padLen)
	}
	ciphertext := make([]byte, len(padded))
	mode.CryptBlocks(ciphertext, padded)
	combined := append(iv, ciphertext...)
	b64 := base64.StdEncoding.EncodeToString(combined)
	return json.Marshal(envelope{C: b64})
}

// DecryptEnvelopeToPlain 解析 body 为 {"c":"..."}，解密得到原始 UTF-8 字节（一般为 JSON）
func DecryptEnvelopeToPlain(body []byte, passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, errors.New("apicrypt: empty passphrase")
	}
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, err
	}
	if env.C == "" {
		return nil, errors.New("apicrypt: missing c")
	}
	raw, err := base64.StdEncoding.DecodeString(env.C)
	if err != nil {
		return nil, err
	}
	if len(raw) < aes.BlockSize || len(raw)%aes.BlockSize != 0 {
		return nil, errors.New("apicrypt: invalid payload length")
	}
	iv := raw[:aes.BlockSize]
	ciphertext := raw[aes.BlockSize:]
	key := DeriveKey(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)
	if len(plain) == 0 {
		return nil, errors.New("apicrypt: empty plaintext")
	}
	padLen := int(plain[len(plain)-1])
	if padLen <= 0 || padLen > aes.BlockSize || padLen > len(plain) {
		return nil, errors.New("apicrypt: invalid padding")
	}
	for i := len(plain) - padLen; i < len(plain); i++ {
		if plain[i] != byte(padLen) {
			return nil, errors.New("apicrypt: invalid padding")
		}
	}
	return plain[:len(plain)-padLen], nil
}
