package factory

import (
	"pelita/entity"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func GenerateAdmin() entity.Admin {
	password := "nopass123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return entity.Admin{
		Username:        gofakeit.Username(),
		Password:        string(hashedPass),
		TelegramUserId:  nil,
		TelegramIsValid: false,
		Email:           gofakeit.Email(),
	}
}
