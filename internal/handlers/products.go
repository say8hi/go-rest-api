package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
	"github.com/say8hi/go-api-test/internal/utils"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productRequest models.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdProduct, err := database.CreateProduct(productRequest)
	if err == database.ErrCategoryDoesntExists {
		utils.SendJSONError(w, "One or more of the categories you specified doesn't exist", http.StatusBadRequest)
		return
	} else if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func GetAllProductsInCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.SendJSONError(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	products, err := database.GetProductsByCategory(categoryID)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.SendJSONError(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var requestProduct models.ProductUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&requestProduct)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateProduct(productID, requestProduct)
	if err == database.ErrCategoryDoesntExists {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	} else if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
	}

	response := models.GeneralResponse{
		Status:  "success",
		Message: "Product updated successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.SendJSONError(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendJSONError(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = database.DeleteProduct(productID)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.GeneralResponse{
		Status:  "success",
		Message: "Product deleted successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		utils.SendJSONError(w, "ID is missing in parameters", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(idStr)
	if err == sql.ErrNoRows {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product, err := database.GetProduct(productID)
	if err == sql.ErrNoRows {
		utils.SendJSONError(w, "product not found", http.StatusNotFound)
		return
	} else if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
