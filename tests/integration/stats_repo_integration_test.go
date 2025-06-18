package repository_test

import (
	"pelita/repository"
	"pelita/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsRepositoryFindMostUsedContext(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewStatsRepository(db)

	// Setup: Prepare Test Data
	admin := tests.CreateTestAdmin(t, db)
	technician := tests.CreateTestTechnician(t, db, admin.ID, "tech@example.com")
	room := tests.CreateTestRoom(t, db)
	asset := tests.CreateTestAsset(t, db, admin.ID)
	placement := tests.CreateTestAssetPlacement(t, db, admin.ID, technician.ID, asset.ID, room.ID)

	// Setup: Create Asset Maintenance with different maintenance_day values
	tests.CreateTestAssetMaintenanceWithDay(t, db, placement.ID, admin.ID, technician.ID, "Mon")
	tests.CreateTestAssetMaintenanceWithDay(t, db, placement.ID, admin.ID, technician.ID, "Mon")
	tests.CreateTestAssetMaintenanceWithDay(t, db, placement.ID, admin.ID, technician.ID, "Tue")
	tests.CreateTestAssetMaintenanceWithDay(t, db, placement.ID, admin.ID, technician.ID, "Wed")

	// Test 1: Should Find Most Used Context by Maintenance Day
	result, err := repo.FindMostUsedContext("asset_maintenances", "maintenance_day")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, len(result) > 0)

	// Test 2: Should Return 'Mon' As Most Frequent Maintenance Day
	assert.Equal(t, "Mon", result[0].Context)
	assert.True(t, result[0].Total >= 2)
}
