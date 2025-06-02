package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	AssetFinding struct {
		ID              uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
		FindingCategory string    `json:"finding_category" gorm:"type:varchar(36);not null"`
		FindingNotes    string    `json:"finding_notes" gorm:"type:varchar(255);not null"`
		FindingImage    *string   `json:"finding_image" gorm:"type:varchar(500);null"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:datetime;not null"`
		// FK - Asset Placement
		AssetPlacementId uuid.UUID      `json:"asset_placement_id" gorm:"not null"`
		AssetPlacement   AssetPlacement `json:"asset_placements" gorm:"foreignKey:AssetPlacementId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Technician
		FindingByTechnician uuid.NullUUID `json:"finding_by_technician" gorm:"null"`
		Technician          Technician    `json:"technicians" gorm:"foreignKey:FindingByTechnician;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - User / Guest
		FindingByUser uuid.NullUUID `json:"finding_by_user" gorm:"null"`
		User          User          `json:"users" gorm:"foreignKey:FindingByUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
