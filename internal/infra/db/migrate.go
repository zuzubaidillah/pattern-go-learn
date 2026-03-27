package db // Merupakan kepanjangan tangan (ekstensi) infra logika yang bersinggungan erat ke basis data.

import (
	"fmt"
	"strings"

	"yourapp/internal/config"
)

// MigrationTarget mendeskripsikan secara konkret "siapa yang akan dimutakhirkan struktur tabelnya".
// Sangat esensil untuk mengisolasi agar satu file direktori sql tidak sembrono berjalan masuk ke database yang salah fungsi.
type MigrationTarget struct {
	Name string // Penamaan rujukan unik modul target, contoh: "mysql_main", "postgres_reporting"
	DSN  string // "Data Source Name", tali penyambung khusus otorisasi driver migration
	Dir  string // Lokasi Direktori (di Hard Disk) di mana file .sql spesifik grup ini kumpul melantai
	Kind string // Jenis mesin Database (Misal: mysql, postgres, sqlite) agar sintaksis Dialect Go tak salah melenceng 
}

// FormatDSN memastikan DSN kompatibel dengan library golang-migrate yang mensyaratkan keras penggunaan prefix skema protokol jaringan.
func FormatDSN(kind, rawDSN string) string {
	if strings.HasPrefix(rawDSN, kind+"://") {
		return rawDSN
	}
	// Contoh kasus kita merangkai: `user:pass@tcp(localhost:3306)/dbname` direndam menjadi `mysql://user:pass@tcp(localhost:3306)/dbname`
	return fmt.Sprintf("%s://%s", kind, rawDSN)
}

// BuildMigrationTargets menyusun dan membongkar rincian konfigurasi YAML / ENV (config.Config)
// menjadikannya deretan instansi MigrationTarget mentah yang tinggal dicaplok oleh program golang CLI.
func BuildMigrationTargets(cfg config.Config) []MigrationTarget {
	var targets []MigrationTarget

	// 1. Target Tabel Utama Aplikasi (Write Primary MySQL) 
	// DI SINI KUNCI ARSITEKTURNYA: Kita hanya memigrasikan database Authoritative/Write. 
	// Database Secondary / Read-Replica TIDAK dibolehkan dimigrasi! (Mereka cuma menyalin data secara logis dari Write).
	if cfg.DB.MySQL.Write.DSN != "" {
		targets = append(targets, MigrationTarget{
			Name: "mysql_main",                               // Label ID pemanggilan bagi Command CLI terminal (-target=mysql_main)
			DSN:  FormatDSN("mysql", cfg.DB.MySQL.Write.DSN), // Memakai identitas rahasia DB yang persis sama dengan perasan Server API (Diolah + Prefix mysql://)
			Dir:  "file://migrations/mysql_main",             // Wajib mutlak memakai prefix protocol 'file://' karena aturan keras golang-migrate driver
			Kind: "mysql",                                    // Petunjuk penyesuaian (adaptation pattern) MySQL golang runtime.
		})
	}

	// Andaikata kedepannya Startup anda berekspansi menambah PostgreSQL analytis, target array di atas bisa ditambahkan. 
	// Sehingga aplikasi monoloth/microservice anda tetap solid multi-migrasi tanpa campur aduk file.

	return targets
}
