package service

import (
	"context"
	"errors"
	"time"

	"pelita/config"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthService interface {
	Register(user *entity.User) (string, error)
	Login(email, password string) (string, error)
	SignOut(token string) error
}

type authService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewAuthService(userRepo repository.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
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

func (s *authService) SignOut(tokenString string) error {
	// Token Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("missing exp in token")
	}

	// Check If Token Expired
	expTime := time.Unix(int64(expFloat), 0)
	duration := time.Until(expTime)
	if duration <= 0 {
		return errors.New("token already expired")
	}

	// Save token to Redis blacklist
	err = s.redisClient.Set(context.Background(), tokenString, "blacklisted", duration).Err()
	if err != nil {
		return errors.New("failed to blacklist token")
	}

	return nil
}
