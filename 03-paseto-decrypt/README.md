# 03 — PASETO v4 Decrypt / Verify dari Ciphertext

## Tujuan
Menerima token PASETO v4 dari luar (misal dari API lain) dan melakukan verifikasi + decode payload-nya.

## Status
⬜ Belum dikerjakan — jadwal: Sprint berikutnya

## Preview
```go
// Skenario: kamu hanya punya public key + token string
// → verifikasi & ambil datanya

publicKeyHex := os.Getenv("PASETO_PUBLIC_KEY")
tokenStr := r.Header.Get("Authorization")

publicKey, _ := paseto.NewV4AsymmetricPublicKeyFromHex(publicKeyHex)
parser := paseto.NewParser()
parser.AddRule(paseto.NotExpired())

token, err := parser.ParseV4Public(publicKey, tokenStr, nil)
```
