package domain

import "context"

// Repository adalah sebuah Kontrak / Perjanjian antarmuka (Interface)
// Di dalam clean architecture, Domain hanya menetapkan *syarat* atau "apa yang harus bisa dilakukan pada User oleh storage".
//
// Siapa yang mengimplementasikannya nanti (Misal SQLRepository atau NoSQLRepository)? Domain tidak mau peduli.
// Hal tersebut membebaskan sistem Service kita dirotasi databasenya kelak tanpa ubah kode Domain.
type Repository interface {
	// FindByID harus mencari user berdasarkan ID unik
	FindByID(ctx context.Context, id int64) (*User, error)
	// FindByEmail harus mencari user berdasarkan alamat Email
	FindByEmail(ctx context.Context, email string) (*User, error)
	// Create harus mendelegasikan perintah menyimpan entitas User baru ke database
	Create(ctx context.Context, user *User) error
	// Update harus memutakhirkan info pada entri User yang sudah ada
	Update(ctx context.Context, user *User) error
}
