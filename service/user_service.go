package service

import (
	"errors"

	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type UserService interface {
	GetMyProfile(user uuid.UUID) (*entity.MyProfile, error)
}

type userService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
}

func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		redisClient: redisClient,
	}
}

func (s *userService) GetMyProfile(userId uuid.UUID) (*entity.MyProfile, error) {
	user, err := s.userRepo.FindById(userId.String())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
