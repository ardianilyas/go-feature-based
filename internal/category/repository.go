package category

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCategory(category *Category) error
	GetCategoryByID(id uuid.UUID) (*Category, error)
	GetAllCategories() ([]*Category, error)
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

func (r *repository) GetAllCategories() ([]*Category, error) {
	var categories []*Category
	if err := config.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *repository) UpdateCategory(category *Category) error {
	return config.DB.Model(&Category{}).Where("id = ?", category.ID).Updates(category).Error
}

func (r *repository) DeleteCategory(id uuid.UUID) error {
	return config.DB.Delete(&Category{}, "id = ?", id).Error
}