package db // Kumpulan blok kode terkait koneksi infrastruktur dasar database SQL

import (
	"database/sql" // Package standar Go untuk berkomunikasi dengan database relasional

	_ "github.com/go-sql-driver/mysql" // Mengimpor driver MySQL secara anonim (_) supaya mendaftarkan dirinya ke package database/sql
	"yourapp/internal/config"          // Mengambil definisi struktur pengaturan koneksi dari konfigurasi environment
)

// NewMySQL adalah fungsi pembuat awal (konstruktor) yang akan mendirikan koneksi tunggal ke satu MySQL URL (Read atau Write)
func NewMySQL(conn config.DatabaseConn) (*sql.DB, error) {
	// Membuka portal database (belum berarti terkoneksi secara harfiah, baru mempersiapkan memori)
	db, err := sql.Open("mysql", conn.DSN)
	if err != nil {
		return nil, err
	}

	// Mengatur batasan agar aplikasi tidak memberondong server MySQL dengan ratusan ribu koneksi
	db.SetMaxOpenConns(conn.MaxOpenConns)       // Kapasitas maksimal koneksi konkuren
	db.SetMaxIdleConns(conn.MaxIdleConns)       // Koneksi cadangan yang biarkan menganggur untuk request berikutnya
	db.SetConnMaxLifetime(conn.ConnMaxLifetime) // Usia maksimal koneksi agar tidak basi / terputus oleh firewall
	db.SetConnMaxIdleTime(conn.ConnMaxIdleTime) // Waktu toleransi koneksi nganggur (langsung dibunuh bila terlampaui)

	// Ping() bertujuan "mengetuk pintu" / mencoba menyambung ke MySQL betulan, memastikan password/user valid.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Jika berhasil, kembalikan objek koneksi database
	return db, nil
}
