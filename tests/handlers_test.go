package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/say8hi/go-api-test/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler_E2E(t *testing.T) {
	requestBody := models.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "John Doe",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:8080/users/create", bytes.NewReader(jsonData))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseBody models.UserInDatabase
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, models.UserInDatabase{1, "testuser", "", "John Doe"}, responseBody)
}
