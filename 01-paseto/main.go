package main

import (
	"fmt"
	"time"

	paseto "aidanwoods.dev/go-paseto"
)

// Keypair menyimpan sepasang kunci untuk signing dan verifikasi.
// SecretKey digunakan untuk sign (jangan pernah dibagikan).
// PublicKey digunakan untuk verify (aman dibagikan ke siapapun).
type Keypair struct {
	SecretKey paseto.V4AsymmetricSecretKey
	PublicKey paseto.V4AsymmetricPublicKey
}

// generateKeypair membuat pasangan kunci baru secara acak.
// Setiap kali dipanggil, hasilnya selalu berbeda.
// Di production: generate sekali, simpan SecretKey di environment variable.
func generateKeypair() Keypair {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	return Keypair{
		SecretKey: secretKey,
		PublicKey: secretKey.Public(),
	}
}

// createToken membuat token baru berisi data pengguna.
// Token ini belum ditandatangani — hanya berisi payload.
func createToken(userID, role, nama string) paseto.Token {
	token := paseto.NewToken()

	// Claim standar PASETO — digunakan untuk validasi waktu
	token.SetIssuedAt(time.Now())                      // kapan token dibuat
	token.SetNotBefore(time.Now())                     // token belum boleh dipakai sebelum waktu ini
	token.SetExpiration(time.Now().Add(1 * time.Hour)) // token kadaluarsa dalam 1 jam

	// Claim custom — data yang kamu simpan di dalam token
	token.Set("user_id", userID)
	token.Set("role", role)
	token.Set("nama", nama)

	return token
}

// signToken menandatangani token menggunakan secret key.
// Hasilnya adalah string token yang bisa dikirim ke client.
// Format: v4.public.<payload>.<signature>
func signToken(token paseto.Token, kp Keypair) string {
	// Parameter nil terakhir adalah "implicit assertion" — bisa diabaikan untuk sekarang
	return token.V4Sign(kp.SecretKey, nil)
}

// verifyToken memverifikasi dan mendekode token string.
// Menggunakan public key — tidak perlu secret key.
// Mengembalikan pointer ke token yang sudah didekode, atau error jika tidak valid.
// *paseto.Token (pointer) karena library mengembalikan referensi, bukan salinan nilai.
func verifyToken(tokenString string, kp Keypair) (*paseto.Token, error) {
	parser := paseto.NewParser()

	// Tambahkan aturan validasi: tolak token yang sudah kadaluarsa
	parser.AddRule(paseto.NotExpired())

	// ParseV4Public akan:
	// 1. Verifikasi signature menggunakan public key
	// 2. Cek semua aturan yang ditambahkan (NotExpired, dll)
	// 3. Decode payload dan kembalikan sebagai Token
	return parser.ParseV4Public(kp.PublicKey, tokenString, nil)
}

func main() {
	fmt.Println("--- PASETO v4 Public: Sign & Verify ---")
	fmt.Println()

	// LANGKAH 1: Generate keypair
	// Anggap ini seperti membuat kunci + gembok.
	// SecretKey = kunci (hanya kamu yang pegang)
	// PublicKey = gembok (bisa dibagikan ke siapapun untuk verifikasi)
	fmt.Println("[1] Generate keypair")
	kp := generateKeypair()
	fmt.Println("    Secret Key (hex):", kp.SecretKey.ExportHex())
	fmt.Println("    Public Key (hex):", kp.PublicKey.ExportHex())
	fmt.Println()

	// LANGKAH 2: Buat token berisi data pengguna
	// Token belum aman di tahap ini — belum ditandatangani
	fmt.Println("[2] Membuat token")
	token := createToken("u-001", "admin", "Dwi Golang")
	fmt.Println("    Token dibuat dengan data: user_id=u-001, role=admin, nama=Dwi Golang")
	fmt.Println()

	// LANGKAH 3: Sign token dengan secret key
	// Setelah di-sign, token tidak bisa diubah tanpa terdeteksi
	fmt.Println("[3] Sign token")
	signed := signToken(token, kp)
	fmt.Println("    Token string:")
	fmt.Println("   ", signed)
	fmt.Println()

	// LANGKAH 4: Verifikasi token
	// Simulasi: server menerima token dari client dan memverifikasinya
	fmt.Println("[4] Verifikasi token")
	parsed, err := verifyToken(signed, kp)
	if err != nil {
		fmt.Println("    GAGAL:", err)
		return
	}

	var userID, role, nama string
	parsed.Get("user_id", &userID)
	parsed.Get("role", &role)
	parsed.Get("nama", &nama)

	fmt.Println("    Token valid!")
	fmt.Println("    user_id :", userID)
	fmt.Println("    role    :", role)
	fmt.Println("    nama    :", nama)
	fmt.Println()

	// LANGKAH 5: Coba manipulasi token — harus ditolak
	// Ini membuktikan bahwa token tidak bisa dipalsukan
	fmt.Println("[5] Test token yang dimanipulasi")
	fakeToken := signed + "diubah"
	_, err = verifyToken(fakeToken, kp)
	if err != nil {
		fmt.Println("    Token palsu berhasil ditolak")
		fmt.Println("    Alasan:", err)
	}
}
