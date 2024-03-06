package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/say8hi/go-api-test/internal/models"
	"github.com/stretchr/testify/assert"
)

var hash = sha256.Sum256([]byte("passwordtestuser"))
var authToken = hex.EncodeToString(hash[:])

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

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var responseBody models.UserInDatabase
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, models.UserInDatabase{ID: 1, Username: "testuser", PasswordHash: "", FullName: "John Doe"}, responseBody)
}

func TestCreateCategoryHandler_E2E(t *testing.T) {
	requestBody := models.CreateCategoryRequest{
		Name:        "testcategory",
		Description: "desc",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:8080/category/create", bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+authToken)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var responseBody models.Category
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, models.Category{ID: 1, Name: "testcategory", Description: "desc"}, responseBody)
}

func TestCreateProductHandler_E2E(t *testing.T) {
	requestBody := models.CreateProductRequest{
		Name:        "testproduct",
		Description: "desc",
		Price:       4.5,
		Categories:  []int{1},
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:8080/product/create", bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "Bearer "+authToken)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var responseBody models.Product
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	assert.NoError(t, err)

	assert.Equal(t, models.Product{
		ID:          1,
		Name:        "testproduct",
		Description: "desc",
		Price:       4.5,
		Categories:  []models.Category{{ID: 1, Name: "testcategory", Description: "desc"}}},
		responseBody)
}
