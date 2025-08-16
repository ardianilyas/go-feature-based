package auth

import (
	"errors"
	"os"

	"github.com/ardianilyas/go-feature-based/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterRequest) error
	Login(input LoginRequest) (*User, string, string, error)
	RefreshToken(refreshToken string) (*User, string, error)
	GetProfile(userID uuid.UUID) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Register(input RegisterRequest) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		ID: uuid.New(),
		Name: input.Name,
		Email: input.Email,
		Password: string(hash),
	}

	return s.repo.CreateUser(&user)
}

func (s *service) Login(input LoginRequest) (*User, string, string, error) {
	user, err := s.repo.GetUserByEmail(input.Email)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), user.Role)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *service) RefreshToken(refreshToken string) (*User, string, error) {
	claims, err := utils.ParseToken(refreshToken, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, "", err
	}

	user, err := s.repo.GetUserByID(claims.ID)
	if err != nil {
		return nil, "", errors.New("user not found")
	}

	newAccessToken, err := utils.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	return user, newAccessToken, nil
}

func (s *service) GetProfile(userID uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}