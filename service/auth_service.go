package service

import (
	"errors"
	"time"

	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(user *entity.User) (string, error) // return token or error
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(user *entity.User) (string, error) {
	// Check duplicate
	existing, err := s.userRepo.FindByUsernameOrEmail(user.Username, user.Email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", errors.New("username or email has already been used")
	}

	// Utils : Hash Password
	if err := utils.HashPassword(user, user.Password); err != nil {
		return "", err
	}

	// Mapping
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.TelegramIsValid = false

	// Query : Create Register
	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	// Utils : Generate Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(email, password string) (string, error) {
	// Query : Check User By Email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email")
	}

	// Utils : Compare Password
	if err := utils.CheckPassword(user, password); err != nil {
		return "", errors.New("invalid password")
	}

	// Utils : Generate Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
