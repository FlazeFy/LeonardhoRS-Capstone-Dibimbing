package service

import (
	"errors"
	"fmt"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type TechnicianService interface {
	GetAllTechnician() ([]entity.Technician, error)
	Create(technician *entity.Technician, adminId uuid.UUID) error
	UpdateById(technician *entity.Technician, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

type technicianService struct {
	technicianRepo repository.TechnicianRepository
}

func NewTechnicianService(technicianRepo repository.TechnicianRepository) TechnicianService {
	return &technicianService{
		technicianRepo: technicianRepo,
	}
}

func (s *technicianService) GetAllTechnician() ([]entity.Technician, error) {
	// Repo : Get All Technician
	technician, err := s.technicianRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if technician == nil {
		return nil, errors.New("technician not found")
	}

	return technician, nil
}

func (s *technicianService) Create(technician *entity.Technician, adminId uuid.UUID) error {
	// Validator
	if technician.Username == "" {
		return errors.New("username is required")
	}
	if technician.Password == "" {
		return errors.New("password is required")
	}
	if technician.Email == "" {
		return errors.New("email is required")
	}

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
	// Validator
	if technician.Username == "" {
		return errors.New("username is required")
	}
	if technician.Password == "" {
		return errors.New("password is required")
	}
	if technician.Email == "" {
		return errors.New("email is required")
	}

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
