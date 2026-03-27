package http // Menyambungkan alur murni Domain Service ke HTTP Framework spesifik (Gin)

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "yourapp/internal/modules/auth/domain"
	"yourapp/internal/pkg/contextutil"
	"yourapp/internal/pkg/response"
)

// Handler menaungi layanan logika bisnis
type Handler struct {
	service authdomain.Service
}

// NewHandler mencetak rupa Handler API
func NewHandler(service authdomain.Service) *Handler {
	return &Handler{service: service}
}

// Login menangkap masukan "email & password" dan mengembalikan Token Akses-Penyegar jika sukses.
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	// Validasi wujud JSON sesuai tag min=6, email format, dll
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	// Teruskan logika ke tulang punggung Service
	tokenPair, err := h.service.Login(c.Request.Context(), authdomain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		// Menggagalkan kredensial 401
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Mengabarkan berhasil 200 OK dengan seragam TokenResponse DTO
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer", // Standardisasi W3C untuk HTTP Headers
	})
}

// Refresh menampung permohonan Token Akses anyar bersenjatakan Refresh Token basi.
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	// Saring permintaan via service
	tokenPair, err := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid refresh token or session expired")
		return
	}

	// Menghidangkan sajian Token siklus terbaru!
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    "Bearer",
	})
}

// Logout melaikkan penghancuran session token dalam RAM Redis.
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid refresh token")
		return // Usaha pemecahan Token yang salah tak digubris
	}

	// Karena sudah Logout, cukuplah membalas salam damai.
	response.Success(c, http.StatusOK, "logout success")
}

// Me mengidentifikasi data suci pribadi berdasarkan titipan JWT Middleware (auth.user_id)
func (h *Handler) Me(c *gin.Context) {
	// Panggilan super bersih via pembongkar tas `contextutil`! Tanpa query database sama sekali!
	userID := contextutil.MustUserID(c)
	email := contextutil.GetEmail(c)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
