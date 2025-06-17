package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
)

// Admin Interface
type AdminService interface {
	GetAllContact() ([]entity.AdminContact, error)
}

// Admin Struct
type adminService struct {
	adminRepo repository.AdminRepository
}

// Admin Constructor
func NewAdminService(adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) GetAllContact() ([]entity.AdminContact, error) {
	// Repo : Get All Asset Placement
	assetPlacement, err := s.adminRepo.FindAllContact()
	if err != nil {
		return nil, err
	}
	if assetPlacement == nil {
		return nil, errors.New("admin not found")
	}

	return assetPlacement, nil
}
