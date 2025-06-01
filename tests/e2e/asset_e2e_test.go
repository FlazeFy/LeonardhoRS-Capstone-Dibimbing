package e2e

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pelita/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Positive - Test Case
func TestAssetPostCreateWithValidInput(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	payload := map[string]string{
		"asset_name":     "Laptop Dell XPS 152",
		"asset_desc":     "High-end developer laptop with 32GB RAM and 1TB SSD.",
		"asset_merk":     "Dell",
		"asset_category": "electronics",
		"asset_price":    "250000",
		"asset_status":   "in_use",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)

	// Set Authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Template Response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset created successfully", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	// Check Data Fields
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, payload["asset_name"], data["asset_name"])
	assert.Equal(t, payload["asset_desc"], data["asset_desc"])
	assert.Equal(t, payload["asset_merk"], data["asset_merk"])
	assert.Equal(t, payload["asset_category"], data["asset_category"])
	assert.Equal(t, payload["asset_price"], data["asset_price"])
	assert.Equal(t, payload["asset_status"], data["asset_status"])

	// Nullable / Optional Fields
	assert.Nil(t, data["asset_image_url"])
	assert.NotEmpty(t, data["created_at"])
	assert.Nil(t, data["updated_at"])
	assert.Nil(t, data["deleted_at"])
	assert.NotEmpty(t, data["created_by"])

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, "", data["asset_name"])
	assert.IsType(t, "", data["asset_desc"])
	assert.IsType(t, "", data["asset_merk"])
	assert.IsType(t, "", data["asset_category"])
	assert.IsType(t, "", data["asset_price"])
	assert.IsType(t, "", data["asset_status"])
	assert.IsType(t, "", data["created_at"])
	assert.IsType(t, "", data["created_by"])
}

func TestAssetFindingPostCreateWithValidInput(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "tester.123@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	payload := map[string]interface{}{
		"finding_category":   "Broken",
		"finding_notes":      "Jatuh",
		"finding_image":      nil,
		"asset_placement_id": "105508b2-e094-472c-85a0-8456a772b17a",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/finding"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)

	// Set Authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Template Response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset finding created successfully", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	// Check Data Fields
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, payload["finding_category"], data["finding_category"])
	assert.Equal(t, payload["finding_notes"], data["finding_notes"])
	assert.Equal(t, payload["asset_placement_id"], data["asset_placement_id"])

	// Nullable / Optional Fields
	assert.Nil(t, data["finding_image"])
	assert.Nil(t, data["finding_by_technician"])
	assert.NotEmpty(t, data["finding_by_user"])
	assert.NotEmpty(t, data["created_at"])

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, "", data["finding_category"])
	assert.IsType(t, "", data["finding_notes"])
	assert.IsType(t, "", data["asset_placement_id"])
	assert.IsType(t, "", data["created_at"])
	assert.IsType(t, "", data["finding_by_user"])
}
