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
	FindByAssetPlacementIdAndRoomId(assetId, roomId uuid.UUID) (*entity.AssetPlacement, error)
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

func (r *assetPlacementRepository) FindByAssetPlacementIdAndRoomId(assetId, roomId uuid.UUID) (*entity.AssetPlacement, error) {
	// Models
	var assetPlacement entity.AssetPlacement

	// Query
	err := r.db.Where("asset_id = ? AND room_id = ?", assetId, roomId).First(&assetPlacement).Error
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
