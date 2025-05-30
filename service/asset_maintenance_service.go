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
	UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
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

	// Repo : Create Asset Maintenance
	if err := s.assetMaintenanceRepo.Create(assetMaintenance, adminId); err != nil {
		return err
	}

	return nil
}

func (s *assetMaintenanceService) UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error {
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

	// Repo : Get Asset Maintenance By Asset Placement Id, Maintenance By, Maintenance Day, and Id
	is_exist, err := s.assetMaintenanceRepo.FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetMaintenance.AssetPlacementId, assetMaintenance.MaintenanceBy, assetMaintenance.MaintenanceDay, assetMaintenance.MaintenanceHourStart, assetMaintenance.MaintenanceHourEnd, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset is already assigned to maintenance by a technician")
	}

	// Repo : Update Asset Maintenance By Id
	if err := s.assetMaintenanceRepo.UpdateById(assetMaintenance, id); err != nil {
		return err
	}

	return nil
}

func (s *assetMaintenanceService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Asset Maintenance By Id
	err := s.assetMaintenanceRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}
