package post

import (
	"errors"
	"net/http"

	"github.com/ardianilyas/go-feature-based/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if !utils.ValidateStruct(c, &req) {
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := &Post{
		Title: req.Title,
		Content: req.Content,
		CategoryID: req.CategoryID,
		UserID: userID,
	}

	createdPost, err := h.service.CreatePost(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post created successfully",
		"post": createdPost,
	})
}

func (h *Handler) GetPostByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (h *Handler) UpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req UpdatePostRequest
	if !utils.ValidateStruct(c, &req) {
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post := &Post{ID: id, Title: req.Title, Content: req.Content, CategoryID: req.CategoryID}

	updatedPost, err := h.service.UpdatePost(post, userID)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post updated successfully",
		"post": updatedPost,
	})
}

func (h *Handler) DeletePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeletePost(id, userID); err != nil {
		if errors.Is(err, ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else if errors.Is(err, ErrForbidden) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
	})
}