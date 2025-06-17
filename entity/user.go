package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID              uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
		Username        string    `json:"username" gorm:"type:varchar(36);not null"`
		Password        string    `json:"password" gorm:"type:varchar(500);not null"`
		Email           string    `json:"email" gorm:"type:varchar(500);not null"`
		TelegramUserId  *string   `json:"telegram_user_id" gorm:"type:varchar(36);null"`
		TelegramIsValid bool      `json:"telegram_is_valid"`
		CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp;not null"`
	}
	UserAuth struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	// For Response Only
	ResponseGetAllUser struct {
		Message  string   `json:"message" example:"user fetched"`
		Status   string   `json:"status" example:"success"`
		Data     []User   `json:"data"`
		Metadata Metadata `json:"metadata"`
	}
	ResponseGetMyProfile struct {
		Message string `json:"message" example:"user fetched"`
		Status  string `json:"status" example:"success"`
		Data    []User `json:"data"`
	}
	ResponsePostLogin struct {
		Message string    `json:"message" example:"user login successfully"`
		Status  string    `json:"status" example:"success"`
		Data    LoginData `json:"data"`
	}
	LoginData struct {
		AccessToken string `json:"access_token" example:"<your_access_token>"`
		Role        string `json:"role" example:"admin"`
	}
	ResponsePostSignOut struct {
		Message string `json:"message" example:"User signed out"`
		Status  string `json:"status" example:"success"`
	}
	RequestPostRegister struct {
		Username        string  `json:"username"`
		Password        string  `json:"password"`
		Email           string  `json:"email"`
		TelegramUserId  *string `json:"telegram_user_id"`
		TelegramIsValid bool    `json:"telegram_is_valid"`
	}
)

// For Generic Interface
func (a *User) GetID() uuid.UUID {
	return a.ID
}
func (a *User) GetPassword() string {
	return a.Password
}
