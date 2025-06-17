package service

import (
	"errors"
	"fmt"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

// Technician Interface
type TechnicianService interface {
	GetAllTechnician(pagination utils.Pagination) ([]entity.Technician, int64, error)
	Create(technician *entity.Technician, adminId uuid.UUID) error
	UpdateById(technician *entity.Technician, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

// Technician Struct
type technicianService struct {
	technicianRepo repository.TechnicianRepository
}

// Technician Constructor
func NewTechnicianService(technicianRepo repository.TechnicianRepository) TechnicianService {
	return &technicianService{
		technicianRepo: technicianRepo,
	}
}

func (s *technicianService) GetAllTechnician(pagination utils.Pagination) ([]entity.Technician, int64, error) {
	// Repo : Get All Technician
	technician, total, err := s.technicianRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if technician == nil {
		return nil, 0, errors.New("technician not found")
	}

	return technician, total, nil
}

func (s *technicianService) Create(technician *entity.Technician, adminId uuid.UUID) error {
	// Repo : Get Technician by email
	is_exist, err := s.technicianRepo.FindByEmail(technician.Email)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("email already used")
	}

	// Repo : Create Technician
	fmt.Println(adminId)
	if err := s.technicianRepo.Create(technician, adminId); err != nil {
		return err
	}

	return nil
}

func (s *technicianService) UpdateById(technician *entity.Technician, id uuid.UUID) error {
	// Repo : Get Technician by email
	is_exist, err := s.technicianRepo.FindByEmailAndId(technician.Email, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("email already used")
	}

	// Repo : Update Technician By Id
	if err := s.technicianRepo.UpdateById(technician, id); err != nil {
		return err
	}

	return nil
}

func (s *technicianService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Technician By Id
	err := s.technicianRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}
