package repository

import (
	"errors"
	"pelita/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetRepository interface {
	FindAll() ([]entity.Asset, error)
	Create(asset *entity.Asset, adminId uuid.UUID) error
	FindByAssetNameCategoryAndMerk(assetName, assetCategory string, assetMerk *string) (*entity.Asset, error)
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

func (r *assetRepository) Create(asset *entity.Asset, adminId uuid.UUID) error {
	asset.ID = uuid.New()
	asset.CreatedBy = adminId

	// Query
	return r.db.Create(asset).Error
}
