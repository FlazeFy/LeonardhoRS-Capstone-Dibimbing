package repository

import (
	"errors"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset Placement Interface
type AssetPlacementRepository interface {
	FindAll(pagination utils.Pagination) ([]entity.AssetPlacement, int64, error)
	Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error
	FindByAssetIdAndRoomId(assetId, assetPlacementId uuid.UUID) (*entity.AssetPlacement, error)
	FindByAssetIdRoomIdAndId(assetId, assetPlacementId uuid.UUID, id uuid.UUID) (*entity.AssetPlacement, error)
	UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error
	DeleteById(id uuid.UUID) error

	// For Seeder
	DeleteAll() error
	FindOneRandom() (*entity.AssetPlacement, error)
}

// Asset Placement Struct
type assetPlacementRepository struct {
	db *gorm.DB
}

// Asset Placement Constructor
func NewAssetPlacementRepository(db *gorm.DB) AssetPlacementRepository {
	return &assetPlacementRepository{db: db}
}

func (r *assetPlacementRepository) FindAll(pagination utils.Pagination) ([]entity.AssetPlacement, int64, error) {
	var total int64

	// Models
	var assetPlacement []entity.AssetPlacement

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	r.db.Model(&entity.AssetPlacement{}).Count(&total)

	// Query
	err := r.db.Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&assetPlacement).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return assetPlacement, total, nil
}

func (r *assetPlacementRepository) FindByAssetIdAndRoomId(assetId, roomId uuid.UUID) (*entity.AssetPlacement, error) {
	// Models
	var assetPlacement entity.AssetPlacement

	// Query
	err := r.db.Where("asset_id = ? AND room_id = ?", assetId, roomId).First(&assetPlacement).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &assetPlacement, err
}

func (r *assetPlacementRepository) FindByAssetIdRoomIdAndId(assetId, roomId, id uuid.UUID) (*entity.AssetPlacement, error) {
	// Models
	var assetPlacement entity.AssetPlacement

	// Query
	err := r.db.Where("asset_id = ? AND room_id = ? AND id != ?", assetId, roomId, id).First(&assetPlacement).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &assetPlacement, err
}

func (r *assetPlacementRepository) Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error {
	now := time.Now()

	assetPlacement.ID = uuid.New()
	assetPlacement.CreatedBy = adminId
	assetPlacement.CreatedAt = now
	assetPlacement.UpdatedAt = nil

	// Query
	return r.db.Create(assetPlacement).Error
}

func (r *assetPlacementRepository) UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error {
	now := time.Now()

	// Query : Check Old Asset Placement
	var existingAssetPlacement entity.AssetPlacement
	if err := r.db.First(&existingAssetPlacement, "id = ?", id).Error; err != nil {
		return err
	}

	// Query : Update
	existingAssetPlacement.UpdatedAt = &now
	existingAssetPlacement.AssetQty = assetPlacement.AssetQty

	if err := r.db.Save(&existingAssetPlacement).Error; err != nil {
		return err
	}

	return nil
}

func (r *assetPlacementRepository) DeleteById(id uuid.UUID) error {
	// Models
	var assetPlacement entity.AssetPlacement

	// Query
	err := r.db.Unscoped().Where("id = ?", id).Delete(&assetPlacement).Error
	if err != nil {
		return err
	}

	return nil
}

// For Seeder
func (r *assetPlacementRepository) DeleteAll() error {
	return r.db.Where("1 = 1").Delete(&entity.AssetPlacement{}).Error
}
func (r *assetPlacementRepository) FindOneRandom() (*entity.AssetPlacement, error) {
	var assetPlacement entity.AssetPlacement

	err := r.db.Order("RAND()").Limit(1).First(&assetPlacement).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &assetPlacement, err
}
