package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	History struct {
		ID           uuid.UUID  `gorm:"type:varchar(36);primaryKey"`
		AdminID      *uuid.UUID `gorm:"type:varchar(36);null"`
		TechnicianID *uuid.UUID `gorm:"type:varchar(36);null"`
		UserID       *uuid.UUID `gorm:"type:varchar(36);null"`
		TypeUser     string     `gorm:"type:varchar(36);not null"`
		TypeHistory  string     `gorm:"type:varchar(255);not null"`
		CreatedAt    time.Time  `gorm:"autoCreateTime"`
	}
)
