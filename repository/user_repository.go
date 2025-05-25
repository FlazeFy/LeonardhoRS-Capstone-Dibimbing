package repository

import (
	"errors"

	"pelita/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsernameOrEmail(username, email string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	// Models
	var user entity.User

	// Query
	err := r.db.Where("username = ? OR email = ?", username, email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	// Models
	var user entity.User

	// Query
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) Create(user *entity.User) error {
	// Query
	return r.db.Create(user).Error
}
