package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
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

func TestCategoryFlow_E2E(t *testing.T) {
	client := &http.Client{}

	sendRequest := func(method, url string, body []byte) (*http.Response, error) {
		req, _ := http.NewRequest(method, url, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+authToken)
		return client.Do(req)
	}

	createCategory := func(name string, description string) models.Category {
		requestBody := models.CreateCategoryRequest{
			Name:        name,
			Description: description,
		}

		jsonData, _ := json.Marshal(requestBody)
		resp, _ := sendRequest(http.MethodPost, "http://0.0.0.0:8080/category/create", jsonData)
		defer resp.Body.Close()

		var responseBody models.Category
		json.NewDecoder(resp.Body).Decode(&responseBody)
		return responseBody
	}

	createdCategory := createCategory("testcategory", "desc")
	t.Run("Create category", func(t *testing.T) {
		assert.Equal(t, models.Category{ID: 1, Name: "testcategory", Description: "desc"}, createdCategory)
	})

	t.Run("Get category", func(t *testing.T) {
		resp, _ := http.Get("http://0.0.0.0:8080/category/1")
		defer resp.Body.Close()

		var responseBody models.Category
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, models.Category{ID: 1, Name: "testcategory", Description: "desc"}, responseBody)
	})

	t.Run("Get all categories", func(t *testing.T) {
		_ = createCategory("testcategory", "desc")
		resp, _ := http.Get("http://0.0.0.0:8080/category/")
		defer resp.Body.Close()

		var responseBody []models.Category
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, []models.Category{{ID: 1, Name: "testcategory", Description: "desc"},
			{ID: 2, Name: "testcategory", Description: "desc"}}, responseBody)
	})

	t.Run("Update category", func(t *testing.T) {
		updatedRequestBody := models.CreateCategoryRequest{Name: "new_test_name", Description: "new_test_desc"}
		jsonData, _ := json.Marshal(updatedRequestBody)

		_, _ = sendRequest(http.MethodPatch, "http://0.0.0.0:8080/category/"+strconv.Itoa(createdCategory.ID), jsonData)

		resp, _ := http.Get("http://0.0.0.0:8080/category/" + strconv.Itoa(createdCategory.ID))
		defer resp.Body.Close()

		var updatedCategory models.Category
		err := json.NewDecoder(resp.Body).Decode(&updatedCategory)
		assert.NoError(t, err)

		assert.Equal(t, "new_test_name", updatedCategory.Name)
		assert.Equal(t, "new_test_desc", updatedCategory.Description)
	})

	t.Run("Delete category", func(t *testing.T) {
		resp, _ := sendRequest(http.MethodDelete, "http://0.0.0.0:8080/category/2", nil)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var deleteResponse map[string]string
		err := json.NewDecoder(resp.Body).Decode(&deleteResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Category deleted successfully", deleteResponse["message"])
	})

	t.Run("Get deleted category", func(t *testing.T) {
		resp, _ := http.Get("http://0.0.0.0:8080/category/2")
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
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
		Categories:  []models.Category{{ID: 1, Name: "new_test_name", Description: "new_test_desc"}}},
		responseBody)
}
