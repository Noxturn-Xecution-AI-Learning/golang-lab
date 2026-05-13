# 01 — PASETO v4 Public

## Tujuan
Belajar membuat dan memverifikasi token aman menggunakan PASETO v4 Public (asymmetric signing dengan Ed25519).

## Konsep Penting

| Istilah | Penjelasan |
|---|---|
| PASETO | Alternatif JWT yang lebih aman — format dan algoritma tidak bisa diubah |
| v4.public | Signing asymmetric dengan Ed25519 — payload terbaca, tapi tidak bisa dipalsukan |
| Secret Key | Private key untuk sign token — jangan pernah dibagikan atau di-commit ke repo |
| Public Key | Untuk verify token — aman dibagikan ke siapapun |
| Sign | Menandatangani token dengan secret key menghasilkan string `v4.public.<payload>.<signature>` |
| Verify | Memastikan signature cocok dengan payload menggunakan public key |
| Pointer `*paseto.Token` | Library mengembalikan referensi ke token, bukan salinan nilai langsung |

## v4.public vs v4.local

- `v4.public` — payload bisa dibaca siapapun (Base64), tapi tidak bisa dipalsukan
- `v4.local` — payload terenkripsi, hanya bisa dibaca yang punya symmetric key

Pakai `v4.public` untuk auth token antar service. Pakai `v4.local` kalau payload harus tersembunyi.

## Struktur Token

```
v4.public . eyJleHAiOi...N9 . hf2odpa...2CQ
    |              |                |
  versi+tipe    payload          signature
               (Base64)          (Ed25519)
```

Payload bisa didekode di https://token.dev — tapi tanpa secret key, signature tidak bisa dibuat ulang.

## Alur Kerja di Aplikasi Nyata

```
1. Server generate keypair sekali → simpan secret key di env variable
2. User login → server buat token berisi user_id, role, dll → sign → kirim ke client
3. Client simpan token (header Authorization)
4. Client kirim request → sertakan token
5. Server verify token pakai public key → kalau valid, proses request
```

## Cara Jalankan

```bash
cd 01-paseto
go run main.go
```

## Dependency

```bash
go get aidanwoods.dev/go-paseto
```

## Output Aktual

```
--- PASETO v4 Public: Sign & Verify ---

[1] Generate keypair
    Secret Key (hex): 7cf52046...dfddf
    Public Key (hex): 9789192b...dfddf

[2] Membuat token
    Token dibuat dengan data: user_id=u-001, role=admin, nama=Dwi Golang

[3] Sign token
    Token string:
    v4.public.eyJleHAiOi...2CQ

[4] Verifikasi token
    Token valid!
    user_id : u-001
    role    : admin
    nama    : Dwi Golang

[5] Test token yang dimanipulasi
    Token palsu berhasil ditolak
    Alasan: bad signature
```

## Catatan Belajar

- Public Key adalah separuh terakhir dari Secret Key — keduanya diturunkan dari algoritma Ed25519
- Token yang dimanipulasi sekecil apapun langsung ditolak dengan error `bad signature`
- Return type `*paseto.Token` (pointer) karena library Go mengembalikan referensi, bukan nilai langsung
- Di production: secret key disimpan di environment variable, bukan di-hardcode di kode