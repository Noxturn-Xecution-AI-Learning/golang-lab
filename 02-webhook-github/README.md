# 02 — Webhook GitHub

## Tujuan
Membuat HTTP server di Golang yang bisa menerima dan memproses event dari GitHub Webhooks secara real-time.

## Konsep Penting

| Istilah | Penjelasan |
|---|---|
| **Webhook** | GitHub kirim HTTP POST ke servermu saat ada event |
| **X-Hub-Signature-256** | Header HMAC untuk verifikasi request dari GitHub |
| **X-GitHub-Event** | Tipe event: `push`, `pull_request`, `issues`, dll |
| **ngrok** | Tool untuk expose localhost ke internet (buat dev) |
| **Secret** | Kunci rahasia untuk validasi request beneran dari GitHub |

## Cara Jalankan

### 1. Jalankan server
```bash
cd 02-webhook-github
go run main.go
# Server jalan di :8080
```

### 2. Expose ke internet dengan ngrok
```bash
# Install ngrok: https://ngrok.com
ngrok http 8080
# Copy URL https://xxxx.ngrok.io
```

### 3. Setup di GitHub
```
Repo → Settings → Webhooks → Add webhook
  Payload URL : https://xxxx.ngrok.io/webhook
  Content type: application/json
  Secret      : rahasia-lokal-123
  Events      : Just the push event
```

### 4. Test — push sesuatu ke repo, lihat log server!

## Environment Variables

```bash
GITHUB_WEBHOOK_SECRET=secret-kamu
PORT=8080
```
