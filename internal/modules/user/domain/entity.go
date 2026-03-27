// Package domain menyimpan definisi inti dari identitas bisnis
// Tanpa bergantung ke logic database spesifik (MySQL/Postgres) atau tipe permintaan/transport (HTTP)
package domain

import "time"

// User adalah entitas utama (Model) yang mencerminkan siapa pengguna ini di dalam alam pikir bisnis (Business Rule)
type User struct {
	ID        int64     // Angka unik primer (Primary key dari DB dipetakan kemari)
	Name      string    // Nama lengkap pendaftar
	Email     string    // Email unik pendaftar
	Status    string    // Status seperti "active", "inactive", "suspended"
	CreatedAt time.Time // Tanda waktu kapan User dibuat
	UpdatedAt time.Time // Tanda waktu perubahan data User (dimutakhirkan)
}
