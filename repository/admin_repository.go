package repository

import (
	"errors"

	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Admin Interface
type AdminRepository interface {
	FindByEmail(email string) (*entity.Admin, error)
	FindAllContact() ([]entity.AdminContact, error)

	// For Seeder
	Create(room *entity.Admin) error
	DeleteAll() error
}

// Admin Struct
type adminRepository struct {
	db *gorm.DB
}

// Admin Constructor
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) FindAllContact() ([]entity.AdminContact, error) {
	// Models
	var admin []entity.AdminContact

	// Query
	err := r.db.Table("admins").
		Select("username, email, telegram_is_valid, telegram_user_id").
		Where("telegram_is_valid = ?", true).
		Where("telegram_user_id IS NOT NULL").
		Order("username ASC").
		Find(&admin).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return admin, err
}

func (r *adminRepository) FindByEmail(email string) (*entity.Admin, error) {
	// Models
	var admin entity.Admin

	// Query
	err := r.db.Where("email = ?", email).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &admin, err
}

// For Seeder
func (r *adminRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&entity.Admin{}).Error
}
func (r *adminRepository) Create(admin *entity.Admin) error {
	admin.ID = uuid.New()

	// Query
	return r.db.Create(admin).Error
}
