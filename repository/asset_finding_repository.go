package repository

import (
	"errors"
	"pelita/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetFindingRepository interface {
	FindAll() ([]entity.AssetFinding, error)
	Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.NullUUID) error
}

type assetFindingRepository struct {
	db *gorm.DB
}

func NewAssetFindingRepository(db *gorm.DB) AssetFindingRepository {
	return &assetFindingRepository{db: db}
}

func (r *assetFindingRepository) FindAll() ([]entity.AssetFinding, error) {
	// Models
	var assetFinding []entity.AssetFinding

	// Query
	err := r.db.Find(&assetFinding).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return assetFinding, err
}

func (r *assetFindingRepository) Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.NullUUID) error {
	now := time.Now()

	assetFinding.ID = uuid.New()
	assetFinding.FindingByTechnician = technicianId
	assetFinding.FindingByUser = userId
	assetFinding.CreatedAt = now

	// Query
	return r.db.Create(assetFinding).Error
}
