package repository

import (
	"errors"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetFindingRepository interface {
	FindAll(pagination utils.Pagination) ([]entity.AssetFinding, int64, error)
	FindAllReport() ([]entity.AssetFindingReport, error)
	FindAllFindingHourTotal() ([]entity.StatsContextTotal, error)
	Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

type assetFindingRepository struct {
	db *gorm.DB
}

func NewAssetFindingRepository(db *gorm.DB) AssetFindingRepository {
	return &assetFindingRepository{db: db}
}

func (r *assetFindingRepository) FindAll(pagination utils.Pagination) ([]entity.AssetFinding, int64, error) {
	var total int64

	// Models
	var assetFinding []entity.AssetFinding

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	r.db.Model(&entity.AssetFinding{}).Count(&total)

	// Query
	err := r.db.Preload("User").
		Preload("Technician").
		Preload("AssetPlacement").
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&assetFinding).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return assetFinding, total, nil
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

func (r *assetFindingRepository) FindAllFindingHourTotal() ([]entity.StatsContextTotal, error) {
	// Models
	var asset []entity.StatsContextTotal

	// Query
	err := r.db.Table("asset_maintenances").
		Select("HOUR(created_at) as context, COUNT(1) as total").
		Group("HOUR(created_at)").
		Order("total DESC").
		Find(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return asset, err
}

func (r *assetFindingRepository) Create(assetFinding *entity.AssetFinding, technicianId, userId uuid.UUID) error {
	now := time.Now()

	assetFinding.ID = uuid.New()
	assetFinding.FindingByTechnician = &technicianId
	assetFinding.FindingByUser = &userId
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
