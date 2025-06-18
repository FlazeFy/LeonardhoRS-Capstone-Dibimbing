package repository_test

import (
	"pelita/entity"
	"pelita/repository"
	"pelita/tests"
	"pelita/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAssetPlacementRepositoryCreateFindUpdateDelete(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetPlacementRepository(db)

	// Setup: Prepare Test Data
	admin := tests.CreateTestAdmin(t, db)
	technician := tests.CreateTestTechnician(t, db, admin.ID, admin.Email)
	asset := tests.CreateTestAsset(t, db, admin.ID)
	room := tests.CreateTestRoom(t, db)
	assetDesc := "Placed in Room A"
	assetQty := 2

	assetPlacement := &entity.AssetPlacement{
		AssetQty:   assetQty,
		AssetDesc:  &assetDesc,
		AssetId:    asset.ID,
		RoomId:     room.ID,
		CreatedBy:  admin.ID,
		AssetOwner: technician.ID,
	}

	// Test 1: Create should succeed
	err := repo.Create(assetPlacement, admin.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, assetPlacement.ID)

	// Test 2: Find All should return the placement
	pagination := utils.Pagination{Page: 1, Limit: 10}
	results, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.True(t, total > 0)
	var found bool
	for _, ap := range results {
		if ap.ID == assetPlacement.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test 3: Find By Asset Id and RoomId should return valid asset
	foundAP, err := repo.FindByAssetIdAndRoomId(asset.ID, room.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundAP)
	assert.Equal(t, assetPlacement.ID, foundAP.ID)

	// Test 4: Find By Asset Id, RoomId and different ID should return valid asset
	differentID := uuid.New()
	foundAP2, err := repo.FindByAssetIdRoomIdAndId(asset.ID, room.ID, differentID)
	assert.NoError(t, err)
	assert.NotNil(t, foundAP2)
	assert.Equal(t, assetPlacement.ID, foundAP2.ID)

	// Test 5: Update By Id should update Asset Qty
	assetPlacement.AssetQty = 5
	err = repo.UpdateById(assetPlacement, assetPlacement.ID)
	assert.NoError(t, err)

	var updated entity.AssetPlacement
	_ = db.First(&updated, "id = ?", assetPlacement.ID).Error
	assert.Equal(t, 5, updated.AssetQty)

	// Test 6: Delete By Id should remove the asset placement
	err = repo.DeleteById(assetPlacement.ID)
	assert.NoError(t, err)

	var check entity.AssetPlacement
	result := db.Unscoped().First(&check, "id = ?", assetPlacement.ID)
	assert.Error(t, result.Error)
}
