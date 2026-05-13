package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// WebhookSecret harus sama dengan yang diset di GitHub
// Di production, pakai environment variable!
var webhookSecret = getEnv("GITHUB_WEBHOOK_SECRET", "rahasia-lokal-123")

// ---------------------------------------------------
// Struct payload GitHub (push event)
// ---------------------------------------------------
type PushPayload struct {
	Ref        string `json:"ref"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
		} `json:"author"`
	} `json:"commits"`
}

func main() {
	port := getEnv("PORT", "8080")

	http.HandleFunc("/webhook", handleWebhook)
	http.HandleFunc("/health", handleHealth)

	fmt.Printf("🚀 Server jalan di http://localhost:%s\n", port)
	fmt.Println("📡 Endpoint webhook: POST /webhook")
	fmt.Println("💡 Tip: Pakai ngrok untuk expose ke internet:")
	fmt.Printf("   ngrok http %s\n\n", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// ---------------------------------------------------
// Handler utama webhook
// ---------------------------------------------------
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Baca body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Gagal baca body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Verifikasi signature dari GitHub
	signature := r.Header.Get("X-Hub-Signature-256")
	if !verifySignature(body, signature) {
		log.Println("❌ Signature tidak valid!")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ambil tipe event
	eventType := r.Header.Get("X-GitHub-Event")
	log.Printf("📨 Event diterima: %s\n", eventType)

	// Handle berdasarkan tipe event
	switch eventType {
	case "push":
		handlePush(body)
	case "ping":
		fmt.Println("🏓 Ping dari GitHub! Webhook berhasil terhubung.")
	default:
		fmt.Printf("ℹ️  Event '%s' diterima tapi belum di-handle\n", eventType)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ---------------------------------------------------
// Handle event push
// ---------------------------------------------------
func handlePush(body []byte) {
	var payload PushPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("❌ Gagal parse payload:", err)
		return
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	fmt.Printf("\n🔔 Push Event!\n")
	fmt.Printf("   Repo   : %s\n", payload.Repository.FullName)
	fmt.Printf("   Branch : %s\n", branch)
	fmt.Printf("   Pusher : %s\n", payload.Pusher.Name)
	fmt.Printf("   Commits: %d\n", len(payload.Commits))

	for i, c := range payload.Commits {
		fmt.Printf("   [%d] %s — %s (%s)\n", i+1, c.ID[:7], c.Message, c.Author.Name)
	}
	fmt.Println()
}

// ---------------------------------------------------
// Verifikasi HMAC-SHA256 signature dari GitHub
// ---------------------------------------------------
func verifySignature(body []byte, signature string) bool {
	if signature == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(webhookSecret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
