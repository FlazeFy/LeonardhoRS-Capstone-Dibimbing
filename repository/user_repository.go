package repository

import (
	"errors"
	"fmt"
	"time"

	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Interface
type UserRepository interface {
	FindByUsernameOrEmail(username, email string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindById(id, role string) (*entity.MyProfile, error)
	Create(user *entity.User) error

	// For Seeder
	DeleteAll() error
}

// User Struct
type userRepository struct {
	db *gorm.DB
}

// User Constructor
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

func (r *userRepository) FindById(id, role string) (*entity.MyProfile, error) {
	// Models
	var user entity.MyProfile
	var tableName = fmt.Sprintf("%ss", role)
	if role == "guest" {
		tableName = "users"
	}

	// Query
	err := r.db.Table(tableName).
		Select("username, email, telegram_is_valid, telegram_user_id, created_at").
		Where("id = ?", id).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &user, err
}

func (r *userRepository) Create(user *entity.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.TelegramIsValid = false
	user.ID = uuid.New()

	// Query
	return r.db.Create(user).Error
}

// For Seeder
func (r *userRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&entity.User{}).Error
}
