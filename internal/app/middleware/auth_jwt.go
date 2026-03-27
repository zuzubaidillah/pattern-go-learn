package middleware // Folder penampung perantara pencegat (interceptor) sebelum Request sampai ke Handler utama

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authdomain "yourapp/internal/modules/auth/domain" // Import kerangka batasan Domain Auth
)

// JWTService adalah kontrak syarat minimal (interface) yang dibutuhkan oleh middleware ini.
// Middleware JWT tidak perduli logika databasenya, ia hanya butuh alat pencabut klaim token (ParseAccessToken).
type JWTService interface {
	ParseAccessToken(token string) (*authdomain.AccessClaims, error)
}

// AuthJWT adalah gin.HandlerFunc (Middleware) yang memblokir orang usil.
// Middleware ini diletakkan sebagai pelindung sebuah grup route HTTP spesifik.
func AuthJWT(authSvc JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Cabut nilai Header dari Klien yang bertuliskan "Authorization"
		header := c.GetHeader("Authorization")
		if header == "" {
			// Jika kosong, langsung potong kawat request (Abort) sambil kirim kode HTTP 401
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			return
		}

		// 2. Format standar HTTP OAuth2 / JWT Token berbentuk string: "Bearer xxxxx.yyyyy.zzzzz"
		parts := strings.SplitN(header, " ", 2)
		// Harus ada belahan 2 spasi, dan kata depannya mutlak 'Bearer' (tanpa peduli besar kecil kapital)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header format"})
			return
		}

		// 3. Minta tolong service JWT penelanjangan string dan memvalidasi gembok sandinya
		claims, err := authSvc.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired access token"})
			return
		}

		// 4. Jika gembok sandi rahasia cocok + waktu belum basi, titipkan (Set) data klaim penting ini ke tas belanja Context
		// Agar HTTP Handler selanjutnya gak usah repot-repot query database dari nol untuk mencari tahu "Siapa gue?"
		c.Set("auth.user_id", claims.UserID)
		c.Set("auth.email", claims.Email)
		c.Set("auth.role", claims.Role)

		// 5. Izinkan request meluncur maju ke handler fungsional (seperti Update Profile/Delete dsb)
		c.Next()
	}
}
