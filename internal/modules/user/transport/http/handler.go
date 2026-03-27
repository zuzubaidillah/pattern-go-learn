// Package http merangkul dan menerjemahkan segala sesuatu perihal teknologi berbasis Web / HTTP
// (Status HTTP seperti 200, 400, 500, request context, param ID dari URL)
package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"yourapp/internal/modules/user/domain"
)

// Handler menjembatani aliran Web dari library `Gin` langsung ke Domain `sistem internal Service`.
type Handler struct {
	service domain.Service
}

func NewHandler(service domain.Service) *Handler {
	return &Handler{service: service}
}

// GetByID memakan parameter URL dinamis (e.g. GET /api/v1/users/99)
func (h *Handler) GetByID(c *gin.Context) {
	// Pengecekan apakah URL "?id=" itu benar-benar murni angka string -> int64.
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	// Teruskan pencarian ke service yang pintar
	user, err := h.service.GetByID(c.Request.Context(), id)

	// Cocokkan status HTTP response ke Klien:
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return // 500 Mesin Server Error
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return // 404 Tak Ditemukan
	}

	// 200 Sukses. Pakai DTO UserResponse agar kolom-kolomnya rapi & terbatas.
	c.JSON(http.StatusOK, UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	})
}

// Create menangani pengikatan Payload Request Tipe JSON (Body)
func (h *Handler) Create(c *gin.Context) {
	var req CreateUserRequest
	// Validasi form otomatis mengikuti tag di `dto.go`
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), domain.CreateUserInput{
		Name:   req.Name,
		Email:  req.Email,
		Status: req.Status,
	})

	if err != nil {
		// Evaluasi manual apabila pentalan service menghasilkan kalimat konflik ganda "email already exists"
		if err.Error() == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"message": err.Error()}) // 409 Conflict
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// 201 Created (Data Baru diciptakan)
	c.JSON(http.StatusCreated, UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	})
}

// Update memadukan Path Parameter ID + Body Request Tipe JSON.
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request", "error": err.Error()})
		return
	}

	// Ubah parameter
	user, err := h.service.Update(c.Request.Context(), domain.UpdateUserInput{
		ID:     id,
		Name:   req.Name,
		Status: req.Status,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	})
}
