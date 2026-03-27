package domain

import "context"

// Service adalah Kontrak yang mewakili murni "Kasus Penggunaan" (Use Case) bisnis kita.
// Service menjabarkan kemampuan apa saja fitur User yang bisa dipanggil oleh HTTP Handler.
type Service interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, req CreateUserInput) (*User, error)
	Update(ctx context.Context, req UpdateUserInput) (*User, error)
}

// CreateUserInput adalah struktur Data Masukan Murni khusus untuk ke Service (Use Case).
//
// Kenapa repot bikin ini lagi? Kenapa HTTP DTO tidak langsung dilempar ke dalam Repository?
// Jawaban: Agar Service kita tetap *mandiri* / suci dari binding atribut web (json, form-data, dll)
// dan bisa dipanggil ulang tanpa error, misal via Event Driven gRPC / Cron Worker lokal!
type CreateUserInput struct {
	Name   string
	Email  string
	Status string
}

// UpdateUserInput mewakili input aman dari dunia luar untuk pengubahan profil User.
type UpdateUserInput struct {
	ID     int64
	Name   string
	Status string
}
