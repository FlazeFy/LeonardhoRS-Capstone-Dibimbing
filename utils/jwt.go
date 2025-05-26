package utils

import (
	"pelita/config"
	"pelita/entity"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     time.Now().Add(config.GetJWTExpirationDuration()).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.GetJWTSecret())
}

func HashPassword(u *entity.User, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPass)

	return nil
}

func CheckPassword(u *entity.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
