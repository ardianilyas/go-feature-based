package internal

import (
	"net/http"

	"github.com/ardianilyas/go-feature-based/internal/auth"
	"github.com/ardianilyas/go-feature-based/internal/category"
	"github.com/ardianilyas/go-feature-based/internal/post"
	"github.com/ardianilyas/go-feature-based/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	testRoutes(r)
	authRoutes(r)
	adminRoutes(r)
	categoriesRoute(r)
	postsRoute(r)
}

func testRoutes(r *gin.Engine) {
	test := r.Group("/test")
	test.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello test"})
	})
}

func authRoutes(r *gin.Engine) {
	authRepo := auth.NewRepository()
	authService := auth.NewService(authRepo)
	auth.RegisterRoutes(r, authService)
}

func adminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(middlewares.JWTAuth(), middlewares.RequireRole("admin"))
	admin.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello admin"})
	})
}

func categoriesRoute(r *gin.Engine) {
	category.RegisterRoutes(r)
}

func postsRoute(r *gin.Engine) {
	post.RegisterRotes(r)
}