package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
)


func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

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
