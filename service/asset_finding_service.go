package service

import (
	"errors"
	"mime/multipart"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

type AssetFindingService interface {
	GetAllAssetFinding(pagination utils.Pagination) ([]entity.AssetFinding, int64, error)
	GetMostContext(targetCol string) ([]entity.StatsContextTotal, error)
	GetFindingHourTotal() ([]entity.StatsContextTotal, error)
	Create(assetFinding *entity.AssetFinding, technicianId, userId *uuid.UUID, file *multipart.FileHeader, fileExt string, fileSize int64) error
	DeleteById(id uuid.UUID) error

	// Scheduler Service
	GetAllAssetFindingReport() ([]entity.AssetFindingReport, error)
}

type assetFindingService struct {
	assetFindingRepo repository.AssetFindingRepository
	statsRepo        repository.StatsRepository
}

func NewAssetFindingService(assetFindingRepo repository.AssetFindingRepository, statsRepo repository.StatsRepository) AssetFindingService {
	return &assetFindingService{
		assetFindingRepo: assetFindingRepo,
		statsRepo:        statsRepo,
	}
}

func (s *assetFindingService) GetAllAssetFinding(pagination utils.Pagination) ([]entity.AssetFinding, int64, error) {
	// Repo : Get All Asset Finding
	assetFinding, total, err := s.assetFindingRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if assetFinding == nil {
		return nil, 0, errors.New("asset finding not found")
	}

	return assetFinding, total, nil
}

func (s *assetFindingService) GetAllAssetFindingReport() ([]entity.AssetFindingReport, error) {
	// Repo : Get All Asset Finding
	assetFinding, err := s.assetFindingRepo.FindAllReport()
	if err != nil {
		return nil, err
	}
	if assetFinding == nil {
		return nil, errors.New("asset finding not found")
	}

	return assetFinding, nil
}

func (s *assetFindingService) Create(assetFinding *entity.AssetFinding, technicianId, userId *uuid.UUID, file *multipart.FileHeader, fileExt string, fileSize int64) error {
	// Validator
	if assetFinding.AssetPlacementId == uuid.Nil {
		return errors.New("asset placement id is required")
	}
	if technicianId == nil && userId == nil {
		return errors.New("technician id and user id is required")
	}

	// Utils : Firebase Upload image
	if file != nil {
		var createdBy uuid.UUID
		if technicianId != nil {
			createdBy = *technicianId
		} else if userId != nil {
			createdBy = *userId
		}
		assetImage, err := utils.UploadFile(createdBy, "asset", file, fileExt)

		if err != nil {
			return errors.New(err.Error())
		}
		assetFinding.FindingImage = &assetImage
	} else {
		assetFinding.FindingImage = nil
	}

	// Repo : Create Asset Finding
	if err := s.assetFindingRepo.Create(assetFinding, *technicianId, *userId); err != nil {
		return err
	}

	return nil
}

func (s *assetFindingService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Asset Finding By Id
	err := s.assetFindingRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *assetFindingService) GetMostContext(targetCol string) ([]entity.StatsContextTotal, error) {
	// Repo : Get Most Context
	asset, err := s.statsRepo.FindMostUsedContext("asset_findings", targetCol)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset finding not found")
	}

	return asset, nil
}

func (s *assetFindingService) GetFindingHourTotal() ([]entity.StatsContextTotal, error) {
	// Repo : Get Finding Hour Total
	asset, err := s.assetFindingRepo.FindAllFindingHourTotal()
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset finding not found")
	}

	return asset, nil
}
