package service

import (
	"errors"
	"pelita/entity"
	"pelita/repository"

	"github.com/google/uuid"
)

type AssetService interface {
	GetAllAsset() ([]entity.Asset, error)
	Create(asset *entity.Asset, adminId uuid.UUID) error
}

type assetService struct {
	assetRepo repository.AssetRepository
}

func NewAssetService(assetRepo repository.AssetRepository) AssetService {
	return &assetService{
		assetRepo: assetRepo,
	}
}

func (s *assetService) GetAllAsset() ([]entity.Asset, error) {
	// Repo : Get All Asset
	asset, err := s.assetRepo.FindAll()
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}

	return asset, nil
}

func (s *assetService) Create(asset *entity.Asset, adminId uuid.UUID) error {
	// Validator
	if asset.AssetName == "" {
		return errors.New("asset name is required")
	}
	if asset.AssetCategory == "" {
		return errors.New("asset category is required")
	}
	if asset.AssetStatus == "" {
		return errors.New("asset status is required")
	}

	// Repo : Get Asset by Asset Name & Category & Merk
	is_exist, err := s.assetRepo.FindByAssetNameCategoryAndMerk(asset.AssetName, asset.AssetCategory, asset.AssetMerk)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset already exist on the same floor")
	}

	// Repo : Create Asset
	if err := s.assetRepo.Create(asset, adminId); err != nil {
		return err
	}

	return nil
}
