package repository

import (
	"errors"

	"pelita/entity"

	"gorm.io/gorm"
)

// Admin Interface
type AdminRepository interface {
	FindByEmail(email string) (*entity.Admin, error)
	FindAllContact() ([]entity.AdminContact, error)
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
