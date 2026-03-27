// Package repository memuat perwujudan konkret (implementasi nyata)
// dari janji-janji / kontrak pada User Repository Interface di package `domain`.
package repository

import (
	"context"
	"database/sql"
	"errors"

	"yourapp/internal/modules/user/domain"
)

// SQLRepository adalah wujud nyata dari domain.Repository untuk database SQL relasional
type SQLRepository struct {
	writeDB *sql.DB // Jalur koneksi untuk operasi ubah data (Set/Update/Insert)
	readDB  *sql.DB // Jalur koneksi yang aman untuk baca skala besar (Select)
}

// NewSQLRepository menciptakan gudang data baru.
func NewSQLRepository(writeDB, readDB *sql.DB) *SQLRepository {
	return &SQLRepository{
		writeDB: writeDB,
		readDB:  readDB,
	}
}

// FindByID mencari user ke Replica DB (readDB).
// Karena operasi ini hanya membaca, Query dialokasikan ke koneksi baca.
func (r *SQLRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	const q = `
		SELECT id, name, email, status, created_at, updated_at
		FROM users
		WHERE id = ? LIMIT 1
	` // Parameter disisipkan via '?' guna menghindari serangan SQL Injection massal.

	var user domain.User
	// QueryRowContext aman digunakan dan digabungkan ke request timeout via `ctx`.
	err := r.readDB.QueryRowContext(ctx, q, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Jika datanya memang terbukti kosong/hapus (bukan crash/putus internet):
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Kembalikan nilai hampa tanpa eror 500
		}
		return nil, err // Lemparkan eror fatalnya ke Service
	}

	return &user, nil // Terlempar ke layer Service bersama representasi objek `User` yang sudah terisi.
}

// FindByEmail tugasnya mirip FindByID
func (r *SQLRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, name, email, status, created_at, updated_at
		FROM users
		WHERE email = ? LIMIT 1
	`

	var user domain.User
	err := r.readDB.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Create melempar perintah INSERT baru ke Database TULIS (writeDB).
func (r *SQLRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (name, email, status, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	` // Kita delegasikan tanggal ke NOW() milik SQL agar lebih presisi.

	res, err := r.writeDB.ExecContext(ctx, q, user.Name, user.Email, user.Status)
	if err != nil {
		return err
	}

	// Setelah tabel ditambahkan, SQL mencetak ID terbarunya (auto_increment).
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Tangkap ID yang disuntik dari DB kemudian kembalikan agar Service tau ID mana
	user.ID = id
	return nil
}

// Update merombak nama/status ke Database Master (writeDB).
func (r *SQLRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users
		SET name = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.writeDB.ExecContext(ctx, q, user.Name, user.Status, user.ID)
	return err // Tidak butuh LastInsertId karena rute pemanggilan datanya cuma sekadar tahu sukses UPDATE atau tidak.
}
