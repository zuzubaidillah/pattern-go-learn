package repository // Menjembatani SQL driver sempit agar tak mengotori Service yang luas

import (
	"context"
	"database/sql"
	authdomain "yourapp/internal/modules/auth/domain"
)

// SQLRepository mengimplementasikan struktur untuk kebutuhan pembacaan (Read).
// Ini hanya terhubung ke Replica / Read DB karena Autentikasi itu cuma operasi periksa dan baca data.
type SQLRepository struct {
	readDB *sql.DB
}

// NewSQLRepository membangun gudang DB untuk modul autentikasi.
func NewSQLRepository(readDB *sql.DB) *SQLRepository {
	return &SQLRepository{readDB: readDB}
}

// FindAuthByEmail mencari ke barisan log users hanya bermodal pencocokan unik email.
func (r *SQLRepository) FindAuthByEmail(ctx context.Context, email string) (*authdomain.AuthUser, error) {
	const q = `
		SELECT id, email, password, role, status
		FROM users
		WHERE email = ? LIMIT 1
	` // Target khusus kolom rahasia seperti 'password', yang absen di Model repository User biasa.

	var row authdomain.AuthUser
	err := r.readDB.QueryRowContext(ctx, q, email).Scan(
		&row.ID,
		&row.Email,
		&row.PasswordHash, // Hash Bcrypt tersalin kesini
		&row.Role,
		&row.Status,
	)

	if err == sql.ErrNoRows {
		// Nilainya murni kosong (Notifikasi aman ke service)
		return nil, nil
	}
	if err != nil {
		// Error hard disk, crash server, dll
		return nil, err
	}

	// Jika sukses mutlak, lemparkan representasi wujud auth ini
	return &row, nil
}
