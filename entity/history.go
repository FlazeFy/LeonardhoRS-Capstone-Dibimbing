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
	AllHistory struct {
		ID          uuid.UUID `gorm:"type:varchar(36);primaryKey"`
		TypeUser    string    `gorm:"type:varchar(36);not null"`
		TypeHistory string    `gorm:"type:varchar(255);not null"`
		CreatedAt   time.Time `gorm:"autoCreateTime"`
		// FK - Admin
		AdminID *uuid.UUID `json:"admin_id"`
		Admin   *Admin     `json:"admins" gorm:"foreignKey:AdminID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - User
		UserID *uuid.UUID `json:"user_id"`
		User   *User      `json:"users" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
		// FK - Technician
		TechnicianID *uuid.UUID  `json:"technician_id"`
		Technician   *Technician `json:"technicians" gorm:"foreignKey:TechnicianID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
)
