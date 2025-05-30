package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	AssetPlacement struct {
		ID        uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
		AssetQty  int        `json:"asset_qty" gorm:"type:int;not null"`
		AssetDesc *string    `json:"asset_desc" gorm:"type:varchar(144)"`
		CreatedAt time.Time  `json:"created_at" gorm:"type:datetime;not null"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime"`
		// FK - Asset
		AssetId uuid.UUID `json:"asset_id" gorm:"not null"`
		Asset   Asset     `json:"-" gorm:"foreignKey:AssetId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Room
		RoomId uuid.UUID `json:"room_id" gorm:"not null"`
		Room   Room      `json:"-" gorm:"foreignKey:RoomId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Admin
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		Admin     Admin     `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Technician
		AssetOwner uuid.UUID  `json:"asset_owner" gorm:"not null"`
		Technician Technician `json:"-" gorm:"foreignKey:AssetOwner;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
