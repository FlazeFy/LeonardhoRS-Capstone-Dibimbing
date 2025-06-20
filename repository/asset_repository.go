package repository

import (
	"errors"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset Interface
type AssetRepository interface {
	FindAll(pagination utils.Pagination) ([]entity.Asset, int64, error)
	Create(asset *entity.Asset, adminId uuid.UUID) error
	FindByAssetPlacementId(id uuid.UUID) (*entity.Asset, error)
	FindByAssetNameCategoryAndMerk(assetName, assetCategory string, assetMerk *string) (*entity.Asset, error)
	FindByAssetNameCategoryMerkAndId(assetName, assetCategory string, assetMerk *string, id uuid.UUID) (*entity.Asset, error)
	FindDeleted() ([]entity.Asset, error)
	UpdateById(asset *entity.Asset, id uuid.UUID) error
	HardDeleteById(id uuid.UUID) error
	SoftDeleteById(id uuid.UUID) error
	RecoverDeletedById(id uuid.UUID) error

	// For Seeder
	DeleteAll() error
	FindOneRandom() (*entity.Asset, error)
}

// Asset Struct
type assetRepository struct {
	db *gorm.DB
}

// Asset Constructor
func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) FindAll(pagination utils.Pagination) ([]entity.Asset, int64, error) {
	var total int64

	// Models
	var asset []entity.Asset

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	r.db.Model(&entity.Asset{}).Where("deleted_at is null").Count(&total)

	// Query
	err := r.db.Order("created_at DESC").
		Where("deleted_at is null").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return asset, total, nil
}

func (r *assetRepository) FindByAssetPlacementId(id uuid.UUID) (*entity.Asset, error) {
	// Models
	var asset entity.Asset

	// Query
	err := r.db.Table("assets").
		Select("assets.id, asset_name, assets.asset_desc, asset_merk, asset_category, asset_price, asset_status, asset_image_url, assets.created_at, assets.updated_at, deleted_at, assets.created_by").
		Joins("JOIN asset_placements ON asset_placements.asset_id = assets.id").
		Where("asset_placements.id = ?", id).
		First(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &asset, err
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

func (r *assetRepository) FindDeleted() ([]entity.Asset, error) {
	// Models
	var asset []entity.Asset

	// Query
	result := r.db.Order("deleted_at DESC").
		Where("deleted_at is not null").
		Find(&asset)

	// Response
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || len(asset) == 0 {
		return nil, errors.New("deleted asset not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return asset, nil
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

func (r *assetRepository) SoftDeleteById(id uuid.UUID) error {
	// Query : Check Old Asset
	var existingAsset entity.Asset
	if err := r.db.First(&existingAsset, "id = ? AND deleted_at is null", id).Error; err != nil {
		return err
	}
	now := time.Now()

	// Query : Update
	existingAsset.DeletedAt = &now

	if err := r.db.Save(&existingAsset).Error; err != nil {
		return err
	}

	return nil
}

func (r *assetRepository) RecoverDeletedById(id uuid.UUID) error {
	// Query : Check Old Asset
	var existingAsset entity.Asset
	if err := r.db.First(&existingAsset, "id = ? AND deleted_at is not null", id).Error; err != nil {
		return err
	}

	// Query : Update
	existingAsset.DeletedAt = nil

	if err := r.db.Save(&existingAsset).Error; err != nil {
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

// For Seeder
func (r *assetRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&entity.Asset{}).Error
}
func (r *assetRepository) FindOneRandom() (*entity.Asset, error) {
	var asset entity.Asset

	err := r.db.Order("RAND()").Limit(1).First(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &asset, err
}
