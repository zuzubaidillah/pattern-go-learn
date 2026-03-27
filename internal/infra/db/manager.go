package db // Package db mewakili logika akses koneksi murni (hanya infra, bukan logika bisnis)

import "database/sql"

// Manager adalah struktur objek yang memegang seluruh kolam (pools) database Anda
// Di sistem skala besar, wajar memiliki dua tipe server database (Write Master & Read Replica)
// dan wajar juga memiliki server mesin dari merek berbeda (MySQL + Postgres) jika ada layanan spesial
type Manager struct {
	MySQLWrite    *sql.DB // Menampung driver asli SQL Pool untuk database Write-only
	MySQLRead     *sql.DB // Menampung driver asli SQL Pool untuk database Read-only
	PostgresWrite *sql.DB // (Opsional) jika kelak pakai postgres
	PostgresRead  *sql.DB // (Opsional) jika kelak pakai postgres
}

// PrimaryUserDB adalah fungsi pembantu (accessor) supaya setiap Service yang memanggil tidak bingung mana DB write
// Fungsi ini hanya menyodorkan objek instance DB Write (MySQL).
func (m *Manager) PrimaryUserDB() *sql.DB {
	return m.MySQLWrite
}

// ReadUserDB adalah fungsi pembantu (accessor) untuk mengembalikan koneksi yang khusus dipakai saat Query SELECT.
func (m *Manager) ReadUserDB() *sql.DB {
	return m.MySQLRead
}
