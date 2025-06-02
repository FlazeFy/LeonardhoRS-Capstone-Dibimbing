package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Room struct {
		ID        uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
		Floor     string    `json:"floor" gorm:"type:varchar(2);not null"`
		RoomName  string    `json:"room_name" gorm:"type:varchar(36);not null"`
		RoomDept  string    `json:"room_dept" gorm:"type:varchar(75);not null"`
		CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	RoomAsset struct {
		Floor         string  `json:"floor"`
		RoomName      string  `json:"room_name"`
		RoomDept      string  `json:"room_dept"`
		AssetName     string  `json:"asset_name"`
		AssetDesc     *string `json:"asset_desc"`
		TotalAsset    int     `json:"total_asset"`
		AssetMerk     *string `json:"asset_merk"`
		AssetCategory string  `json:"asset_category"`
	}
)
