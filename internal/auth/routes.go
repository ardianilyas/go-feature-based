package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, service Service) {
	h := NewHandler(service)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}