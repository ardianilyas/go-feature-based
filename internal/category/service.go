package category

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/pkg/pagination"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CreateCategory(category *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(page, limit int, baseURL string) (pagination.PaginatedResult[Category], error)
	UpdateCategory(category *Category) error
	DeleteCategory(id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateCategory(category *Category) error {
	err := s.repo.CreateCategory(category)
	if err == nil {
		config.Log.WithFields(logrus.Fields{
			"category_id": category.ID,
			"name":        category.Name,
		}).Info("Category created")
	}
	return err
}

func (s *service) GetCategoryByID(id uuid.UUID) (*Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *service) GetAllCategories(page, limit int, baseURL string) (pagination.PaginatedResult[Category], error) {
	return s.repo.GetAllCategories(page, limit, baseURL)
}

func (s *service) UpdateCategory(category *Category) error {
	oldCategory, err := s.repo.GetCategoryByID(category.ID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateCategory(category)
	if err == nil {
		config.Log.WithFields(logrus.Fields{
			"category_id": category.ID,
			"old_name":    oldCategory.Name,
			"new_name":    category.Name,
		}).Info("Category updated")
	}
	return err
}

func (s *service) DeleteCategory(id uuid.UUID) error {
	category, err := s.repo.GetCategoryByID(id)
	if err != nil {
		return err
	}
	if err := s.repo.DeleteCategory(id); err != nil {
		return err
	}
	config.Log.WithFields(logrus.Fields{
		"category_id": id,
		"name":        category.Name,
	}).Info("Category deleted")
	return nil
}