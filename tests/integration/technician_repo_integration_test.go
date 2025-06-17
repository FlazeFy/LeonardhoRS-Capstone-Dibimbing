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
	"gorm.io/gorm"
)

func createTestAdmin(t *testing.T, db *gorm.DB) uuid.UUID {
	admin := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin",
		Password:        "hashed_password",
		Email:           "admin@example.com",
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := db.Create(&admin).Error
	assert.NoError(t, err)
	return admin.ID
}

func createTestTechnician(t *testing.T, db *gorm.DB, createdBy uuid.UUID, email string) entity.Technician {
	technician := entity.Technician{
		ID:              uuid.New(),
		Username:        "tech_user",
		Password:        "hashed_password",
		Email:           email,
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
		CreatedBy:       createdBy,
	}
	err := db.Create(&technician).Error
	assert.NoError(t, err)
	return technician
}

func TestTechnicianRepositoryCreateAndFind(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewTechnicianRepository(db)
	adminId := createTestAdmin(t, db)

	// Test 1: Create Technician
	tech := entity.Technician{
		Username:        "tech1",
		Password:        "pass1",
		Email:           "tech1@example.com",
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := repo.Create(&tech, adminId)
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
	adminId := createTestAdmin(t, db)

	// Create multiple technicians
	for i := 1; i <= 3; i++ {
		tech := entity.Technician{
			ID:              uuid.New(),
			Username:        "tech" + uuid.NewString()[0:5],
			Password:        "pass",
			Email:           uuid.NewString()[0:8] + "@example.com",
			TelegramIsValid: true,
			CreatedAt:       time.Now(),
			CreatedBy:       adminId,
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
	adminId := createTestAdmin(t, db)
	original := createTestTechnician(t, db, adminId, "update_test@example.com")

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
	adminId := createTestAdmin(t, db)
	tech := createTestTechnician(t, db, adminId, "delete_test@example.com")

	// Test 1: Should Delete technician by id
	err := repo.DeleteById(tech.ID)
	assert.NoError(t, err)

	// Test 2: Should not find deleted technician
	result, err := repo.FindById(tech.ID)
	assert.NoError(t, err)
	assert.Nil(t, result)
}
