package post

import (
	"github.com/ardianilyas/go-feature-based/pkg/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRotes(r *gin.Engine) {
	h := NewHandler(NewService(NewRepository()))

	posts := r.Group("/posts")
	posts.Use(middlewares.JWTAuth())
	{
		posts.POST("/", h.CreatePost)
		posts.GET("/", h.GetAllPosts)
		posts.GET("/:id", h.GetPostByID)
		posts.PUT("/:id", h.UpdatePost)
		posts.DELETE("/:id", h.DeletePost)
	}
}