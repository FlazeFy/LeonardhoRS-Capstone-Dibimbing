package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type AssetPlacementService interface {
	GetAllAssetPlacement() ([]entity.AssetPlacement, error)
	Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error
	UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

type assetPlacementService struct {
	assetPlacementRepo repository.AssetPlacementRepository
}

func NewAssetPlacementService(assetPlacementRepo repository.AssetPlacementRepository) AssetPlacementService {
	return &assetPlacementService{
		assetPlacementRepo: assetPlacementRepo,
	}
}

func (s *assetPlacementService) GetAllAssetPlacement() ([]entity.AssetPlacement, error) {
	// Repo : Get All AssetPlacement
	assetPlacement, err := s.assetPlacementRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if assetPlacement == nil {
		return nil, errors.New("asset placement not found")
	}

	return assetPlacement, nil
}

func (s *assetPlacementService) Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error {
	// Validator
	if assetPlacement.AssetId == uuid.Nil {
		return errors.New("asset placement name is required")
	}
	if assetPlacement.RoomId == uuid.Nil {
		return errors.New("asset placement category is required")
	}

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
	// Validator
	if assetPlacement.AssetId == uuid.Nil {
		return errors.New("asset placement name is required")
	}
	if assetPlacement.RoomId == uuid.Nil {
		return errors.New("asset placement category is required")
	}

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
