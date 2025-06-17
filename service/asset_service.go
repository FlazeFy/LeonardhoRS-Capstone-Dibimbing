package service

import (
	"errors"
	"mime/multipart"
	"pelita/entity"
	"pelita/repository"
	"pelita/utils"

	"github.com/google/uuid"
)

// Asset Interface
type AssetService interface {
	GetAllAsset(pagination utils.Pagination) ([]entity.Asset, int64, error)
	GetDeleted() ([]entity.Asset, error)
	GetMostContext(targetCol string) ([]entity.StatsContextTotal, error)
	Create(asset *entity.Asset, adminId uuid.UUID, file *multipart.FileHeader, fileExt string, fileSize int64) error
	UpdateById(asset *entity.Asset, id uuid.UUID) error
	HardDeleteById(id uuid.UUID) error
	SoftDeleteById(id uuid.UUID) error
	RecoverDeletedById(id uuid.UUID) error
}

// Asset Struct
type assetService struct {
	assetRepo repository.AssetRepository
	statsRepo repository.StatsRepository
}

// Asset Constructor
func NewAssetService(assetRepo repository.AssetRepository, statsRepo repository.StatsRepository) AssetService {
	return &assetService{
		assetRepo: assetRepo,
		statsRepo: statsRepo,
	}
}

func (s *assetService) GetAllAsset(pagination utils.Pagination) ([]entity.Asset, int64, error) {
	// Repo : Get All Asset
	asset, total, err := s.assetRepo.FindAll(pagination)
	if err != nil {
		return nil, 0, err
	}
	if asset == nil {
		return nil, 0, errors.New("asset not found")
	}

	return asset, total, nil
}

func (s *assetService) GetDeleted() ([]entity.Asset, error) {
	// Repo : Get All Deleted Asset
	asset, err := s.assetRepo.FindDeleted()
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("deleted asset not found")
	}

	return asset, nil
}

func (s *assetService) Create(asset *entity.Asset, adminId uuid.UUID, file *multipart.FileHeader, fileExt string, fileSize int64) error {
	// Repo : Get Asset by Asset Name & Category & Merk
	is_exist, err := s.assetRepo.FindByAssetNameCategoryAndMerk(asset.AssetName, asset.AssetCategory, asset.AssetMerk)
	if err != nil {
		return err
	}
	if is_exist != nil {
		return errors.New("asset already exist on the same floor")
	}

	// Utils : Firebase Upload image
	if file != nil {
		assetImage, err := utils.UploadFile(adminId, "asset", file, fileExt)
		if err != nil {
			return errors.New(err.Error())
		}
		asset.AssetImageURL = &assetImage
	} else {
		asset.AssetImageURL = nil
	}

	// Repo : Create Asset
	if err := s.assetRepo.Create(asset, adminId); err != nil {
		return err
	}

	return nil
}

func (s *assetService) UpdateById(asset *entity.Asset, id uuid.UUID) error {
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

func (s *assetService) GetMostContext(targetCol string) ([]entity.StatsContextTotal, error) {
	// Repo : Get My History
	asset, err := s.statsRepo.FindMostUsedContext("assets", targetCol)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.New("asset not found")
	}

	return asset, nil
}
