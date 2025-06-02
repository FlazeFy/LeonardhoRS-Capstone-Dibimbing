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
	DeleteById(id uuid.UUID) error
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
	err := r.db.Preload("User").
		Preload("Technician").
		Preload("AssetPlacement").
		Order("created_at DESC").
		Find(&assetFinding).Error

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

func (r *assetFindingRepository) DeleteById(id uuid.UUID) error {
	// Models
	var assetFinding entity.AssetFinding

	// Query
	err := r.db.Unscoped().Where("id = ?", id).Delete(&assetFinding).Error
	if err != nil {
		return err
	}

	return nil
}
