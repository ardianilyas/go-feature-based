package category

import "github.com/google/uuid"

type Service interface {
	CreateCategory(category *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories() ([]*Category, error)
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
	return s.repo.CreateCategory(category)
}

func (s *service) GetCategoryByID(id uuid.UUID) (*Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *service) GetAllCategories() ([]*Category, error) {
	return s.repo.GetAllCategories()
}

func (s *service) UpdateCategory(category *Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *service) DeleteCategory(id uuid.UUID) error {
	return s.repo.DeleteCategory(id)
}