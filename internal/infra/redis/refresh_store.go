package redisx // Sengaja digabung ke folder infra/redis sebelumnya supaya tidak membuat cluster file yang mencar-mencar

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9" // Driver murni redis langsung dari pembuat library v9
)

// RefreshStore bertugas merekam histori "Sesi Device Login" yang umurnya sangat lama (Refresh Token).
// Apa filosofi menyimpannya dalam redis? Ini disebut 'Stateful'. Membantu bila user tercuri gawai aslinya:
// kita bisa suruh API menghapus kunci ini di DB redis, jadi secara magis akses pencuri tertolak mentah-mentah
// saat permohonan token akses (Access Token) anyar diajukan kembali nantinya.
type RefreshStore struct {
	client *goredis.Client
}

// NewRefreshStore menelurkan agen pintu loket logis Redis khusus untuk perpanjangan Token
func NewRefreshStore(client *goredis.Client) *RefreshStore {
	return &RefreshStore{client: client}
}

// Save meletakkan kode kunci rahasia Sesi (sessionID) dan mengawinkannya pada umur kadaluwarsa panjang (1 minggu TTL misalnya)
func (s *RefreshStore) Save(ctx context.Context, sessionID string, userID int64, ttl time.Duration) error {
	// Dimasukkan imbuhan nama "refresh:[kode]" agar folder RAM redis rapi, gak tabrakan secara liar dengan modul cache User Module.
	key := fmt.Sprintf("refresh:%s", sessionID)

	// Value (isi) yang disimpan ke string nilainya hanya ID user yang berasosiasi untuk menghemat Memori Redis RAM server Anda secara optimal.
	return s.client.Set(ctx, key, strconv.FormatInt(userID, 10), ttl).Err()
}

// Exists memastikan secara kritis apakah Session ID yang melayang dari JSON Body Klien masih bernyawa / belum pernah dipotong paksa.
func (s *RefreshStore) Exists(ctx context.Context, sessionID string) (bool, error) {
	key := fmt.Sprintf("refresh:%s", sessionID)
	// Return angka n. Jumlah baris Redis yang berhasil terkena Scan Exists.
	n, err := s.client.Exists(ctx, key).Result()

	// Jika  n == 1, maka benar kunci belum basi mati.
	return n > 0, err
}

// Delete membunuh mutlak sesi login ini dari alam data RAM redis API kita.
// Berguna untuk mengimplementasi fitur "Logout" akun klien. Token curian jadi rongsokan tidak berguna lagi seketika.
func (s *RefreshStore) Delete(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("refresh:%s", sessionID)
	return s.client.Del(ctx, key).Err()
}
