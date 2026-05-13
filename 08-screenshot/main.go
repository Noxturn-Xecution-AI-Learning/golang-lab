package main

import (
	"fmt"
	"time"

	"github.com/kbinani/screenshot"
	"image/png"
	"os"
)

func main() {
	fmt.Println("Mengambil screenshot layar...")

	n := screenshot.NumActiveDisplays()
	fmt.Printf("Jumlah monitor: %d\n", n)

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Printf("Gagal capture display %d: %v\n", i, err)
			continue
		}

		filename := fmt.Sprintf("screenshot_%d_%s.png", i, time.Now().Format("20060102_150405"))
		file, err := os.Create(filename)
		if err != nil {
			fmt.Printf("Gagal buat file: %v\n", err)
			continue
		}
		defer file.Close()

		png.Encode(file, img)
		fmt.Printf("Screenshot display %d disimpan: %s (%dx%d)\n",
			i, filename, bounds.Dx(), bounds.Dy())
	}
}
