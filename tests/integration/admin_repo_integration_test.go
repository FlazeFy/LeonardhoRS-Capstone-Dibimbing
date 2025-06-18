package repository_test

import (
	"pelita/entity"
	"pelita/repository"
	"pelita/tests"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminRepositoryFindByEmail(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAdminRepository(db)

	// Setup: Prepare Test Data
	admin := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin_user",
		Password:        "hashed_password",
		Email:           "admin@gmail.com",
		TelegramUserId:  nil,
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := db.Create(&admin).Error
	assert.NoError(t, err, "failed to insert admin")

	// Test 1: Should find the admin by valid email (existing)
	result, err := repo.FindByEmail("admin@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, admin.ID, result.ID)

	// Test 2: Should not find admin by invalid email (non-existing)
	result, err = repo.FindByEmail("missing@gmail.com")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestAdminRepositoryFindAllContact(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAdminRepository(db)
	telegramId := "12341"

	admin1 := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin1",
		Password:        "pass1",
		Email:           "admin1@gmail.com",
		TelegramUserId:  &telegramId,
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	admin2 := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin2",
		Password:        "pass2",
		Email:           "admin2@gmail.com",
		TelegramUserId:  &telegramId,
		TelegramIsValid: false,
		CreatedAt:       time.Now(),
	}
	admin3 := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin3",
		Password:        "pass3",
		Email:           "admin3@gmail.com",
		TelegramUserId:  nil,
		TelegramIsValid: false,
		CreatedAt:       time.Now(),
	}

	db.Create(&admin1)
	db.Create(&admin2)
	db.Create(&admin3)

	// Test 1: Only admin1 have valid telegram user id and telegram is valid
	contacts, err := repo.FindAllContact()
	assert.NoError(t, err)
	assert.Len(t, contacts, 1)

	// Test 2: Ensure order by admin username in ascending
	assert.Equal(t, "admin1", contacts[0].Username)
}
