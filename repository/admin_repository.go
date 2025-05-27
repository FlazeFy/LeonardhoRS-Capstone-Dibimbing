package repository

import (
	"errors"

	"pelita/entity"

	"gorm.io/gorm"
)

type AdminRepository interface {
	FindByEmail(email string) (*entity.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
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
