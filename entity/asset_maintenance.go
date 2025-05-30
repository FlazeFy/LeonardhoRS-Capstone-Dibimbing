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
)
