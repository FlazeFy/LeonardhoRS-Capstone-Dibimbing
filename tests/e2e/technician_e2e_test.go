package e2e

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"pelita/tests"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Positive - Test Case
func TestTechnicianDeleteWithValidId(t *testing.T) {
	// Pre - Condition : Need To Login First as Admin
	username := "flazen.edu@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Test Data
	id := "125c0721-bd25-4808-9570-11444e7d0b29"

	// Exec
	url := "http://127.0.0.1:9000/api/v1/technician/" + id
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
	assert.Equal(t, "technician deleted", result["message"])
}
