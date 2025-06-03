package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
)

type AdminService interface {
	GetAllContact() ([]entity.AdminContact, error)
}

type adminService struct {
	adminRepo repository.AdminRepository
}

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
