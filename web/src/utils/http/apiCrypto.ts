/**
 * 与 pkg/apicrypt 一致：AES-256-CBC + PKCS7，密钥 SHA256(UTF8)；IV 前置后 Base64；信封 { c }
 */
import CryptoJS from 'crypto-js'

export const API_CRYPTO_HEADER = 'X-Shh-Encrypted'
export const API_CRYPTO_HEADER_VALUE = '1'

function getKeyRaw(): string {
  // .env 中若未对含 $ 的值加双引号，dotenv 可能把 $xxx 当变量展开，导致与后端 encryptKey 不一致
  return (import.meta.env.VITE_API_ENCRYPT_KEY as string | undefined)?.trim() ?? ''
}

export function isApiCryptoEnabled(): boolean {
  return getKeyRaw() !== ''
}

function deriveKey(passphrase: string) {
  return CryptoJS.SHA256(CryptoJS.enc.Utf8.parse(passphrase))
}

/** 明文 UTF-8 字符串 -> 与 Go EncryptToEnvelopeJSON 中 c 字段相同的 Base64 */
export function encryptUtf8ToEnvelopeB64(plain: string): string {
  const key = deriveKey(getKeyRaw())
  const iv = CryptoJS.lib.WordArray.random(128 / 8)
  const enc = CryptoJS.AES.encrypt(plain, key, {
    iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  })
  const combined = iv.concat(enc.ciphertext)
  return CryptoJS.enc.Base64.stringify(combined)
}

export function decryptEnvelopeB64ToUtf8(b64: string): string {
  const key = deriveKey(getKeyRaw())
  const combined = CryptoJS.enc.Base64.parse(b64)
  const iv = CryptoJS.lib.WordArray.create(combined.words.slice(0, 4), 16)
  const ciphertext = CryptoJS.lib.WordArray.create(combined.words.slice(4), combined.sigBytes - 16)
  const dec = CryptoJS.AES.decrypt({ ciphertext } as any, key, {
    iv,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7
  })
  return dec.toString(CryptoJS.enc.Utf8)
}

function headerGet(headers: any, name: string): string | undefined {
  if (!headers) return undefined
  if (typeof headers.get === 'function') {
    const v = headers.get(name) ?? headers.get(name.toLowerCase())
    if (v != null) return String(v)
  }
  const lower = name.toLowerCase()
  const direct = headers[name] ?? headers[lower]
  if (direct != null) return String(direct)
  return undefined
}

export function isEncryptedResponseHeader(headers: any): boolean {
  return headerGet(headers, API_CRYPTO_HEADER) === API_CRYPTO_HEADER_VALUE
}

/** 从响应体取出信封字段 c（支持 Axios 传入字符串或已解析的 { c } 对象） */
function extractEnvelopeC(data: unknown): string | undefined {
  if (typeof data === 'string') {
    try {
      const o = JSON.parse(data) as { c?: unknown }
      return typeof o?.c === 'string' ? o.c : undefined
    } catch {
      return undefined
    }
  }
  if (data !== null && typeof data === 'object' && 'c' in data) {
    const c = (data as { c?: unknown }).c
    return typeof c === 'string' ? c : undefined
  }
  return undefined
}

/** transformResponse 阶段：若为密文信封则解密为明文 JSON 字符串（供后续 JSON.parse 成业务体） */
export function maybeDecryptResponseBody(data: unknown, headers: any): unknown {
  if (!isApiCryptoEnabled()) return data
  if (!isEncryptedResponseHeader(headers)) return data

  const cVal = extractEnvelopeC(data)
  if (cVal === undefined) return data

  try {
    return decryptEnvelopeB64ToUtf8(cVal)
  } catch {
    return data
  }
}
