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
var serverURL = "http://0.0.0.0:8081"

func TestCreateUserHandler_E2E(t *testing.T) {
	requestBody := models.CreateUserRequest{
		Username: "testuser",
		Password: "password",
		FullName: "John Doe",
	}

	jsonData, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, serverURL+"/users/create", bytes.NewReader(jsonData))
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
		resp, _ := sendRequest(http.MethodPost, serverURL+"/category/create", jsonData)
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
		resp, _ := http.Get(serverURL + "/category/" + strconv.Itoa(createdCategory.ID))
		defer resp.Body.Close()

		var responseBody models.Category
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, models.Category{ID: 1, Name: "testcategory", Description: "desc"}, responseBody)
	})

	t.Run("Get all categories", func(t *testing.T) {
		_ = createCategory("testcategory2", "desc")
		resp, _ := http.Get(serverURL + "/category/")
		defer resp.Body.Close()

		var responseBody []models.Category
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, []models.Category{{ID: 1, Name: "testcategory", Description: "desc"},
			{ID: 2, Name: "testcategory2", Description: "desc"}}, responseBody)
	})

	t.Run("Update category", func(t *testing.T) {
		updatedRequestBody := models.CreateCategoryRequest{Name: "new_test_name", Description: "new_test_desc"}
		jsonData, _ := json.Marshal(updatedRequestBody)

		_, _ = sendRequest(http.MethodPatch, serverURL+"/category/"+strconv.Itoa(createdCategory.ID), jsonData)

		resp, _ := http.Get(serverURL + "/category/" + strconv.Itoa(createdCategory.ID))
		defer resp.Body.Close()

		var updatedCategory models.Category
		err := json.NewDecoder(resp.Body).Decode(&updatedCategory)
		assert.NoError(t, err)

		assert.Equal(t, "new_test_name", updatedCategory.Name)
		assert.Equal(t, "new_test_desc", updatedCategory.Description)
	})

	categoryToDelete := createCategory("delname", "deldesc")
	t.Run("Delete category", func(t *testing.T) {
		resp, _ := sendRequest(http.MethodDelete, serverURL+"/category/"+strconv.Itoa(categoryToDelete.ID), nil)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var deleteResponse map[string]string
		err := json.NewDecoder(resp.Body).Decode(&deleteResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Category deleted successfully", deleteResponse["message"])
	})

	t.Run("Get deleted category", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/category/" + strconv.Itoa(categoryToDelete.ID))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Get all categories2", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/category/")
		defer resp.Body.Close()

		var responseBody []models.Category
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, []models.Category{{ID: 1, Name: "new_test_name", Description: "new_test_desc"},
			{ID: 2, Name: "testcategory2", Description: "desc"}}, responseBody)
	})
}

func TestProductFlow_E2E(t *testing.T) {
	client := &http.Client{}

	sendRequest := func(method, url string, body []byte) (*http.Response, error) {
		req, _ := http.NewRequest(method, url, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+authToken)
		return client.Do(req)
	}

	createProduct := func(name, description string, price float64, categories []string) models.Product {
		requestBody := models.CreateProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
			Categories:  categories,
		}
		jsonData, _ := json.Marshal(requestBody)
		resp, _ := sendRequest(http.MethodPost, serverURL+"/product/create", jsonData)
		defer resp.Body.Close()

		var responseBody models.Product
		json.NewDecoder(resp.Body).Decode(&responseBody)
		return responseBody
	}

	createdProduct := createProduct("testproduct", "desc", 9.99, []string{"new_test_name", "testcategory2"})

	t.Run("Create product", func(t *testing.T) {
		assert.Equal(t, models.Product{
			ID:          1,
			Name:        "testproduct",
			Description: "desc",
			Price:       9.99,
			Categories: []models.Category{
				{ID: 1, Name: "new_test_name", Description: "new_test_desc"},
				{ID: 2, Name: "testcategory2", Description: "desc"},
			},
		},
			createdProduct)
	})

	t.Run("Get product", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/product/" + strconv.Itoa(createdProduct.ID))
		defer resp.Body.Close()

		var responseBody models.Product
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, models.Product{
			ID:          1,
			Name:        "testproduct",
			Description: "desc",
			Price:       9.99,
			Categories: []models.Category{
				{ID: 1, Name: "new_test_name", Description: "new_test_desc"},
				{ID: 2, Name: "testcategory2", Description: "desc"},
			},
		},
			responseBody)
	})

	t.Run("Get all products in category", func(t *testing.T) {
		_ = createProduct("second", "desc", 5.5, []string{"new_test_name"})

		resp, _ := http.Get(serverURL + "/category/1/products")
		defer resp.Body.Close()

		var responseBody []models.Product
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		assert.NoError(t, err)
		assert.Equal(t, []models.Product{
			{
				ID:          1,
				Name:        "testproduct",
				Description: "desc",
				Price:       9.99,
				Categories: []models.Category{
					{ID: 1, Name: "new_test_name", Description: "new_test_desc"},
				},
			},
			{
				ID:          2,
				Name:        "second",
				Description: "desc",
				Price:       5.5,
				Categories: []models.Category{
					{ID: 1, Name: "new_test_name", Description: "new_test_desc"},
				},
			},
		},
			responseBody)
	})

	t.Run("Update product", func(t *testing.T) {
		newName := "new_test_name"
		newDescription := "new_test_desc"
		newPrice := 10.99
		updatedRequestBody := models.CreateProductRequest{
			Name:        "new_test_name",
			Description: "new_test_desc",
			Price:       10.99,
			Categories:  []string{"new_test_name"},
		}
		jsonData, _ := json.Marshal(updatedRequestBody)
		resp, _ := sendRequest(http.MethodPatch, serverURL+"/product/"+strconv.Itoa(createdProduct.ID), jsonData)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var deleteResponse map[string]string
		err := json.NewDecoder(resp.Body).Decode(&deleteResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Product updated successfully", deleteResponse["message"])

		resp, _ = http.Get(serverURL + "/product/" + strconv.Itoa(createdProduct.ID))
		defer resp.Body.Close()

		var updatedProduct models.Product
		err = json.NewDecoder(resp.Body).Decode(&updatedProduct)
		assert.NoError(t, err)
		assert.Equal(t, newName, updatedProduct.Name)
		assert.Equal(t, newDescription, updatedProduct.Description)
		assert.Equal(t, newPrice, updatedProduct.Price)
		assert.Equal(t, 1, len(updatedProduct.Categories))
		assert.Equal(t, []models.Category{{ID: 1, Name: "new_test_name", Description: "new_test_desc"}}, updatedProduct.Categories)
	})

	t.Run("Delete product", func(t *testing.T) {
		resp, _ := sendRequest(http.MethodDelete, serverURL+"/product/"+strconv.Itoa(createdProduct.ID), nil)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var deleteResponse map[string]string
		err := json.NewDecoder(resp.Body).Decode(&deleteResponse)
		assert.NoError(t, err)
		assert.Equal(t, "Product deleted successfully", deleteResponse["message"])
	})

	t.Run("Get deleted product", func(t *testing.T) {
		resp, _ := http.Get(serverURL + "/product/" + strconv.Itoa(createdProduct.ID))
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
