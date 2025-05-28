package repository

import (
	"errors"

	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TechnicianRepository interface {
	FindByEmail(email string) (*entity.Technician, error)
	FindAll() ([]entity.Technician, error)
	Create(technician *entity.Technician, adminId uuid.UUID) error
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

func (r *technicianRepository) FindAll() ([]entity.Technician, error) {
	// Models
	var technician []entity.Technician

	// Query
	err := r.db.Find(&technician).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return technician, err
}

func (r *technicianRepository) Create(technician *entity.Technician, adminId uuid.UUID) error {
	technician.ID = uuid.New()
	technician.CreatedBy = adminId

	// Query
	return r.db.Create(technician).Error
}
