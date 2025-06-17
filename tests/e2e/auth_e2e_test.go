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
func TestAuthPostRegisterWithValidInput(t *testing.T) {
	// Test Data
	payload := map[string]string{
		"username": "tester123",
		"password": "nopass123",
		"email":    "tester.123@gmail.com",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/auth/register"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
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
	assert.Equal(t, "user registered", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	accessToken, ok := data["access_token"].(string)
	assert.True(t, ok, "access_token should be a string")
	assert.NotEmpty(t, accessToken)
}

func TestAuthPostLoginWithValidInput(t *testing.T) {
	// Test Data
	payload := map[string]string{
		"password": "nopass123",
		"email":    "tester.123@gmail.com",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/auth/login"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Validate Template Response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "user login", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	accessToken, ok := data["access_token"].(string)
	assert.True(t, ok, "access_token should be a string")
	assert.NotEmpty(t, accessToken)

	role, ok := data["role"].(string)
	assert.True(t, ok, "role should be a string")
	assert.NotEmpty(t, role)
}

func TestAuthPostSignOutWithValidInput(t *testing.T) {
	// Pre - Condition : Need To Login First
	username := "tester.123@gmail.com"
	password := "nopass123"
	token, _ := tests.GetAuthTokenAndRole(t, username, password)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/auth/signout"
	req, err := http.NewRequest("POST", url, nil)
	assert.NoError(t, err)

	// Set Authorization
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Exec
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Validate Template Response
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "user signout", result["message"])
}
