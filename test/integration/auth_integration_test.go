package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8099/auth" // Adjust as per your docker-compose/services
)

func TestRegisterAndLoginIntegration(t *testing.T) {
	// Start with a clean database or use test containers as needed.
	// This example assumes the auth service is running and accessible.
	//for uniqueness
	randEmail := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	randName := fmt.Sprintf("name%d", time.Now().UnixNano())
	// Step 1: Register a new user
	registerPayload := map[string]string{
		"name":     randName,
		"email":    randEmail,
		"password": "TestPassword123!",
	}
	body, _ := json.Marshal(registerPayload)

	regResp, err := http.Post(baseURL+"/register", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)

	// Read response body for debugging
	respBody, _ := io.ReadAll(regResp.Body)
	regResp.Body.Close()

	// Print response body if there's an error
	if regResp.StatusCode != http.StatusOK {
		t.Logf("Response Body: %s", string(respBody))
	}
	assert.Equal(t, http.StatusOK, regResp.StatusCode)

	// Step 2: Login as the new user
	loginPayload := map[string]string{
		"identifier": randEmail,
		"password":   "TestPassword123!",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginResp, err := http.Post(baseURL+"/login", "application/json", bytes.NewReader(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	// Step 3: Decode login response and check for token
	var loginData struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginData.AccessToken)
	assert.NotEmpty(t, loginData.RefreshToken)
}

func TestInvalidLoginIntegration(t *testing.T) {
	// Step 1: Try to login with invalid credentials
	loginPayload := map[string]string{
		"email":    "nonexistent_user@example.com",
		"password": "wrongpassword",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginResp, err := http.Post(baseURL+"/login", "application/json", bytes.NewReader(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, loginResp.StatusCode)
}
