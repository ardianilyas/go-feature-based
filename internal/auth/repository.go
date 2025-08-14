package auth

import (
	"github.com/ardianilyas/go-feature-based/config"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (r *repository) CreateUser(user *User) error {
	return config.DB.Create(user).Error
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	if err := config.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}