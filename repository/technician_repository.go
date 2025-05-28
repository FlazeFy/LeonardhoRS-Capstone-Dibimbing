package repository

import (
	"errors"

	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TechnicianRepository interface {
	FindByEmail(email string) (*entity.Technician, error)
	FindByEmailAndId(email string, id uuid.UUID) (*entity.Technician, error)
	FindAll() ([]entity.Technician, error)
	Create(technician *entity.Technician, adminId uuid.UUID) error
	DeleteById(id uuid.UUID) error
	UpdateById(technician *entity.Technician, adminId uuid.UUID) error
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

func (r *technicianRepository) FindByEmailAndId(email string, id uuid.UUID) (*entity.Technician, error) {
	// Models
	var technician entity.Technician

	// Query
	err := r.db.Where("email = ? AND id != ?", email, id).First(&technician).Error
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

func (r *technicianRepository) UpdateById(technician *entity.Technician, id uuid.UUID) error {
	// Query : Check Old Technician
	var existingTechnician entity.Technician
	if err := r.db.First(&existingTechnician, "id = ?", id).Error; err != nil {
		return err
	}

	// Query : Update
	technician.ID = id
	technician.CreatedAt = existingTechnician.CreatedAt
	technician.CreatedBy = existingTechnician.CreatedBy

	if err := r.db.Save(&technician).Error; err != nil {
		return err
	}

	return nil
}

func (r *technicianRepository) DeleteById(id uuid.UUID) error {
	// Models
	var technician entity.Technician

	// Query
	err := r.db.Unscoped().Where("id = ?", id).Delete(&technician).Error
	if err != nil {
		return err
	}

	return nil
}
