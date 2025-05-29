package repository

import (
	"errors"
	"pelita/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetRepository interface {
	FindAll() ([]entity.Asset, error)
	Create(asset *entity.Asset, adminId uuid.UUID) error
	FindByAssetNameCategoryAndMerk(assetName, assetCategory string, assetMerk *string) (*entity.Asset, error)
	FindByAssetNameCategoryMerkAndId(assetName, assetCategory string, assetMerk *string, id uuid.UUID) (*entity.Asset, error)
	UpdateById(asset *entity.Asset, id uuid.UUID) error
	HardDeleteById(id uuid.UUID) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) FindAll() ([]entity.Asset, error) {
	// Models
	var asset []entity.Asset

	// Query
	err := r.db.Find(&asset).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return asset, err
}

func (r *assetRepository) FindByAssetNameCategoryAndMerk(assetName, assetCategory string, assetMerk *string) (*entity.Asset, error) {
	// Models
	var asset entity.Asset

	// Query
	err := r.db.Where("asset_name = ? AND asset_category = ? AND asset_merk = ?", assetName, assetCategory, assetMerk).First(&asset).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &asset, err
}

func (r *assetRepository) FindByAssetNameCategoryMerkAndId(assetName, assetCategory string, assetMerk *string, id uuid.UUID) (*entity.Asset, error) {
	// Models
	var asset entity.Asset

	// Query
	err := r.db.Where("asset_name = ? AND asset_category = ? AND asset_merk = ? AND id != ?", assetName, assetCategory, assetMerk, id).First(&asset).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &asset, err
}

func (r *assetRepository) Create(asset *entity.Asset, adminId uuid.UUID) error {
	now := time.Now()

	asset.ID = uuid.New()
	asset.CreatedBy = adminId
	asset.CreatedAt = now
	asset.UpdatedAt = nil
	asset.DeletedAt = nil

	// Query
	return r.db.Create(asset).Error
}

func (r *assetRepository) UpdateById(asset *entity.Asset, id uuid.UUID) error {
	// Query : Check Old Asset
	var existingAsset entity.Asset
	if err := r.db.First(&existingAsset, "id = ?", id).Error; err != nil {
		return err
	}
	now := time.Now()

	// Query : Update
	asset.ID = id
	asset.CreatedAt = existingAsset.CreatedAt
	asset.UpdatedAt = &now

	if err := r.db.Save(&asset).Error; err != nil {
		return err
	}

	return nil
}

func (r *assetRepository) HardDeleteById(id uuid.UUID) error {
	// Models
	var asset entity.Asset

	// Query
	err := r.db.Unscoped().Where("id = ? AND deleted_at is not null", id).Delete(&asset).Error
	if err != nil {
		return err
	}

	return nil
}
