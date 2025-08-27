package category

import (
	"github.com/ardianilyas/go-feature-based/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	h := NewHandler(NewService(NewRepository()))

	categories := r.Group("/categories")
	categories.Use(middlewares.JWTAuth(), middlewares.RequireRole("admin"))
	{
		categories.POST("/", h.CreateCategory)
		categories.GET("/", h.GetAllCategories)
		categories.GET("/:id", h.GetCategoryByID)
		categories.PUT("/:id", h.UpdateCategory)
		categories.DELETE("/:id", h.DeleteCategory)
	}
}