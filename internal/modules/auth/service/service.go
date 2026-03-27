package service // Inti gravitasi dari logika masuk/keluarnya otentikasi user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"yourapp/internal/config"
	authdomain "yourapp/internal/modules/auth/domain"
)

// Antarmuka (Interface) berikut dicomot (diinjeksi) kepemilikannya oleh Service,
// agar Service tak pedulikan wujud aslinya (Bcrypt kah, Redis kah, SQL kah), yang penting fungsinya seragam.

type RefreshStore interface {
	Save(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, sessionID string) (bool, error)
	Delete(ctx context.Context, sessionID string) error
}

type PasswordVerifier interface {
	Verify(hashedPassword, plainPassword string) bool
}

type AuthRepository interface {
	FindAuthByEmail(ctx context.Context, email string) (*authdomain.AuthUser, error)
}

// Service menyimpan ikatan tali-temali alat konfigurasinya.
type Service struct {
	cfg           config.JWTConfig
	userRepo      AuthRepository
	refreshStore  RefreshStore
	passwordCheck PasswordVerifier
}

// New merakit jalinan logika dan memori (Service Layer).
func New(
	cfg config.JWTConfig,
	userRepo AuthRepository,
	refreshStore RefreshStore,
	passwordCheck PasswordVerifier,
) *Service {
	return &Service{
		cfg:           cfg,
		userRepo:      userRepo,
		refreshStore:  refreshStore,
		passwordCheck: passwordCheck,
	}
}

// Login mengelola jalur kedatangan User pertama kali. Merupakan perwujudan `domain.Service`.
func (s *Service) Login(ctx context.Context, req authdomain.LoginInput) (*authdomain.TokenPair, error) {
	// 1. Cari data murni email di Database
	user, err := s.userRepo.FindAuthByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// Cegat jika user tak ada atau akunnya tidak berstatus aktif
	if user == nil || user.Status != "active" {
		// Sengaja kita membalas error samar demi alasan keamanan mutlak.
		// Jangan spesifik membalas "Email tak ditemukan", karena mudah diserang Enumerasi Target (Scraping)!
		return nil, errors.New("invalid credentials")
	}

	// 2. Adu ketikan kata sandi (Plain) terhadap Enkripsi rumit di tabel (Hash).
	if !s.passwordCheck.Verify(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// 3. Jika mulus tanpa cacat, racik sepasang token kunci akses untuknya.
	return s.generateTokenPair(ctx, user.ID, user.Email, user.Role)
}

// Refresh tugasnya cuma satu: Memberi AccessToken baru BILA RefreshToken masih sah belom basi.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (*authdomain.TokenPair, error) {
	claims := &authdomain.RefreshClaims{}

	// Buka gembok token menggunakan kunci rahasia "Refresh" (Sangat berbeda jalurnya dengan alat cek Access!)
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.RefreshSecret), nil
	}, jwt.WithAudience(s.cfg.Audience), jwt.WithIssuer(s.cfg.Issuer))

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Pastikan jenisnya benar-benar "refresh", mencegah hacker curang memutarbalikkan pakai AccessToken
	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Tanyakan log Redis: Apakah sesi ID token ini masih terdaftar atau jangan-jangan udah di-logout minggat?
	exists, err := s.refreshStore.Exists(ctx, claims.SessionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("refresh session revoked") // Tertolak mutlak
	}

	// Karena norma kemanan modern mengharuskan setiap Refresh cuma boleh dipakai SEKALI SAJA (Rotasi Keamanan Oauth2), kita bantai kunci redis lamanya.
	_ = s.refreshStore.Delete(ctx, claims.SessionID)

	// Berkelanjutan ke pembuahan Pasangan Token (Token Pair) baru yang tereset kembali umur tayangnya.
	return s.generateTokenPair(ctx, claims.UserID, "", "")
}

// ParseAccessToken mendedah token JWT demi memuluskan persyaratan Middleware pelindung API kita.
func (s *Service) ParseAccessToken(tokenStr string) (*authdomain.AccessClaims, error) {
	claims := &authdomain.AccessClaims{}

	// Buka kunci token memakai rahasia Khusus "Access"
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.AccessSecret), nil
	}, jwt.WithAudience(s.cfg.Audience), jwt.WithIssuer(s.cfg.Issuer))

	if err != nil || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

// Logout mengakhiri paksa sesi melalui pengeboman cache sesi di RAM Redis.
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	claims := &authdomain.RefreshClaims{}

	// Ekstraksi info Session ID milik token itu sendiri (Pastikan yang cabut token adalah pembawa Token yang murni sah)
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(s.cfg.RefreshSecret), nil
	}, jwt.WithAudience(s.cfg.Audience), jwt.WithIssuer(s.cfg.Issuer))

	if err != nil || !token.Valid {
		return errors.New("invalid refresh token")
	}

	// Hapus jejak redis-nya! Otomatis device apapun yang simpan refresh token session ID tersebut bakal ditolak mesin JWT.
	return s.refreshStore.Delete(ctx, claims.SessionID)
}

// generateTokenPair meracik panasea ramuan rahasia keamanannya. (Fungsi internal kelas Private).
func (s *Service) generateTokenPair(ctx context.Context, userID int64, email, role string) (*authdomain.TokenPair, error) {
	now := time.Now()
	accessExp := now.Add(s.cfg.AccessTokenTTL)   // Misal setel +15 menit dari sekarang
	refreshExp := now.Add(s.cfg.RefreshTokenTTL) // Misal setel +7 Hari dari sekarang

	// Bikin Identifier Sesi Acak sepanjang 16 karakter yang didedikasikan mandiri untuk Token Refresh ini (guna penyimpanan log Redis)
	sessionID, err := randomHex(16)
	if err != nil {
		return nil, err
	}

	// 1. Perancangan Dokumen Identitas (Access Claims)
	accessClaims := authdomain.AccessClaims{
		UserID: userID, // Dicetak jelas
		Email:  email,
		Role:   role,
		Type:   "access", // Penanda statik
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.Issuer,
			Subject:   toString(userID),
			Audience:  jwt.ClaimStrings{s.cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(accessExp), // Penempelan Umur Expired
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now), // Hanya berlaku sejagat alam go kita dan harus dimulai dari tanda waktu detik ini terjadi!
		},
	}

	// 2. Perancangan Surat Kuasa Tinggal (Refresh Claims). Tubuhnya minimalis tidak seperti Access Token.
	refreshClaims := authdomain.RefreshClaims{
		UserID:    userID,
		SessionID: sessionID,
		Type:      "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.cfg.Issuer,
			Subject:   toString(userID),
			Audience:  jwt.ClaimStrings{s.cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Tanda tangani / Segel surat klaim terbuat ke dalam pita Enkripsi terlipat berbasis Kriptografi HS256 + Pasword Khusus Secret!
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.cfg.AccessSecret))
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.cfg.RefreshSecret))
	if err != nil {
		return nil, err
	}

	// Lempar cetakan "Sesi Aktif" perintis ke sistem antrean Redis seumur waktu yang dicanangkan di konfigurasi lokal
	if err := s.refreshStore.Save(ctx, sessionID, userID, s.cfg.RefreshTokenTTL); err != nil {
		return nil, err
	}

	return &authdomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// randomHex mencetak bilangan string alfanumerik secara instan murni acak kriptografis.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// toString pembungkus sepele mengubah int.
func toString(v int64) string {
	return fmt.Sprintf("%d", v)
}
