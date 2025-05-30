package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Asset struct {
		ID            uuid.UUID  `json:"id" gorm:"type:varchar(36);primaryKey"`
		AssetName     string     `json:"asset_name" gorm:"type:varchar(144);not null"`
		AssetDesc     *string    `json:"asset_desc" gorm:"type:varchar(500);null"`
		AssetMerk     *string    `json:"asset_merk" gorm:"type:varchar(75);null"`
		AssetCategory string     `json:"asset_category" gorm:"type:varchar(36);not null"`
		AssetPrice    *string    `json:"asset_price" gorm:"type:varchar(9);null"`
		AssetStatus   string     `json:"asset_status" gorm:"type:varchar(36);not null"`
		AssetImageURL *string    `json:"asset_image_url" gorm:"type:varchar(1000);null"`
		CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
		UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
		DeletedAt     *time.Time `json:"deleted_at" gorm:"type:datetime;default:null"`
		// FK - Admin
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		Admin     Admin     `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
