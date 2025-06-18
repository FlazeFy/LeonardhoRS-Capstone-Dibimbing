package repository

import (
	"pelita/entity"
	"pelita/repository"
	"pelita/tests"
	"pelita/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoomRepositoryCRUD(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewRoomRepository(db)

	// Setup: Create test room
	room := &entity.Room{
		Floor:    "1",
		RoomName: "Room A",
		RoomDept: "IT",
	}
	err := repo.Create(room)
	assert.NoError(t, err)

	// Setup: Create test dependencies
	admin := tests.CreateTestAdmin(t, db)
	asset := tests.CreateTestAsset(t, db, admin.ID)
	technician := tests.CreateTestTechnician(t, db, admin.ID, "tech@gmail.com")
	tests.CreateTestAssetPlacement(t, db, admin.ID, technician.ID, asset.ID, room.ID)

	// Test 1: Should Find All Rooms
	pagination := utils.Pagination{Page: 1, Limit: 10}
	rooms, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.True(t, total > 0)
	var exists bool
	for _, r := range rooms {
		if r.ID == room.ID {
			exists = true
			break
		}
	}
	assert.True(t, exists)

	// Test 2: Should Find By Room Name And Floor
	found, err := repo.FindByRoomNameAndFloor("Room A", "1")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, room.ID, found.ID)

	// Test 3: Should Return Nil When Find By Room Name Floor And Id With Same Id
	dupe, err := repo.FindByRoomNameFloorAndId("Room A", "1", room.ID)
	assert.NoError(t, err)
	assert.Nil(t, dupe)

	// Test 4: Should Update Room By Id
	room.RoomDept = "Engineering"
	err = repo.UpdateById(room, room.ID)
	assert.NoError(t, err)

	updated, err := repo.FindByRoomNameAndFloor("Room A", "1")
	assert.NoError(t, err)
	assert.Equal(t, "Engineering", updated.RoomDept)

	// Test 5: Should Find Room Asset By Floor And Room Name
	assetResults, err := repo.FindRoomAssetByFloorAndRoomName(room.Floor, room.RoomName)
	assert.NoError(t, err)
	assert.NotEmpty(t, assetResults)
	exists = false
	for _, a := range assetResults {
		if a.RoomName == "Room A" {
			exists = true
			break
		}
	}
	assert.True(t, exists)

	// Test 6: Should Find Room Asset Short By Floor And Room Name
	shortResults, err := repo.FindRoomAssetShortByFloorAndRoomName(room.Floor, room.RoomName)
	assert.NoError(t, err)
	assert.NotEmpty(t, shortResults)
	exists = false
	for _, a := range shortResults {
		if a.RoomName == "Room A" {
			exists = true
			break
		}
	}
	assert.True(t, exists)

	// Test 7: Should Return All Room Assets When Room Name is "all"
	allAssets, err := repo.FindRoomAssetByFloorAndRoomName("1", "all")
	assert.NoError(t, err)
	assert.NotEmpty(t, allAssets)

	allShortAssets, err := repo.FindRoomAssetShortByFloorAndRoomName("1", "all")
	assert.NoError(t, err)
	assert.NotEmpty(t, allShortAssets)

	// Test 8: Should Delete Room By Id
	err = repo.DeleteById(room.ID)
	assert.NoError(t, err)

	var deleted entity.Room
	result := db.Unscoped().First(&deleted, "id = ?", room.ID)
	assert.Error(t, result.Error)
}
