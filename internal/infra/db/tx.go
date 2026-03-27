package db // Folder logika database fundamental tingkat rendah

import (
	"context"
	"database/sql"
)

// WithTransaction adalah sebuah "Fungsi Pembungkus" (Helper) yang memastikan kumpulan
// perintah query yang dioper ke dalamnya berjalan dalam satu siklus Transaksi DB Atomic (Rollback atau Commit).
//
// Cara kerjanya:
//  1. Ia menerima object DB Write dan sebuah fungsi lokal (fn) berisi instruksi SQL
//  2. Ia membuka jalur transaksi (`BeginTx`)
func WithTransaction(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	// Memulai blok transaksi
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err // Gagal memulai transaksi, kembalikan error
	}

	// Blok `defer` akan ditaruh di antrean dan selalu dijalankan saat fungsi (WithTransaction) mau selesai.
	// Jika tiba-tiba panic / berhenti, Rollback() akan dipanggil otomatis agar data setengah jalan tidak tersimpan.
	defer func() {
		_ = tx.Rollback() // Mengabaikan error lanjutan apabila transaksi sudah di-commit duluan
	}()

	// Menjalankan fungsi instruksi lokal (bisnis logic) yang dipassing dari luar (parameter `fn`)
	// Jika operasinya gagal/mendapat error, maka kita lempar (return) errornya dan `defer` di atas bertugas nge-rollback!
	if err := fn(tx); err != nil {
		return err
	}

	// Jika operasinya sukses tanpa error, lakukan tindakan Commit (simpan mutlak ke dalam Hard Disk DB).
	// Jika commit sukses, status selesai.
	return tx.Commit()
}
