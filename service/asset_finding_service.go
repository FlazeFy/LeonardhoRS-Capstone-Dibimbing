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
	GetAllAssetFinding() ([]entity.AssetFinding, error)
	Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.NullUUID, file *multipart.FileHeader, fileExt string, fileSize int64) error
	DeleteById(id uuid.UUID) error

	// Scheduler Service
	GetAllAssetFindingReport() ([]entity.AssetFindingReport, error)
}

type assetFindingService struct {
	assetFindingRepo repository.AssetFindingRepository
}

func NewAssetFindingService(assetFindingRepo repository.AssetFindingRepository) AssetFindingService {
	return &assetFindingService{
		assetFindingRepo: assetFindingRepo,
	}
}

func (s *assetFindingService) GetAllAssetFinding() ([]entity.AssetFinding, error) {
	// Repo : Get All Asset Finding
	assetFinding, err := s.assetFindingRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if assetFinding == nil {
		return nil, errors.New("asset finding not found")
	}

	return assetFinding, nil
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

func (s *assetFindingService) Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.NullUUID, file *multipart.FileHeader, fileExt string, fileSize int64) error {
	// Validator
	if assetFinding.AssetPlacementId == uuid.Nil {
		return errors.New("asset placement id is required")
	}
	if !technicianId.Valid && !userId.Valid {
		return errors.New("technician id or user id is required")
	}

	// Utils : Firebase Upload image
	if file != nil {
		var createdBy uuid.UUID
		if technicianId.Valid {
			createdBy = technicianId.UUID
		} else if userId.Valid {
			createdBy = userId.UUID
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
	if err := s.assetFindingRepo.Create(assetFinding, technicianId, userId); err != nil {
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
