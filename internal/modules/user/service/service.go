// Package service menghubungkan kebutuhan Bisnis dan meramu interaksi antara:
// (User, Redis Cache, Database).
// Alur berpikir Clean Architecture diletakkan sebagian besar di sini secara steril (tanpa ada referensi HTTP).
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	redisx "yourapp/internal/infra/redis"
	"yourapp/internal/modules/user/domain"
)

// Service menyimpan ikatan objek alat / senjata apa yang dibutuhkannya (repository & cache).
type Service struct {
	repo  domain.Repository // Kontrak tak terikat tipe SQL/Mongo, yang penting cocok format Input/Output
	cache *redisx.Cache     // Integrasi Redis
}

// New menghidupkan dan menyuntik objek layanan ini persediaan alat-alatnya.
func New(repo domain.Repository, cache *redisx.Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

// GetByID adalah proses cerdas pencarian user hibrida.
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	var cached domain.User

	// 1. Tanya Redis: "Hei Redis, apa kamu simpan text json dengan nama `user:x`?"
	found, err := s.cache.GetJSON(ctx, cacheKey, &cached)
	if err == nil && found {
		return &cached, nil // Cache Hit! (Berhasil ambil, potong waktu proses 50x lebih cepat)
	}

	// 2. Jika tak ada di memori Redis Cache(Miss!), perintahkan Repository untuk mencarinya ke dalam lubang Database SQL Read.
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Benar-benar user invalid di universe ini.
	}

	// 3. Simpan di Redis Cache dengan durasi 5 menit.
	// Kita abaikan eror simpan (`_ = s.cache...`) supaya fungsi utama "Mengambil User" tetap sukses
	// (Cache hanyalah optimasi tambahan, tidak boleh bikin transaksi eror mutlak).
	_ = s.cache.SetJSON(ctx, cacheKey, user, 5*time.Minute)

	// 4. Setelah memanaskan Cache, berikan hasilnya.
	return user, nil
}

// Create memeriksa apakah Email tersebut unik. Email tidak boleh duplikat.
func (s *Service) Create(ctx context.Context, req domain.CreateUserInput) (*domain.User, error) {
	// Teropong Database, apakah email "A" udah ada?
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already exists") // Mencegah duplikasi data sebelum kena reject SQL
	}

	// Petakan `req` menjadi object domain baru
	user := &domain.User{
		Name:   req.Name,
		Email:  req.Email,
		Status: req.Status,
	}

	// Mintakan repository mendaftarkannya ke database SQL Writer
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err // Mungkin tabel belum tersedia/Crash.
	}

	return user, nil
}

// Update mutlak wajib melakukan Pembersihan Cache (Cache Invalidation) agar Redis tak tertinggal oleh versi MySQL.
func (s *Service) Update(ctx context.Context, req domain.UpdateUserInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil // Bila ID asal ketik, hentikan
	}

	// Pematrian status ke database model di memori lokal Go
	user.Name = req.Name
	user.Status = req.Status

	// Perintahkan Repository melempar SQL Query UPDATE yang asli di mesin MySQL
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	// [PENTING]: Hapus/Lumpuhkan (Invalidate) redis lama versi sebelum diubah tadi!
	// Jika tidak, API GetByID selanjutnya akan tetap membalas versi usang selama > 5 Menit.
	_ = s.cache.Delete(ctx, fmt.Sprintf("user:%d", user.ID))

	// Balas data baru ke HTTP
	return user, nil
}
