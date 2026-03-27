package domain // Package domain menampung kontrak dan bentuk dasar Entitas abstrak

import "github.com/golang-jwt/jwt/v5"

// AccessClaims merupakan struktur izin JWT khusus untuk Token Akses.
// Kita memasukkan spesifikasi klaim `RegisteredClaims` milik library JWT (seperti kapan dibuat, audiensnya siapa)
// lalu menambahkan informasi unik bisnis kita sendiri: UserID, Email, dan Role.
type AccessClaims struct {
	UserID int64  `json:"user_id"` // Sangat penting agar handler mudah mengambil profil yang bersangkutan
	Email  string `json:"email"`
	Role   string `json:"role"`
	Type   string `json:"type"` // Ditandai hardcode "access" untuk pengamanan
	jwt.RegisteredClaims
}

// RefreshClaims khusus digunakan untuk validasi perpanjangan token lama.
type RefreshClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"` // ID Unik untuk dicocokkan ke baris tabel Redis
	Type      string `json:"type"`       // Ditandai hardcode "refresh" supaya token ini tidak bisa di-bypass dipakai ke /me
	jwt.RegisteredClaims
}

// TokenPair adalah model sederhana bungkus balasan ke Endpoint Login/Refresh.
type TokenPair struct {
	AccessToken  string // Dikirim ke Authorization header klien nantinya
	RefreshToken string // Disimpan di storage klien dan sesekali ditukarkan
}
