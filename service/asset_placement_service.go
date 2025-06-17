package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

// Asset Placement Interface
type AssetPlacementService interface {
	GetAllAssetPlacement(pagination utils.Pagination) ([]entity.AssetPlacement, int64, error)
	Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error
	UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

// Asset Placement Struct
type assetPlacementService struct {
	assetPlacementRepo repository.AssetPlacementRepository
}

// Asset Placement Constructor
func NewAssetPlacementService(assetPlacementRepo repository.AssetPlacementRepository) AssetPlacementService {
	return &assetPlacementService{
		assetPlacementRepo: assetPlacementRepo,
	}
}

func (s *assetPlacementService) GetAllAssetPlacement(pagination utils.Pagination) ([]entity.AssetPlacement, int64, error) {
	// Repo : Get All Asset Placement
	assetPlacement, total, err := s.assetPlacementRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if assetPlacement == nil {
		return nil, 0, errors.New("asset placement not found")
	}

	return assetPlacement, total, nil
}

func (s *assetPlacementService) Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error {
	// Repo : Get Asset Placement by Room Id and Asset Id
	is_exist, err := s.assetPlacementRepo.FindByAssetIdAndRoomId(assetPlacement.AssetId, assetPlacement.RoomId)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset is already placed in the room")
	}

	// Repo : Create Asset Placement
	if err := s.assetPlacementRepo.Create(assetPlacement, adminId); err != nil {
		return err
	}

	return nil
}

func (s *assetPlacementService) UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error {
	// Repo : Get Asset by Asset Name & Floor
	is_exist, err := s.assetPlacementRepo.FindByAssetIdRoomIdAndId(assetPlacement.AssetId, assetPlacement.RoomId, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset already exist on the same floor")
	}

	// Repo : Update Asset By Id
	if err := s.assetPlacementRepo.UpdateById(assetPlacement, id); err != nil {
		return err
	}

	return nil
}

func (s *assetPlacementService) DeleteById(id uuid.UUID) error {
	// Repo : Delete Asset Placement By Id
	err := s.assetPlacementRepo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}
