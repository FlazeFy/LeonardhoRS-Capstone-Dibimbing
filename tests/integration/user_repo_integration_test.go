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

func TestUserRepositoryCreateAndFind(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewUserRepository(db)

	// Setup: Prepare Test Data
	telegramId := "12341"
	user := &entity.User{
		ID:              uuid.New(),
		Username:        "testuser",
		Password:        "hashed_pass",
		Email:           "testuser@example.com",
		TelegramUserId:  &telegramId,
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}

	// Test 1: Create user should not return error
	err := repo.Create(user)
	assert.NoError(t, err, "creating user should not error")

	// Test 2: Should find by exact email
	foundByEmail, err := repo.FindByEmail("testuser@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundByEmail)
	assert.Equal(t, user.ID, foundByEmail.ID)

	// Test 3: Should find by username or email
	foundByEither, err := repo.FindByUsernameOrEmail("testuser", "notexist@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundByEither)
	assert.Equal(t, user.ID, foundByEither.ID)

	// Test 4: Should return nil for non-existing email
	notFound, err := repo.FindByEmail("unknown@example.com")
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestUserRepositoryFindById(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewUserRepository(db)
	telegramId := "12341"
	user := &entity.User{
		ID:              uuid.New(),
		Username:        "testuser",
		Password:        "pass",
		Email:           "testuser@example.com",
		TelegramUserId:  &telegramId,
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}

	err := db.Create(user).Error
	assert.NoError(t, err, "failed to insert user")

	// Test 1: Should find user profile by ID with guest role
	result, err := repo.FindById(user.ID.String(), "guest")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "testuser@example.com", result.Email)

	// Test 2: Should return nil for unknown ID
	result, err = repo.FindById(uuid.New().String(), "guest")
	assert.NoError(t, err)
	assert.Nil(t, result)
}
