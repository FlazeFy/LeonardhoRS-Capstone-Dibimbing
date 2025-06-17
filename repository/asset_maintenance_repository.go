package repository

import (
	"errors"
	"fmt"
	"pelita/entity"
	"pelita/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Asset Maintenance Interface
type AssetMaintenanceRepository interface {
	FindAll(pagination utils.Pagination) ([]entity.AssetMaintenance, int64, error)
	FindAllSchedule() ([]entity.AssetMaintenanceSchedule, error)
	Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error
	FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time) (*entity.AssetMaintenance, error)
	FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time, id uuid.UUID) (*entity.AssetMaintenance, error)
	UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

// Asset Maintenance Struct
type assetMaintenanceRepository struct {
	db *gorm.DB
}

// Asset Maintenance Constructor
func NewAssetMaintenanceRepository(db *gorm.DB) AssetMaintenanceRepository {
	return &assetMaintenanceRepository{db: db}
}

func (r *assetMaintenanceRepository) FindAll(pagination utils.Pagination) ([]entity.AssetMaintenance, int64, error) {
	var total int64

	// Models
	var assetMaintenance []entity.AssetMaintenance

	// Pagination
	offset := (pagination.Page - 1) * pagination.Limit
	r.db.Model(&entity.AssetMaintenance{}).Count(&total)

	// Query
	err := r.db.Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&assetMaintenance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, err
	}

	return assetMaintenance, total, nil
}

func (r *assetMaintenanceRepository) FindAllSchedule() ([]entity.AssetMaintenanceSchedule, error) {
	// Models
	var asset []entity.AssetMaintenanceSchedule

	// Query
	err := r.db.Table("asset_maintenances").
		Select("maintenance_day, maintenance_hour_start, maintenance_hour_end, maintenance_notes, asset_qty, asset_name, asset_category, username, email, telegram_user_id, telegram_is_valid").
		Joins("JOIN asset_placements ON asset_maintenances.asset_placement_id = asset_placements.id").
		Joins("JOIN assets ON assets.id = asset_placements.asset_id").
		Joins("JOIN technicians ON technicians.id = asset_maintenances.maintenance_by").
		Order("FIELD(maintenance_day, 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'), maintenance_hour_start ASC").
		Find(&asset).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return asset, err
}

func (r *assetMaintenanceRepository) FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time) (*entity.AssetMaintenance, error) {
	// Models
	var existingRecords []entity.AssetMaintenance

	// Query
	err := r.db.Where("asset_placement_id = ? AND maintenance_by = ? AND maintenance_day = ?", assetPlacementId, maintenanceBy, maintenanceDay).
		Find(&existingRecords).Error

	if err != nil {
		return nil, err
	}

	newStart := maintenanceHourStart.Time
	newEnd := maintenanceHourEnd.Time

	// Check For Overlap Time
	for _, record := range existingRecords {
		existingStart := record.MaintenanceHourStart.Time
		existingEnd := record.MaintenanceHourEnd.Time

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return nil, fmt.Errorf("time range overlaps with existing maintenance window (%s - %s)",
				existingStart.Format("15:04"), existingEnd.Format("15:04"))
		}
	}

	return nil, nil
}

func (r *assetMaintenanceRepository) FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time, id uuid.UUID) (*entity.AssetMaintenance, error) {
	// Models
	var existingRecords []entity.AssetMaintenance

	// Query
	err := r.db.Where("asset_placement_id = ? AND maintenance_by = ? AND maintenance_day = ? AND id != ?", assetPlacementId, maintenanceBy, maintenanceDay, id).
		Find(&existingRecords).Error

	if err != nil {
		return nil, err
	}

	newStart := maintenanceHourStart.Time
	newEnd := maintenanceHourEnd.Time

	// Check For Overlap Time
	for _, record := range existingRecords {
		existingStart := record.MaintenanceHourStart.Time
		existingEnd := record.MaintenanceHourEnd.Time

		if newStart.Before(existingEnd) && newEnd.After(existingStart) {
			return nil, fmt.Errorf("time range overlaps with existing maintenance window (%s - %s)",
				existingStart.Format("15:04"), existingEnd.Format("15:04"))
		}
	}

	return nil, nil
}

func (r *assetMaintenanceRepository) Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error {
	now := time.Now()

	assetMaintenance.ID = uuid.New()
	assetMaintenance.CreatedBy = adminId
	assetMaintenance.CreatedAt = now

	// Query
	return r.db.Create(assetMaintenance).Error
}

func (r *assetMaintenanceRepository) UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error {
	now := time.Now()

	// Query : Check Old Asset Maintenance
	var existingAssetMaintenance entity.AssetMaintenance
	if err := r.db.First(&existingAssetMaintenance, "id = ?", id).Error; err != nil {
		return err
	}

	// Query : Update
	existingAssetMaintenance.UpdatedAt = &now
	existingAssetMaintenance.MaintenanceDay = assetMaintenance.MaintenanceDay
	existingAssetMaintenance.MaintenanceHourStart = assetMaintenance.MaintenanceHourStart
	existingAssetMaintenance.MaintenanceHourEnd = assetMaintenance.MaintenanceHourEnd

	if err := r.db.Save(&existingAssetMaintenance).Error; err != nil {
		return err
	}

	return nil
}

func (r *assetMaintenanceRepository) DeleteById(id uuid.UUID) error {
	// Models
	var assetMaintenance entity.AssetMaintenance

	// Query
	err := r.db.Unscoped().Where("id = ?", id).Delete(&assetMaintenance).Error
	if err != nil {
		return err
	}

	return nil
}
