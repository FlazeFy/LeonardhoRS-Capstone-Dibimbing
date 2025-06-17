package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Technician struct {
		ID              uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Password        string    `json:"password" gorm:"type:varchar(500);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
		// FK - User
		CreatedBy uuid.UUID `json:"created_by" gorm:"not null"`
		Admin     Admin     `json:"-" gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	}
	// For Response Only
	ResponseGetAllTechnician struct {
		Message  string       `json:"message" example:"technician fetched"`
		Status   string       `json:"status" example:"success"`
		Data     []Technician `json:"data"`
		Metadata Metadata     `json:"metadata"`
	}
	ResponseUpdateTechnicianById struct {
		Message string `json:"message" example:"technician updated"`
		Status  string `json:"status" example:"success"`
	}
	ResponsePostTechnician struct {
		Message string `json:"message" example:"technician created"`
		Status  string `json:"status" example:"success"`
	}
	RequestPostUpdateTechnicianById struct {
		Username        string  `json:"username"`
		Password        string  `json:"password"`
		Email           string  `json:"email"`
		TelegramUserId  *string `json:"telegram_user_id"`
		TelegramIsValid bool    `json:"telegram_is_valid"`
	}
	ResponseDeleteTechnicianById struct {
		Message string `json:"message" example:"technician deleted"`
		Status  string `json:"status" example:"success"`
	}
)

// For Generic Interface
func (a *Technician) GetID() uuid.UUID {
	return a.ID
}
func (a *Technician) GetPassword() string {
	return a.Password
}
