package repository

import (
	"errors"
	"fmt"
	"pelita/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssetMaintenanceRepository interface {
	FindAll() ([]entity.AssetMaintenance, error)
	Create(assetMaintenance *entity.AssetMaintenance, adminId uuid.UUID) error
	FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time) (*entity.AssetMaintenance, error)
	FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetPlacementId, maintenanceBy uuid.UUID, maintenanceDay string, maintenanceHourStart, maintenanceHourEnd entity.Time, id uuid.UUID) (*entity.AssetMaintenance, error)
	UpdateById(assetMaintenance *entity.AssetMaintenance, id uuid.UUID) error
	DeleteById(id uuid.UUID) error
}

type assetMaintenanceRepository struct {
	db *gorm.DB
}

func NewAssetMaintenanceRepository(db *gorm.DB) AssetMaintenanceRepository {
	return &assetMaintenanceRepository{db: db}
}

func (r *assetMaintenanceRepository) FindAll() ([]entity.AssetMaintenance, error) {
	// Models
	var assetMaintenance []entity.AssetMaintenance

	// Query
	err := r.db.Find(&assetMaintenance).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return assetMaintenance, err
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
