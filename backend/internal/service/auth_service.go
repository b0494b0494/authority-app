package service

import (
	"errors"

	"authority-app/backend/internal/models"
	"authority-app/backend/internal/repository"
	"authority-app/backend/internal/util"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, password string) error {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{Username: username, Password: hashedPassword}
	return s.userRepo.CreateUser(&user)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if err := util.ComparePassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	tokenString, err := util.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
