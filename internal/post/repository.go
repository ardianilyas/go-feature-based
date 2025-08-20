package post

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/google/uuid"
)

type Repository interface {
	CreatePost(post *Post) error
	GetPostByID(id uuid.UUID) (*Post, error)
	GetAllPosts() ([]*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id uuid.UUID) error
}

type repository struct {}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreatePost(post *Post) error {
	return config.DB.Create(post).Error
}

func (r *repository) GetPostByID(id uuid.UUID) (*Post, error) {
	var post Post
	if err := config.DB.Preload("User").Preload("Category").First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *repository) GetAllPosts() ([]*Post, error) {
	var posts []*Post
	if err := config.DB.Preload("User").Preload("Category").Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *repository) UpdatePost(post *Post) error {
	return config.DB.Model(&Post{}).Where("id = ?", post.ID).Updates(post).Error
}

func (r *repository) DeletePost(id uuid.UUID) error {
	return config.DB.Delete(&Post{}, "id = ?", id).Error
}