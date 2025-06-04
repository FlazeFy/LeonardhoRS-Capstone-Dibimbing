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
	FindAllReport() ([]entity.AssetFindingReport, error)
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

func (r *assetFindingRepository) FindAllReport() ([]entity.AssetFindingReport, error) {
	// Models
	var assetFinding []entity.AssetFindingReport

	// Query
	err := r.db.Table("asset_findings").
		Select("asset_name,finding_category, finding_notes, asset_findings.created_at, floor, room_name, username, email").
		Joins("JOIN asset_placements ON asset_findings.asset_placement_id = asset_placements.id").
		Joins("JOIN assets ON asset_placements.asset_id = assets.id").
		Joins("JOIN rooms ON rooms.id = asset_placements.room_id").
		Joins("JOIN asset_maintenances ON asset_maintenances.asset_placement_id = asset_placements.id").
		Joins("JOIN technicians ON technicians.id = asset_maintenances.maintenance_by").
		Order("asset_findings.created_at DESC").
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
