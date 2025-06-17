package repository_test

import (
	"fmt"
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

func createTestHistory(t *testing.T, db *gorm.DB, h entity.History) entity.History {
	h.ID = uuid.New()
	h.CreatedAt = time.Now()
	err := db.Create(&h).Error
	assert.NoError(t, err)
	return h
}

func TestHistoryRepositoryFindAll(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewHistoryRepository(db)

	// Models
	admin := tests.CreateTestAdmin(t, db)
	tech := tests.CreateTestTechnician(t, db, admin.ID, "create_test@example.com")
	user := tests.CreateTestUser(t, db)

	// Create multiple histories
	for i := 0; i < 5; i++ {
		createTestHistory(t, db, entity.History{
			AdminID:     &admin.ID,
			TypeUser:    "admin",
			TypeHistory: fmt.Sprintf("Admin Log %d", i),
		})
		createTestHistory(t, db, entity.History{
			TechnicianID: &tech.ID,
			TypeUser:     "technician",
			TypeHistory:  fmt.Sprintf("Technician Log %d", i),
		})
		createTestHistory(t, db, entity.History{
			UserID:      &user.ID,
			TypeUser:    "guest",
			TypeHistory: fmt.Sprintf("Guest Log %d", i),
		})
	}

	// Pagination
	pagination := utils.Pagination{Page: 1, Limit: 4}

	// Query
	histories, total, err := repo.FindAll(pagination)

	// Assert
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, int(total), 10)
	assert.Len(t, histories, 4)
}

func TestHistoryRepositoryFindMyAdmin(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewHistoryRepository(db)

	// Models
	admin := tests.CreateTestAdmin(t, db)

	// Create history for admin
	for i := 0; i < 3; i++ {
		createTestHistory(t, db, entity.History{
			AdminID:     &admin.ID,
			TypeUser:    "admin",
			TypeHistory: fmt.Sprintf("AdminLog%d", i),
		})
	}

	// Pagination
	pagination := utils.Pagination{Page: 1, Limit: 2}

	// Query
	history, total, err := repo.FindMy(pagination, admin.ID, "admin")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, history, 2)
	for _, h := range history {
		assert.Equal(t, "admin", h.TypeUser)
		assert.Equal(t, admin.ID, *h.AdminID)
	}
}

func TestHistoryRepositoryFindMyTechnician(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewHistoryRepository(db)

	// Models
	admin := tests.CreateTestAdmin(t, db)
	tech := tests.CreateTestTechnician(t, db, admin.ID, "create_test@example.com")

	// Create history for technician
	createTestHistory(t, db, entity.History{
		TechnicianID: &tech.ID,
		TypeUser:     "technician",
		TypeHistory:  "TechLogged",
	})

	// Pagination
	pagination := utils.Pagination{Page: 1, Limit: 5}

	// Query
	history, total, err := repo.FindMy(pagination, tech.ID, "technician")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, history, 1)
	assert.Equal(t, tech.ID, *history[0].TechnicianID)
	assert.Equal(t, "technician", history[0].TypeUser)
}

func TestHistoryRepository_FindMy_Guest(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewHistoryRepository(db)

	// Models
	user := tests.CreateTestUser(t, db)

	// Create history for guest
	createTestHistory(t, db, entity.History{
		UserID:      &user.ID,
		TypeUser:    "guest",
		TypeHistory: "UserLogin",
	})

	// Pagination
	pagination := utils.Pagination{Page: 1, Limit: 5}

	// Query
	history, total, err := repo.FindMy(pagination, user.ID, "guest")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, history, 1)
	assert.Equal(t, user.ID, *history[0].UserID)
	assert.Equal(t, "guest", history[0].TypeUser)
}
