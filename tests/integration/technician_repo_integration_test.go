package repository_test

import (
	"pelita/entity"
	"pelita/repository"
	"pelita/tests"
	"pelita/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTechnicianRepositoryCreateAndFind(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewTechnicianRepository(db)

	// Setup: Prepare Test Data
	admin := tests.CreateTestAdmin(t, db)

	// Test 1: Create Technician
	tech := entity.Technician{
		Username:        "tech1",
		Password:        "pass1",
		Email:           "tech1@example.com",
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(&tech, admin.ID)
	assert.NoError(t, err)

	// Test 2: Should find by email
	foundByEmail, err := repo.FindByEmail("tech1@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, foundByEmail)
	assert.Equal(t, tech.Email, foundByEmail.Email)

	// Test 3: Should find by id
	foundById, err := repo.FindById(tech.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundById)
	assert.Equal(t, tech.ID, foundById.ID)

	// Test 4: Should find by email and unique id
	result, err := repo.FindByEmailAndId("tech1@example.com", uuid.New())
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Test 5: Should find by email and id (same ID should return nil)
	result, err = repo.FindByEmailAndId("tech1@example.com", tech.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestTechnicianRepositoryFindAll(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewTechnicianRepository(db)
	admin := tests.CreateTestAdmin(t, db)

	// Create multiple technicians
	for i := 1; i <= 3; i++ {
		tech := entity.Technician{
			ID:              uuid.New(),
			Username:        "tech" + uuid.NewString()[0:5],
			Password:        "pass",
			Email:           uuid.NewString()[0:8] + "@example.com",
			TelegramIsValid: true,
			CreatedAt:       time.Now(),
			CreatedBy:       admin.ID,
		}
		db.Create(&tech)
	}

	// Test 1: Should Get all technician with limit and offset
	pagination := utils.Pagination{Page: 1, Limit: 2}
	result, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(total), 3)
	assert.Len(t, result, 2)
}

func TestTechnicianRepositoryUpdateById(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewTechnicianRepository(db)
	admin := tests.CreateTestAdmin(t, db)
	original := tests.CreateTestTechnician(t, db, admin.ID, "update_test@example.com")

	// Update fields
	original.Username = "updated_username"
	original.TelegramIsValid = false

	// Test 1: Should update by id
	err := repo.UpdateById(&original, original.ID)
	assert.NoError(t, err)
	updated, err := repo.FindById(original.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updated_username", updated.Username)
	assert.False(t, updated.TelegramIsValid)
}

func TestTechnicianRepositoryDeleteById(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewTechnicianRepository(db)
	admin := tests.CreateTestAdmin(t, db)
	tech := tests.CreateTestTechnician(t, db, admin.ID, "delete_test@example.com")

	// Test 1: Should Delete technician by id
	err := repo.DeleteById(tech.ID)
	assert.NoError(t, err)

	// Test 2: Should not find deleted technician
	result, err := repo.FindById(tech.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}
