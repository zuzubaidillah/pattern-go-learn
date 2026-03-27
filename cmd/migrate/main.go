package main // Entrypoint (Pelatuk peluncur eksekusi) khusus program Terminal Migrasi. Ini BERDIRI SENDIRI MURNI terpisah dari Router API Gin.

import (
	"flag"
	"fmt"
	"log"

	// Driver-driver spesifik milik golang-migrate wajib dipanggil (diimpor paksa init-nya via blank identifier "_") 
	// jika tidak, mesin migrasi "New Migrate()" tidak akan tahu cara menyedot file lokal dan koneksi SQL-nya.
	_ "github.com/go-sql-driver/mysql"                      // Driver pangkalan MYSQL orisinil Golang
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // Pendaftaran ekstensi dialek migrasi khusus MySQL  
	_ "github.com/golang-migrate/migrate/v4/source/file"    // Pendaftaran sinkronisasi pembacaan file lokal "file://"

	"github.com/golang-migrate/migrate/v4"
	"yourapp/internal/app/bootstrap"
	"yourapp/internal/infra/db"
)

func main() {
	// 1. Memanggil pendeteksi bendera (Flag) pemantik terminal dari eksekutor Linux/MacOS.
	// Contoh tata bahasa panggilannya: $ go run cmd/migrate/main.go -target=mysql_main -cmd=up
	targetName := flag.String("target", "mysql_main", "Tentukan basis data parsial mana yang mau di-migrate saat ini.")
	command := flag.String("cmd", "up", "Tentukan poros aksinya (up, down, drop)")
	flag.Parse()

	// 2. Impor struktur rahasia dasar sistem kita menggunakan alat pemanggil yang sama dengan API Utama.
	cfg := bootstrap.MustLoadConfig()

	// 3. Susun daftar inventaris target sinkronisasi schema yang "legal" / dilegitimasi sistem di fungsi BuildMigrationTargets kita.
	targets := db.BuildMigrationTargets(cfg)

	// 4. Cari spesifikasi target yang kebetulan senada dan jodoh dengan parameter Flag (-target) di atas.
	var finalTarget *db.MigrationTarget
	for _, t := range targets {
		if t.Name == *targetName {
			finalTarget = &t
			break // Langsung rem putaran loop, efisiensi!  
		}
	}

	// Kalau pengetik typo "-target=mysql_min" misalnya.
	if finalTarget == nil {
		log.Fatalf("❌ Tolak Akses: Target database '%s' tak ditemukan pendaftarannya di skema Konfigurasi Go Anda!", *targetName)
	}

	// 5. Instansiasi jembatan ajaib (The Core Engine) golang-migrate yang menjahit "Folder direktori File Source" vs "Tebing Dialek MySQL".
	// Bila di sini tersandung Error Koneksi Refused, artinya password / port docker database Anda tewas.
	m, err := migrate.New(finalTarget.Dir, finalTarget.DSN)
	if err != nil {
		log.Fatalf("❌ Gagal meracik pabrik alat pemigrasian (New Migrate Generator): %v", err)
	}

	fmt.Printf("🎯 Validasi Target Database Terkunci: [%s]\n", finalTarget.Name)
	fmt.Printf("📝 Validasi Perintah Migrasi Terdeteksi: [%s]\n", *command)

	// 6. Jalankan Penyalur Eksekutor Command (Switch / Case Routing murni gaya Console OS)
	switch *command {
	case "up":
		// 'Up' berarti mengeksekusi semua lembar arsip kode .up.sql yang terutang serentak, dan menaikkan level database schema secara konstan.
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("❌ Proses penyuntikan Migrasi UP Gagal Tersendat / Bentrok: %v", err)
		}
		
		// Apabila tak ada lembar migrasi baru (ErrNoChange), anggap itu kesuksesan wajar yang damai.
		if err == migrate.ErrNoChange {
			fmt.Println("✅ Santai! Basis data sudah berada di puncak suksesi versi teratas (Tak ada perombakan schema).")
		} else {
			fmt.Println("✅ WOW! Formasi serentak Migrasi Tabel DB-UP Anda Sukses Berjalan Sempurna!")
		}
		
	case "down":
		// 'Down' berarti membunuh versi schema yang sudah dirakit SATU KALI langkah mundur ke belakang. (Aturan standard safety -1 Step).
		// AWAS! Menggunakan tool ini sembarangan dapat mengakibatkan kelenyapan fatal seluruh isi tabel di ranah Production!
		if err := m.Steps(-1); err != nil {
			log.Fatalf("❌ Eksekusi Proses Migrasi DOWN Gagal Tersendat: %v", err)
		}
		fmt.Println("✅ Tarian Proses Mundur (Down/Rollback) tepat satu langkah file berhasil dihapus/drop bersih!")

	case "drop":
		// 'Drop' itu mode Nuklir Mutlak. Membumihanguskan mutlak tanpa belas kasihan. Hati-hati kepencet!
		if err := m.Drop(); err != nil {
			log.Fatalf("❌ Invasi Proses DROP Nuke Gagal Tersendat: %v", err)
		}
		fmt.Println("☠️ Operasi Drop Database Total sukses diterapkan bak Gurun Sahara.")

	default:
		// Saringan penangkal Hacker Script Kiddies yang iseng masukin nama aneh.
		log.Fatalf("❌ Komando '%s' tidak disokong! Cuma mengenal format standar lazim (up / down / drop)", *command)
	}
}
