package factory

import (
	"pelita/entity"

	"github.com/brianvoe/gofakeit/v6"
	"golang.org/x/crypto/bcrypt"
)

func GenerateTechnician() entity.Technician {
	password := "nopass123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return entity.Technician{
		Username:        gofakeit.Username(),
		Password:        string(hashedPass),
		TelegramUserId:  nil,
		TelegramIsValid: false,
		Email:           gofakeit.Email(),
	}
}
