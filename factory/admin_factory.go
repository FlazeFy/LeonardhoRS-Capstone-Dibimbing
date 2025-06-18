package factory

import (
	"pelita/entity"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAdmin() entity.Admin {
	password := "nopass123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return entity.Admin{
		ID:              uuid.New(),
		Username:        gofakeit.Username(),
		Password:        string(hashedPass),
		TelegramUserId:  nil,
		TelegramIsValid: false,
		Email:           gofakeit.Email(),
		CreatedAt:       time.Now(),
	}
}
