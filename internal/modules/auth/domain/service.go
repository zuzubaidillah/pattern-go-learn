package domain

import "context"

// Service menentukan perihal apa saja yang BOLEH dilakukan sistem modul ini
// (Login dari awal, Melakukan Tukar Tambah Token, Mencabut Klaim Token, dan Logout Paksa)
type Service interface {
	Login(ctx context.Context, req LoginInput) (*TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	ParseAccessToken(token string) (*AccessClaims, error)
	Logout(ctx context.Context, refreshToken string) error
}

// LoginInput mewakili data mentah kredensial dari handler yang telah diverifikasi kelengkapannya.
type LoginInput struct {
	Email    string
	Password string
}
