package main // Menandakan bahwa ini adalah paket utama (entry point) aplikasi

import (
	"log" // Digunakan untuk mencetak log sistem standar ke terminal

	"yourapp/internal/app/bootstrap" // Mengimpor package bootstrap yang mengatur konfigurasi dan injeksi dependensi
)

// main adalah fungsi pertama yang dieksekusi saat aplikasi Go dijalankan
func main() {
	// Memanggil fungsi NewApp() untuk mempersiapkan semua kebutuhan server (Database, Logger, Middleware)
	app, err := bootstrap.NewApp()

	// Mengecek apakah ada error selama proses setup aplikasi
	if err != nil {
		// Jika terjadi error, paksa program berhenti dengan status log.Fatal
		log.Fatal(err)
	}

	// app.Run() bertugas menyalakan HTTP Server agar mulai mendengarkan port (misalnya :8080)
	// Kita periksa kembali nilainya, apabila server berhenti dengan error (misal port bentrok), ia akan di-log.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
