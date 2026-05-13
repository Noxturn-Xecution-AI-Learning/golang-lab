package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	fmt.Println("🖱️  Menggerakkan kursor...")

	// Ambil posisi saat ini
	x, y := robotgo.Location()
	fmt.Printf("Posisi awal: (%d, %d)\n", x, y)

	// Gerak ke kanan perlahan
	for i := 0; i < 5; i++ {
		robotgo.Move(x+(i*50), y)
		time.Sleep(300 * time.Millisecond)
	}

	// Kembali ke posisi awal
	robotgo.Move(x, y)
	fmt.Println("Kursor kembali ke posisi awal")
}
