# 01 — PASETO v4 Public

## Tujuan
Belajar membuat dan memverifikasi token aman menggunakan PASETO v4 Public (asymmetric signing dengan Ed25519).

## Konsep Penting

| Istilah | Penjelasan |
|---|---|
| **PASETO** | Alternatif JWT yang lebih aman & opinionated |
| **v4.public** | Signing asymmetric dengan Ed25519 |
| **Secret Key** | Private key untuk *sign* token (jangan pernah expose!) |
| **Public Key** | Untuk *verify* token, aman dibagikan |
| **Sign** | Membuat token + tanda tangan kriptografis |
| **Verify** | Memastikan token belum diubah & masih valid |

## v4.public vs v4.local

- `v4.public` → payload **bisa dibaca** siapapun, tapi tidak bisa dipalsukan
- `v4.local` → payload **terenkripsi**, hanya bisa dibaca yang punya key
