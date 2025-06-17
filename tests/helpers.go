package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pelita/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
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
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err, "failed to connect test DB")

	err = db.AutoMigrate(&entity.Admin{})
	assert.NoError(t, err, "failed to migrate admin schema")

	return db
}
