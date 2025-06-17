package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	Metadata struct {
		Limit      int `json:"limit"`
		Page       int `json:"page"`
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
	}
	// All Role
	Account interface {
		GetID() uuid.UUID
		GetPassword() string
	}
	MyProfile struct {
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	// For Response
	ResponseBadRequest struct {
		Message string `json:"message" example:"targetCol is not valid"`
		Status  string `json:"status" example:"failed"`
	}
	ResponseNotFound struct {
		Message string `json:"message" example:"asset not found"`
		Status  string `json:"status" example:"failed"`
	}
)
