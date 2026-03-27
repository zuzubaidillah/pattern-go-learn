package http

import "github.com/gin-gonic/gin"

// RegisterRoutes adalah kolektor URL Endpoint untuk entitas Users (!= Controller).
// rg = Engine Inti Router Gin
// h  = Penanggung Jawab Proses Logika
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	// Menambahkan prefix sub-url grup kepada setiap alamat khusus User API ("/api/v1/uers")
	users := rg.Group("/users")
	{
		users.GET("/:id", h.GetByID) // GET /api/v1/users/5 -> Menampilkan
		users.POST("", h.Create)     // POST /api/v1/users -> Membuat
		users.PUT("/:id", h.Update)  // PUT /api/v1/users/5 -> Mengubah Keseluruhan Identitas
	}
	// Di masa depan Anda dapat memblok akses Endpoint Modul ini dengan cara menyisipkan Middleware
	// Token Otorisasi di tingkat `users.Use(middleware)` pada baris grup router!
}
