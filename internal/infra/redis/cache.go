package redisx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache merupakan pelapis objek client (Wrapper Object). Tujuannya menyembunyikan
// sintaks redis murni menjadi fungsi yang lebih masuk akal dimengerti alur bisnis (SetJSON, dsb).
type Cache struct {
	client *redis.Client
}

// NewCache membungkus klien murni.
func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// SetJSON adalah cara praktis melempar Struct GO atau Map GO apapun dan menyimpannya
// dengan kadaluwarsa (TTL / Time To Live) menjadi teks JSON murni langsung ke RAM Redis.
func (c *Cache) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value) // Ekstraksi isi variabel `value` secara reflektif menjadi bentuk teks `[{..}]`
	if err != nil {
		return err
	}
	// Perintah redis standar SET dengan atribut waktu nyala
	return c.client.Set(ctx, key, b, ttl).Err()
}

// GetJSON bertugas kebalikannya. Mencari Key redis di RAM, menyedot string JSON nya,
// lalu ditembakkan (Unmarshal) isinya untuk menyetel struktur kosong dari objek di pointer `dest`.
func (c *Cache) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	val, err := c.client.Get(ctx, key).Result()

	// Secara spesial mendeteksi jika kuncinya tidak tersimpan di memori cache (Redis nil)
	if err == redis.Nil {
		return false, nil // Data tak ada
	}
	if err != nil {
		return false, err // Redis error lain misal internet mati
	}

	// Jika berhasil menyedot, salin/decode datanya ke rumah memori tujuan (dest)
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}
	return true, nil // Status data sukses ditemukan
}

// Delete biasa dipakai ketika ada update data user (Invalidation), sehingga cache usang harus dihancurkan.
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err() // Dukung variadic `keys ...string` agar bisa hapus borongan multiple kunci
}
