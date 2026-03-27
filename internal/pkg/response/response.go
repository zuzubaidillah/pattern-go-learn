package response // Package ini bertujuan membakukan bentuk objek JSON balasan API

import "github.com/gin-gonic/gin"

// Success dipakai oleh semua API Handlers yang ingin membalas JSON dengan status sukses (200, 201)
// Hasilnya akan selalu berbentuk: { "data": [...] }
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, gin.H{"data": data}) // Mengikat HTTP Status Code & JSON Body menggunakan fitur default Gin
}

// Error dipakai API Handlers yang ingin membalas pesan gagal (400, 404, 500)
// Hasilnya seragam seperti: { "error": "pesan teks" }
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
