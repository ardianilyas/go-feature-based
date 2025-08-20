package post

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("forbidden")
)

type Service interface {
	CreatePost(post *Post) (*Post, error)
	GetPostByID(id uuid.UUID) (*Post, error)
	GetAllPosts() ([]*Post, error)
	UpdatePost(post *Post, userID uuid.UUID) (*Post, error)
	DeletePost(id, userID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreatePost(post *Post) (*Post, error) {
	if err := s.repo.CreatePost(post); err != nil {
		return nil, err
	}

	return s.repo.GetPostByID(post.ID)
}

func (s *service) GetPostByID(id uuid.UUID) (*Post, error) {
	post, err := s.repo.GetPostByID(id)
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	return post, nil
}

func (s *service) GetAllPosts() ([]*Post, error) {
	return s.repo.GetAllPosts()
}

func (s *service) UpdatePost(post *Post, userID uuid.UUID) (*Post, error) {
	existingPost, err := s.repo.GetPostByID(post.ID)
	if err != nil {
		return nil, ErrPostNotFound
	}
	if existingPost.UserID != userID {
		return nil, ErrForbidden
	}
	if err := s.repo.UpdatePost(post); err != nil {
		return nil, err
	}
	return s.repo.GetPostByID(post.ID)
} 

func (s *service) DeletePost(id, userID uuid.UUID) error {
	existingPost, err := s.repo.GetPostByID(id)
	if err != nil {
		return ErrPostNotFound
	}
	if existingPost.UserID != userID {
		return ErrForbidden
	}
	return s.repo.DeletePost(id)
}