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

func TestAssetRepositoryCreateAndFind(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetRepository(db)
	admin := tests.CreateTestAdmin(t, db)
	assetDesc := "test desc"
	assetName := "test asset"
	assetCategory := "test category"
	assetStatus := "new"
	assetMerk := "merk"
	assetPrice := "10000"
	assetImageUrl := "http://example.com/image.jpg"

	asset := &entity.Asset{
		AssetName:     assetName,
		AssetDesc:     &assetDesc,
		AssetMerk:     &assetMerk,
		AssetCategory: assetCategory,
		AssetPrice:    &assetPrice,
		AssetStatus:   assetStatus,
		AssetImageURL: &assetImageUrl,
	}

	// Test 1: Create should succeed
	err := repo.Create(asset, admin.ID)
	assert.NoError(t, err)

	// Test 2: Find All should return the created asset
	pagination := utils.Pagination{Page: 1, Limit: 4}
	result, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	assert.True(t, total > 0)
	assert.NotEmpty(t, result)

	// Test 3: Find By Asset Name Category And Merk should return asset
	found, err := repo.FindByAssetNameCategoryAndMerk(assetName, assetCategory, &assetMerk)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, asset.ID, found.ID)

	// Test 4: Find By Asset Name Category Merk And Id (wrong ID) should return the asset
	differentID := uuid.New()
	found2, err := repo.FindByAssetNameCategoryMerkAndId(assetName, assetCategory, &assetMerk, differentID)
	assert.NoError(t, err)
	assert.NotNil(t, found2)

	// Test 5: Find Deleted should return nothing
	deleted, err := repo.FindDeleted()
	assert.Error(t, err)
	assert.Nil(t, deleted)
}

func TestAssetRepositorySoftAndRecoverDelete(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetRepository(db)
	adminId := uuid.New()
	asset := &entity.Asset{
		ID:            uuid.New(),
		AssetName:     "test asset",
		AssetCategory: "test category",
		AssetStatus:   "new",
		CreatedAt:     time.Now(),
		CreatedBy:     adminId,
	}

	err := db.Create(asset).Error
	assert.NoError(t, err)

	// Test 1: Asset should be deleted by id
	err = repo.SoftDeleteById(asset.ID)
	assert.NoError(t, err)

	// Test 2: Asset should not be returned in Find All
	pagination := utils.Pagination{Page: 1, Limit: 4}
	assets, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	for _, a := range assets {
		assert.NotEqual(t, asset.ID, a.ID)
	}
	assert.True(t, total >= 0)

	// Test 3: Deleted asset should appear in FindDeleted
	deleted, err := repo.FindDeleted()
	assert.NoError(t, err)
	var found bool
	for _, a := range deleted {
		if a.ID == asset.ID {
			found = true
			break
		}
	}
	assert.True(t, found)

	// Test 4: RecoverDeletedById should restore the asset
	err = repo.RecoverDeletedById(asset.ID)
	assert.NoError(t, err)

	// Test 5: Asset should now appear in Find All again
	all, total, err := repo.FindAll(pagination)
	assert.NoError(t, err)
	var exists bool
	for _, a := range all {
		if a.ID == asset.ID {
			exists = true
			break
		}
	}
	assert.True(t, exists)
}

func TestAssetRepositoryUpdateAndHardDelete(t *testing.T) {
	db := tests.SetupTestDB(t)
	repo := repository.NewAssetRepository(db)
	adminId := uuid.New()
	asset := &entity.Asset{
		ID:            uuid.New(),
		AssetName:     "test asset",
		AssetCategory: "test category",
		AssetStatus:   "new",
		CreatedAt:     time.Now(),
		CreatedBy:     adminId,
	}

	err := db.Create(asset).Error
	assert.NoError(t, err)

	// Test 1: Asset should be Update By Id
	asset.AssetName = "UpdatedName"
	err = repo.UpdateById(asset, asset.ID)
	assert.NoError(t, err)

	var updated entity.Asset
	_ = db.First(&updated, "id = ?", asset.ID).Error
	assert.Equal(t, "UpdatedName", updated.AssetName)

	// Test 2: Asset should be deleted by id
	err = repo.SoftDeleteById(asset.ID)
	assert.NoError(t, err)

	// Test 3:  Asset should be permanentally deleted by id
	err = repo.HardDeleteById(asset.ID)
	assert.NoError(t, err)

	var check entity.Asset
	result := db.Unscoped().First(&check, "id = ?", asset.ID)
	assert.Error(t, result.Error)
}
