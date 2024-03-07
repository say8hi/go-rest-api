package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
)


func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productRequest models.CreateProductRequest
	err := json.NewDecoder(r.Body).Decode(&productRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdProduct, err := database.CreateProduct(productRequest)
	if err == database.ErrCategoryDoesntExists {
		http.Error(w, "One or more of the categories you specified doesn't exist", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
  }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func GetAllProductsInCategoryHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  idStr, ok := vars["id"]
  if !ok {
      http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
      return
  }

  categoryID, err := strconv.Atoi(idStr)
  if err != nil {
      http.Error(w, "Invalid ID format", http.StatusBadRequest)
      return
  }

  categories, err := database.GetProductsByCategory(categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  idStr, ok := vars["id"]
  if !ok {
      http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
      return
  }

  productID, err := strconv.Atoi(idStr)
  if err != nil {
      http.Error(w, "Invalid ID format", http.StatusBadRequest)
      return
  }

	var requestProduct models.ProductUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&requestProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateProduct(productID, requestProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Product updated successfully")
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  idStr, ok := vars["id"]
  if !ok {
      http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
      return
  }

  productID, err := strconv.Atoi(idStr)
  if err != nil {
      http.Error(w, "Invalid ID format", http.StatusBadRequest)
      return
  }

	err = database.DeleteProduct(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Product deleted successfully")
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  idStr, ok := vars["id"]
  if !ok {
      http.Error(w, "ID is missing in parameters", http.StatusBadRequest)
      return
  }

  productID, err := strconv.Atoi(idStr)
  if err != nil {
      http.Error(w, "Invalid ID format", http.StatusBadRequest)
      return
  }

	product, err := database.GetProduct(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

