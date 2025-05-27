package repository

import (
	"errors"

	"pelita/entity"

	"gorm.io/gorm"
)

type TechnicianRepository interface {
	FindByEmail(email string) (*entity.Technician, error)
}

type technicianRepository struct {
	db *gorm.DB
}

func NewTechnicianRepository(db *gorm.DB) TechnicianRepository {
	return &technicianRepository{db: db}
}

func (r *technicianRepository) FindByEmail(email string) (*entity.Technician, error) {
	// Models
	var technician entity.Technician

	// Query
	err := r.db.Where("email = ?", email).First(&technician).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &technician, err
}
