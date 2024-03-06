package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
)


func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

  var requestCategory models.CreateCategoryRequest
  err := json.NewDecoder(r.Body).Decode(&requestCategory)
  if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdCategory, err := database.CreateCategory(requestCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

