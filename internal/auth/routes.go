package auth

import (
	"github.com/ardianilyas/go-feature-based/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, service Service) {
	h := NewHandler(service)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}

	protected := r.Group("/auth")
	protected.Use(middlewares.JWTAuth())
	{
		protected.GET("/profile", h.Profile)
	}

	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuth(), middlewares.RequireRole("admin"))
	admin.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello Admin!"})
	})
}