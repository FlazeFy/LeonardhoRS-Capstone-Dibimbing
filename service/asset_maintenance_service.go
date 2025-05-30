package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type AssetMaintenanceService interface {
	GetAllAssetMaintenance() ([]entity.AssetMaintenance, error)
	Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error
}

type assetMaintenanceService struct {
	assetMaintenanceRepo repository.AssetMaintenanceRepository
}

func NewAssetMaintenanceService(assetMaintenanceRepo repository.AssetMaintenanceRepository) AssetMaintenanceService {
	return &assetMaintenanceService{
		assetMaintenanceRepo: assetMaintenanceRepo,
	}
}

func (s *assetMaintenanceService) GetAllAssetMaintenance() ([]entity.AssetMaintenance, error) {
	// Repo : Get All Asset Maintenance
	assetMaintenance, err := s.assetMaintenanceRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if assetMaintenance == nil {
		return nil, errors.New("asset placement not found")
	}

	return assetMaintenance, nil
}

func (s *assetMaintenanceService) Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error {
	// Validator
	if assetMaintenance.AssetPlacementId == uuid.Nil {
		return errors.New("asset placement id is required")
	}
	if assetMaintenance.MaintenanceBy == uuid.Nil {
		return errors.New("asset maintenance by is required")
	}
	if assetMaintenance.MaintenanceDay == "" {
		return errors.New("asset maintenance day is required")
	}

	// Repo : Get Asset Maintenance By Asset Placement Id, Maintenance By, and Maintenance Day
	is_exist, err := s.assetMaintenanceRepo.FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetMaintenance.AssetPlacementId, assetMaintenance.MaintenanceBy, assetMaintenance.MaintenanceDay, assetMaintenance.MaintenanceHourStart, assetMaintenance.MaintenanceHourEnd)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset is already assigned to maintenance by a technician")
	}

	// Repo : Create Asset Placement
	if err := s.assetMaintenanceRepo.Create(assetMaintenance, adminId); err != nil {
		return err
	}

	return nil
}
