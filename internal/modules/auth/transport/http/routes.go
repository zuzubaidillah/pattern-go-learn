package http

import (
	"github.com/gin-gonic/gin"
	"yourapp/internal/app/middleware"
)

// RegisterRoutes membuhul rute-rute publik maupun rute terikat Middleware.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, authSvc middleware.JWTService) {
	// Endpoint Publik (Masuk)
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.Login)     // Meracik kunci
		auth.POST("/refresh", h.Refresh) // Menukar tambah kunci
		auth.POST("/logout", h.Logout)   // Menghancurkan kunci
	}

	// Endpoint Privat (Keamanan Berlapis)
	private := rg.Group("/me")
	private.Use(middleware.AuthJWT(authSvc)) // Menjaga gerbang `/me` dengan Satpam JWT!
	{
		private.GET("", h.Me) // Membedah diri sendiri
	}
}
