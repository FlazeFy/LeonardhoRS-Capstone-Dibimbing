package unit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"pelita/entity"
	"pelita/middleware"
	"pelita/utils"
)

func TestHashPassword(t *testing.T) {
	// Test Data
	user := &entity.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test123@example.com",
	}
	rawPassword := "nopass123"

	// Exec
	err := utils.HashPassword(user, rawPassword)

	// Test 1 : Not returned an error, empty value, have different char after hash, and less than 500 character
	assert.NoError(t, err, "hashing password should not return an error")
	assert.NotEmpty(t, user.Password, "hashed password should not be empty")
	assert.NotEqual(t, rawPassword, user.Password, "hashed password should not equal raw password")
	assert.Less(t, len(user.Password), 500, "hashed password should be less than 500 characters")

	// Test 2 : Hashed password should match original
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rawPassword))
	assert.NoError(t, err, "hashed password should match original password")
}

func TestGenerateAndValidateToken(t *testing.T) {
	// Test Data
	userID := uuid.New()
	role := "guest"

	// Exec
	token, err := utils.GenerateToken(userID, role)

	// Test 1 : Not returned an error and not empty value
	assert.NoError(t, err, "token generation should not return error")
	assert.NotEmpty(t, token, "token should not be empty")

	// Test 2: Token should be valid
	parsedUserID, err := middleware.ValidateToken(token)
	assert.NoError(t, err, "token should be valid")

	// Test 3: Parsed User ID from token should be same with raw
	assert.Equal(t, userID, parsedUserID, "parsed user ID should match original user ID")
}
