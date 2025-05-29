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
	UpdateById(asset *entity.Asset, id uuid.UUID) error
	HardDeleteById(id uuid.UUID) error
	SoftDeleteById(id uuid.UUID) error
	RecoverDeletedById(id uuid.UUID) error
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

func (s *assetService) UpdateById(asset *entity.Asset, id uuid.UUID) error {
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

	// Repo : Get Asset by Asset Name & Floor
	is_exist, err := s.assetRepo.FindByAssetNameCategoryMerkAndId(asset.AssetName, asset.AssetCategory, asset.AssetMerk, id)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset already exist on the same floor")
	}

	// Repo : Update Asset By Id
	if err := s.assetRepo.UpdateById(asset, id); err != nil {
		return err
	}

	return nil
}

func (s *assetService) HardDeleteById(id uuid.UUID) error {
	// Repo : Delete Asset By Id
	err := s.assetRepo.HardDeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *assetService) SoftDeleteById(id uuid.UUID) error {
	// Repo : Delete Asset By Id
	err := s.assetRepo.SoftDeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *assetService) RecoverDeletedById(id uuid.UUID) error {
	// Repo : Recover Asset By Id
	err := s.assetRepo.RecoverDeletedById(id)
	if err != nil {
		return err
	}

	return nil
}
