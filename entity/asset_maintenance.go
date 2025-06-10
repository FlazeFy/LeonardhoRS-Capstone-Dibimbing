package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	AssetMaintenance struct {
		ID                   uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
		MaintenanceDay       string     `json:"maintenance_day" gorm:"type:varchar(3);not null"`
		MaintenanceHourStart Time       `json:"maintenance_hour_start" gorm:"not null"`
		MaintenanceHourEnd   Time       `json:"maintenance_hour_end" gorm:"not null"`
		MaintenanceNotes     *string    `json:"maintenance_notes" gorm:"type:varchar(144);null"`
		CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
		UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime;null"`
		// FK - Asset Placement
		AssetPlacementId uuid.UUID      `json:"asset_placement_id" gorm:"not null"`
		AssetPlacement   AssetPlacement `json:"-" gorm:"foreignKey:AssetPlacementId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Admin
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		Admin     Admin     `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Technician
		MaintenanceBy uuid.UUID  `json:"maintenance_by" gorm:"not null"`
		Technician    Technician `json:"-" gorm:"foreignKey:MaintenanceBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	AssetMaintenanceSchedule struct {
		MaintenanceDay       string  `json:"maintenance_day"`
		MaintenanceHourStart Time    `json:"maintenance_hour_start"`
		MaintenanceHourEnd   Time    `json:"maintenance_hour_end"`
		MaintenanceNotes     *string `json:"maintenance_notes"`
		// FK - Asset Placement
		AssetQty int `json:"asset_qty"`
		// FK - Asset
		AssetName     string `json:"asset_name"`
		AssetCategory string `json:"asset_category"`
		// FK - Technician
		Username        string  `json:"username"`
		Email           string  `json:"email"`
		TelegramUserId  *string `json:"telegram_user_id"`
		TelegramIsValid bool    `json:"telegram_is_valid"`
	}
	// For Response Only
	ResponseGetAllAssetMaintenance struct {
		Message  string             `json:"message" example:"asset maintenance fetched"`
		Status   string             `json:"status" example:"success"`
		Data     []AssetMaintenance `json:"data"`
		Metadata Metadata           `json:"metadata"`
	}
	ResponseGetAllAssetMaintenanceSchedule struct {
		Message string                     `json:"message" example:"asset maintenance schedule fetched"`
		Status  string                     `json:"status" example:"success"`
		Data    []AssetMaintenanceSchedule `json:"data"`
	}
)
