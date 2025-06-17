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

// Auth Interface
type AuthService interface {
	Register(user *entity.User) (string, error)
	Login(email, password string) (string, string, error)
	SignOut(token string) error
}

// Auth Struct
type authService struct {
	userRepo       repository.UserRepository
	adminRepo      repository.AdminRepository
	technicianRepo repository.TechnicianRepository
	redisClient    *redis.Client
}

// Auth Constructor
func NewAuthService(userRepo repository.UserRepository, adminRepo repository.AdminRepository, technicianRepo repository.TechnicianRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepo:       userRepo,
		adminRepo:      adminRepo,
		technicianRepo: technicianRepo,
		redisClient:    redisClient,
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

	// Repo : Create Register
	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	// Utils : Generate Token
	token, err := utils.GenerateToken(user.ID, "user")
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(email, password string) (string, string, error) {
	// Model
	var account entity.Account
	var role string

	// Repo : Check Admin By Email
	admin, err := s.adminRepo.FindByEmail(email)
	if err != nil {
		return "", "", err
	}
	if admin != nil {
		account = admin
		role = "admin"
	}

	// Repo : Check Technician By Email
	if account == nil {
		technician, err := s.technicianRepo.FindByEmail(email)
		if err != nil {
			return "", "", err
		}
		if technician != nil {
			account = technician
			role = "technician"
		}
	}

	// Repo : Check User (Guest) By Email
	if account == nil {
		user, err := s.userRepo.FindByEmail(email)
		if err != nil {
			return "", "", err
		}
		if user != nil {
			account = user
			role = "guest"
		}
	}

	if account == nil {
		return "", "", errors.New("account not found")
	}

	// Utils : Compare Password
	if err := utils.CheckPassword(account, password); err != nil {
		return "", "", errors.New("invalid password")
	}

	// Utils : Generate Token
	token, err := utils.GenerateToken(account.GetID(), role)
	if err != nil {
		return "", "", err
	}

	return token, role, nil
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
