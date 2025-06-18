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
)

func TestAssetMaintenanceRepositoryCRUD(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetMaintenanceRepository(db)

	// Setup: Prepare Test Data
	admin := tests.CreateTestAdmin(t, db)
	technician := tests.CreateTestTechnician(t, db, admin.ID, admin.Email)
	asset := tests.CreateTestAsset(t, db, admin.ID)
	room := tests.CreateTestRoom(t, db)
	assetPlacement := tests.CreateTestAssetPlacement(t, db, admin.ID, technician.ID, asset.ID, room.ID)
	note := "Routine check"
	day := "Mon"
	start := entity.Time{Time: time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)}
	end := entity.Time{Time: time.Date(0, 1, 1, 11, 0, 0, 0, time.UTC)}

	maintenance := &entity.AssetMaintenance{
		MaintenanceDay:       day,
		MaintenanceHourStart: start,
		MaintenanceHourEnd:   end,
		MaintenanceNotes:     &note,
		AssetPlacementId:     assetPlacement.ID,
		MaintenanceBy:        technician.ID,
	}

	// Test 1: Create should succeed
	err := repo.Create(maintenance, admin.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, maintenance.ID)

	// Test 2: Find All should return the maintenance
	pagination := utils.Pagination{Page: 1, Limit: 10}
	results, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.True(t, total > 0)

	var found bool
	for _, m := range results {
		if m.ID == maintenance.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test 3: Find All Schedule should return results
	schedules, err := repo.FindAllSchedule()
	assert.NoError(t, err)
	fmt.Println(err)
	assert.NotEmpty(t, schedules)

	var scheduleFound bool
	for _, s := range schedules {
		if s.MaintenanceDay == day && s.AssetName == asset.AssetName {
			scheduleFound = true
			break
		}
	}
	assert.True(t, scheduleFound)

	// Test 4: Find By Asset Placement Id Maintenance By And Maintenance Day with overlapping time should error
	overlapStart := entity.Time{Time: time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)}
	overlapEnd := entity.Time{Time: time.Date(0, 1, 1, 12, 0, 0, 0, time.UTC)}
	foundAP, err := repo.FindByAssetPlacementIdMaintenanceByAndMaintenanceDay(assetPlacement.ID, technician.ID, day, overlapStart, overlapEnd)
	assert.Error(t, err)
	assert.Nil(t, foundAP)

	// Test 5: Find By Asset Placement Id Maintenance By Maintenance Day And Id with same time should error if id is different
	otherID := uuid.New()
	foundAP2, err := repo.FindByAssetPlacementIdMaintenanceByMaintenanceDayAndId(assetPlacement.ID, technician.ID, day, start, end, otherID)
	assert.Error(t, err)
	assert.Nil(t, foundAP2)

	// Test 6: Update By Id should change the hours
	newStart := entity.Time{Time: time.Date(0, 1, 1, 13, 0, 0, 0, time.UTC)}
	newEnd := entity.Time{Time: time.Date(0, 1, 1, 15, 0, 0, 0, time.UTC)}

	maintenance.MaintenanceHourStart = newStart
	maintenance.MaintenanceHourEnd = newEnd
	err = repo.UpdateById(maintenance, maintenance.ID)
	assert.NoError(t, err)

	var updated entity.AssetMaintenance
	err = db.First(&updated, "id = ?", maintenance.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, newStart.Time.Hour(), updated.MaintenanceHourStart.Time.Hour())

	// Test 7: Delete By Id should delete maintenance
	err = repo.DeleteById(maintenance.ID)
	assert.NoError(t, err)

	var check entity.AssetMaintenance
	res := db.Unscoped().First(&check, "id = ?", maintenance.ID)
	assert.Error(t, res.Error)
}
