package repository

import (
	"errors"
	"pelita/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetPlacementRepository interface {
	FindAll() ([]entity.AssetPlacement, error)
	Create(assetPlacement *entity.AssetPlacement, adminId uuid.UUID) error
	FindByAssetIdAndRoomId(assetId, assetPlacementId uuid.UUID) (*entity.AssetPlacement, error)
	FindByAssetIdRoomIdAndId(assetId, assetPlacementId uuid.UUID, id uuid.UUID) (*entity.AssetPlacement, error)
	UpdateById(assetPlacement *entity.AssetPlacement, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

type assetPlacementRepository struct {
	db *gorm.DB
}

func NewAssetPlacementRepository(db *gorm.DB) AssetPlacementRepository {
	return &assetPlacementRepository{db: db}
}

func (r *assetPlacementRepository) FindAll() ([]entity.AssetPlacement, error) {
	// Models
	var assetPlacement []entity.AssetPlacement

	// Query
	err := r.db.Find(&assetPlacement).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return assetPlacement, err
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
