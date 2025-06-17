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
	assert.Equal(t, "asset created", result["message"])

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
	assert.Nil(t, data["updated_at"])
	assert.Nil(t, data["deleted_at"])

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
	assert.Equal(t, "asset finding created", result["message"])

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

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, "", data["finding_category"])
	assert.IsType(t, "", data["finding_notes"])
	assert.IsType(t, "", data["asset_placement_id"])
	assert.IsType(t, "", data["created_at"])
	assert.IsType(t, "", data["finding_by_user"])
}

func TestAssetMaintenancePostCreateWithValidInput(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	payload := map[string]interface{}{
		"maintenance_day":        "Sun",
		"maintenance_hour_start": "13:00:00",
		"maintenance_hour_end":   "15:00:00",
		"maintenance_notes":      "test",
		"asset_placement_id":     "105508b2-e094-472c-85a0-8456a772b17a",
		"maintenance_by":         "465c0721-bd25-4808-9570-11444e7d0b29",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/maintenance"
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
	assert.Equal(t, "asset maintenance created", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	// Check Data Fields
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, payload["maintenance_day"], data["maintenance_day"])
	assert.Equal(t, payload["maintenance_hour_start"], data["maintenance_hour_start"])
	assert.Equal(t, payload["maintenance_hour_end"], data["maintenance_hour_end"])
	assert.Equal(t, payload["maintenance_notes"], data["maintenance_notes"])
	assert.Equal(t, payload["asset_placement_id"], data["asset_placement_id"])

	// Nullable / Optional Fields
	assert.NotEmpty(t, data["updated_at"])

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, "", data["maintenance_day"])
	assert.IsType(t, "", data["maintenance_hour_start"])
	assert.IsType(t, "", data["maintenance_hour_end"])
	assert.IsType(t, "", data["maintenance_notes"])
	assert.IsType(t, "", data["asset_placement_id"])
	assert.IsType(t, "", data["created_at"])
	assert.IsType(t, "", data["updated_at"])
	assert.IsType(t, "", data["created_by"])
	assert.IsType(t, "", data["maintenance_by"])
}

func TestAssetPlacementPostCreateWithValidInput(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	payload := map[string]interface{}{
		"asset_qty":   3,
		"asset_desc":  "buat kerja",
		"asset_id":    "114a7d9c-2760-4c2f-acc2-c22ebf3fd516",
		"room_id":     "59c87a7b-6666-401d-956f-cfb0d8a20c54",
		"asset_owner": "465c0721-bd25-4808-9570-11444e7d0b29",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/placement"
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
	assert.Equal(t, "asset placement created", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	// Check Data Fields
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, payload["asset_qty"], int(data["asset_qty"].(float64)))
	assert.Equal(t, payload["asset_desc"], data["asset_desc"])
	assert.Equal(t, payload["asset_id"], data["asset_id"])
	assert.Equal(t, payload["room_id"], data["room_id"])

	// Nullable / Optional Fields
	assert.NotEmpty(t, data["updated_at"])

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, float64(0), data["asset_qty"])
	assert.IsType(t, "", data["asset_desc"])
	assert.IsType(t, "", data["asset_id"])
	assert.IsType(t, "", data["room_id"])
	assert.IsType(t, "", data["created_at"])
	assert.IsType(t, "", data["updated_at"])
	assert.IsType(t, "", data["created_by"])
	assert.IsType(t, "", data["asset_owner"])
}

func TestAssetSoftDeleteWithValidId(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	id := "b1e8541d-23b6-4a83-9634-ed01c739a3d8"

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/" + id
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset deleted", result["message"])
}

func TestAssetHardDeleteWithValidId(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	id := "b1e8541d-23b6-4a83-9634-ed01c739a3d8"

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/destroy/" + id
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset permanentally deleted", result["message"])
}

func TestAssetFindingDeleteWithValidId(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	id := "f6cc7afa-46fa-430e-9568-51be8b49837d"

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/finding/" + id
	req, err := http.NewRequest("DELETE", url, nil)
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset finding deleted", result["message"])
}

func TestAssetRecoverWithValidId(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	id := "921d4b03-bf28-444d-b97f-30017cb82c93"

	// Exec
	url := "http://127.0.0.1:9000/api/v1/asset/recover/" + id
	req, err := http.NewRequest("PUT", url, nil)
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "asset recovered", result["message"])
}
