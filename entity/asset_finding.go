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
		FindingByTechnician *uuid.UUID `json:"finding_by_technician" gorm:"null"`
		Technician          Technician `json:"technicians" gorm:"foreignKey:FindingByTechnician;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - User / Guest
		FindingByUser *uuid.UUID `json:"finding_by_user" gorm:"null"`
		User          User       `json:"users" gorm:"foreignKey:FindingByUser;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	AssetFindingReport struct {
		AssetName       string    `json:"asset_name"`
		FindingCategory string    `json:"finding_category"`
		FindingNotes    string    `json:"finding_notes"`
		CreatedAt       time.Time `json:"created_at"`
		// FK - Asset Placement
		Floor    string `json:"floor"`
		RoomName string `json:"room_name"`
		// FK - Technician
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	// For Response Only
	ResponseGetAllAssetFinding struct {
		Message  string         `json:"message" example:"asset finding fetched"`
		Status   string         `json:"status" example:"success"`
		Data     []AssetFinding `json:"data"`
		Metadata Metadata       `json:"metadata"`
	}
	ResponseGetFindingHourTotal struct {
		Message string              `json:"message" example:"asset finding fetched"`
		Status  string              `json:"status" example:"success"`
		Data    []StatsContextTotal `json:"data"`
	}
	ResponseDeleteAssetFindingById struct {
		Message string `json:"message" example:"asset finding deleted"`
		Status  string `json:"status" example:"success"`
	}
	ResponseCreateAssetFinding struct {
		Message string `json:"message" example:"asset finding created"`
		Status  string `json:"status" example:"success"`
	}
)
