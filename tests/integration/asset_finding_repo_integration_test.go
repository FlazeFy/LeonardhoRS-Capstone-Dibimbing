package repository_test

import (
	"pelita/entity"
	"pelita/repository"
	"pelita/tests"
	"pelita/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetFindingRepositoryCreateAndFind(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetFindingRepository(db)

	// Setup dependencies
	admin := tests.CreateTestAdmin(t, db)
	user := tests.CreateTestUser(t, db)
	technician := tests.CreateTestTechnician(t, db, admin.ID, "admin@gmail.com")
	asset := tests.CreateTestAsset(t, db, admin.ID)
	room := tests.CreateTestRoom(t, db)
	placement := tests.CreateTestAssetPlacement(t, db, admin.ID, technician.ID, asset.ID, room.ID)

	// Test 1: Should Create Asset Finding
	finding := &entity.AssetFinding{
		FindingCategory:  "Damaged",
		FindingNotes:     "finding notes",
		FindingImage:     nil,
		AssetPlacementId: placement.ID,
	}

	err := repo.Create(finding, technician.ID, user.ID)
	assert.NoError(t, err)

	// Test 2: Should Find All Asset Finding
	pagination := utils.Pagination{Page: 1, Limit: 4}
	found, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.True(t, total > 0)
	var exists bool
	for _, f := range found {
		if f.ID == finding.ID {
			exists = true
			break
		}
	}
	assert.True(t, exists)

	// Test 3: Should Find All Report
	reports, err := repo.FindAllReport()
	assert.NoError(t, err)
	assert.NotNil(t, reports)

	// Test 4: Should Find All Finding Hour Total
	stats, err := repo.FindAllFindingHourTotal()
	assert.NoError(t, err)
	assert.NotNil(t, stats)

	// Test 5: Should Delete Asset Finding By Id
	err = repo.DeleteById(finding.ID)
	assert.NoError(t, err)

	var deleted entity.AssetFinding
	result := db.Unscoped().First(&deleted, "id = ?", finding.ID)
	assert.Error(t, result.Error)
}
