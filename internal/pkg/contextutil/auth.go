package contextutil // Menampung alat-alat perkakas khusus urusan membongkar Gin Context (tas map data transit Request)

import "github.com/gin-gonic/gin"

// MustUserID memaksa merogoh data nilai UserID dari tas c.Set() yang sebelumnya dimasukkan secara senyap oleh Middleware JWT.
// Ini bikin handler jadi rapi karena tinggal memanggil contextutil.MustUserID(c).
func MustUserID(c *gin.Context) int64 {
	v, _ := c.Get("auth.user_id") // Mengambil properti abstrak yang tertanam di memori gin context (sebagai interface{})
	id, _ := v.(int64)            // Konversi paksa dari objek samar 'any/interface' ke integer 64 (Safe Cast: kalau gagal, id otomatis 0)
	return id
}

// GetEmail mengambil properti pembungkus String alamat e-mail user.
func GetEmail(c *gin.Context) string {
	v, _ := c.Get("auth.email")
	email, _ := v.(string)
	return email
}
