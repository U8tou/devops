package apicrypt

import (
	"bytes"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	key := "shh_encrypt_key_test"
	plain := []byte(`{"userName":"a","password":"b"}`)
	env, err := EncryptToEnvelopeJSON(plain, key)
	if err != nil {
		t.Fatal(err)
	}
	out, err := DecryptEnvelopeToPlain(env, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("got %q want %q", out, plain)
	}
}

// 与常见 setting.yaml 中带 $、! 的密钥一致，避免前后端约定偏差
func TestRoundTripSpecialCharsKey(t *testing.T) {
	key := "!jyyg$Pp"
	plain := []byte(`{"ids":["1","2"]}`)
	env, err := EncryptToEnvelopeJSON(plain, key)
	if err != nil {
		t.Fatal(err)
	}
	out, err := DecryptEnvelopeToPlain(env, key)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, plain) {
		t.Fatalf("got %q want %q", out, plain)
	}
}
