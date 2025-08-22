package category

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/ardianilyas/go-feature-based/pkg/pagination"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCategory(category *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories(page, limit int, baseURL string) (pagination.PaginatedResult[Category], error)
	UpdateCategory(category *Category) error
	DeleteCategory(id uuid.UUID) error
}

type repository struct {}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateCategory(category *Category) error {
	return config.DB.Create(category).Error
}

func (r *repository) GetCategoryByID(id uuid.UUID) (*Category, error) {
	var category Category
	if err := config.DB.First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *repository) GetAllCategories(page, limit int, baseURL string) (pagination.PaginatedResult[Category], error) {
	return pagination.Paginate[Category](config.DB, page, limit, baseURL)
}

func (r *repository) UpdateCategory(category *Category) error {
	return config.DB.Model(&Category{}).Where("id = ?", category.ID).Updates(category).Error
}

func (r *repository) DeleteCategory(id uuid.UUID) error {
	return config.DB.Delete(&Category{}, "id = ?", id).Error
}