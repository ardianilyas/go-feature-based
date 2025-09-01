package post

import (
	"errors"

	"github.com/ardianilyas/go-feature-based/config"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

	config.Log.WithFields(logrus.Fields{
		"post_id": post.ID,
		"title":   post.Title,
		"user_id": post.UserID,
	}).Info("Post created")
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

	config.Log.WithFields(logrus.Fields{
		"post_id":         post.ID,
		"user_id":         userID,
		"old_title":       existingPost.Title,
		"new_title":       post.Title,
		"old_category_id": existingPost.CategoryID,
		"new_category_id": post.CategoryID,
	}).Info("Post updated")
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
	if err := s.repo.DeletePost(id); err != nil {
		return err
	}
	config.Log.WithFields(logrus.Fields{
		"post_id": id,
		"title":   existingPost.Title,
		"user_id": userID,
	}).Info("Post deleted")
	return nil
}