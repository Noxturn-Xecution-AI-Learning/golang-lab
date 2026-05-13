package main

import (
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"
)

func main() {
	fmt.Println("=== PASETO v4 Public: Sign & Verify ===\n")

	// ---------------------------------------------------
	// STEP 1: Generate keypair
	// Di production, private key disimpan di env variable!
	// ---------------------------------------------------
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public()

	fmt.Println("✅ Keypair berhasil di-generate")
	fmt.Println("🔑 Public Key (hex):", publicKey.ExportHex())
	fmt.Println()

	// ---------------------------------------------------
	// STEP 2: Buat & Sign token
	// ---------------------------------------------------
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(1 * time.Hour))

	// Set payload custom
	token.Set("user_id", "u-001")
	token.Set("role", "admin")
	token.Set("nama", "Budi Golang")

	// Sign → hasilkan token string
	signed := token.V4Sign(secretKey, nil)
	fmt.Println("✅ Token berhasil di-sign:")
	fmt.Println(signed)
	fmt.Println()

	// ---------------------------------------------------
	// STEP 3: Verify & Decode token
	// ---------------------------------------------------
	fmt.Println("=== Verifikasi Token ===\n")

	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsed, err := parser.ParseV4Public(publicKey, signed, nil)
	if err != nil {
		fmt.Println("❌ Token tidak valid:", err)
		return
	}

	// Ambil data dari token
	var userID, role, nama string
	parsed.Get("user_id", &userID)
	parsed.Get("role", &role)
	parsed.Get("nama", &nama)

	fmt.Println("✅ Token valid!")
	fmt.Println("👤 User ID :", userID)
	fmt.Println("🔐 Role    :", role)
	fmt.Println("📛 Nama    :", nama)

	// ---------------------------------------------------
	// STEP 4: Coba verifikasi token yang diubah (harus gagal)
	// ---------------------------------------------------
	fmt.Println()
	fmt.Println("=== Test Token Palsu ===\n")
	fakeToken := signed + "tampered"
	_, err = parser.ParseV4Public(publicKey, fakeToken, nil)
	if err != nil {
		fmt.Println("✅ Token palsu berhasil ditolak:", err)
	}
}
