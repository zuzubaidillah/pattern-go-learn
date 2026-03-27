package redisx // Sengaja dinamakan redisx (redis extended) agar tidak bertabrakan dengan package "github.com/redis/go-redis"

import (
	"context"

	"github.com/redis/go-redis/v9" // Client standar komunitas go-redis versi 9 (terbaru)
	"yourapp/internal/config"      // Integrator pembaca config YAML
)

// NewClient membangun akses penghubung ke server Redis
func NewClient(cfg config.RedisConfig) (*redis.Client, error) {
	// Menginstansiasi client menggunakan otentikasi (password & alamat)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB, // Biasanya DB 0
	})

	// Melakukan PING dengan konteks kosong untuk memvalidasi server redisnya nyata-nyata hidup
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil // Kembalikan koneksi mentah Redis
}
