package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pelita/config"
	"pelita/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// For E2E Test
func GetAuthTokenAndRole(t *testing.T, email, password string) (string, string) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonPayload, _ := json.Marshal(payload)

	url := "http://127.0.0.1:9000/api/v1/auth/login"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "user login", result["message"])

	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	accessToken, ok := data["access_token"].(string)
	assert.True(t, ok, "access_token should be a string")
	assert.NotEmpty(t, accessToken)

	role, ok := data["role"].(string)
	assert.True(t, ok, "role should be a string")
	assert.NotEmpty(t, role)

	return accessToken, role
}

// For Integration Test
func SetupTestDB(t *testing.T) *gorm.DB {
	// Load Env
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("error loading ENV")
	}

	db := config.ConnectTestDatabase(t)

	err = db.Migrator().DropTable(
		&entity.Admin{},
		&entity.User{},
		&entity.Technician{},
		&entity.History{},
		&entity.Asset{},
		&entity.AssetPlacement{},
		&entity.AssetMaintenance{},
		&entity.AssetFinding{},
	)
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&entity.Admin{},
		&entity.User{},
		&entity.Technician{},
		&entity.History{},
		&entity.Asset{},
		&entity.AssetPlacement{},
		&entity.AssetMaintenance{},
		&entity.AssetFinding{},
	)
	assert.NoError(t, err)

	return db
}

func CreateTestAdmin(t *testing.T, db *gorm.DB) entity.Admin {
	admin := entity.Admin{
		ID:              uuid.New(),
		Username:        "admin_test",
		Password:        "hashed_password",
		Email:           "admin@test.com",
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := db.Create(&admin).Error
	assert.NoError(t, err)
	return admin
}

func CreateTestTechnician(t *testing.T, db *gorm.DB, createdBy uuid.UUID, email string) entity.Technician {
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

func CreateTestUser(t *testing.T, db *gorm.DB) entity.User {
	user := entity.User{
		ID:              uuid.New(),
		Username:        "user_test",
		Password:        "hashed_password",
		Email:           "user@test.com",
		TelegramIsValid: true,
		CreatedAt:       time.Now(),
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)
	return user
}

func CreateTestAsset(t *testing.T, db *gorm.DB, adminID uuid.UUID) *entity.Asset {
	assetName := "Test Asset"
	assetCategory := "Test Category"
	assetStatus := "new"
	assetMerk := "Test Merk"
	assetPrice := "12345"
	assetImageUrl := "http://example.com/test.jpg"
	assetDesc := "Test asset description"

	asset := &entity.Asset{
		ID:            uuid.New(),
		AssetName:     assetName,
		AssetCategory: assetCategory,
		AssetStatus:   assetStatus,
		AssetDesc:     &assetDesc,
		AssetMerk:     &assetMerk,
		AssetPrice:    &assetPrice,
		AssetImageURL: &assetImageUrl,
		CreatedAt:     time.Now(),
		CreatedBy:     adminID,
	}

	err := db.Create(asset).Error
	assert.NoError(t, err)

	return asset
}

func CreateTestRoom(t *testing.T, db *gorm.DB) *entity.Room {
	room := &entity.Room{
		ID:        uuid.New(),
		Floor:     "2",
		RoomName:  "Test Room A",
		RoomDept:  "Engineering",
		CreatedAt: time.Now(),
	}

	err := db.Create(room).Error
	assert.NoError(t, err)

	return room
}

func CreateTestAssetPlacement(t *testing.T, db *gorm.DB, adminId, technicianId, assetId, roomId uuid.UUID) *entity.AssetPlacement {
	placement := &entity.AssetPlacement{
		ID:         uuid.New(),
		AssetQty:   3,
		AssetDesc:  nil,
		AssetId:    assetId,
		RoomId:     roomId,
		CreatedBy:  adminId,
		AssetOwner: technicianId,
		CreatedAt:  time.Now(),
	}

	err := db.Create(placement).Error
	assert.NoError(t, err)
	return placement
}
