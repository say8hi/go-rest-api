package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
	"github.com/say8hi/go-api-test/internal/utils"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var requestCategory models.CreateCategoryRequest
	err := json.NewDecoder(r.Body).Decode(&requestCategory)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdCategory, err := database.CreateCategory(requestCategory)
  if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
      utils.SendJSONError(w, "This category name is already exist.", http.StatusBadRequest)
    return
  } else if err != nil{
      utils.SendJSONError(w, "Database error.", http.StatusInternalServerError)
    return
  }

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	var requestCategory models.CategoryUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&requestCategory)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateCategory(categoryID, requestCategory)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

  response := models.GeneralResponse{
      Status:  "success",
      Message: "Category updated successfully",
  }

  jsonResponse, err := json.Marshal(response)
  if err != nil {
      utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
      return
  }

	w.WriteHeader(http.StatusOK)
  w.Write(jsonResponse)
}

func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
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

	err = database.DeleteCategory(categoryID)
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

  response := models.GeneralResponse{
      Status:  "success",
      Message: "Category deleted successfully",
  }

  jsonResponse, err := json.Marshal(response)
  if err != nil {
      utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
      return
  }

	w.WriteHeader(http.StatusOK)
  w.Write(jsonResponse)
}

func GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
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

	category, err := database.GetCategoryByID(categoryID)
  if err == sql.ErrNoRows{
		utils.SendJSONError(w, "category not found", http.StatusNotFound)
		return
  } else if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(category)
}

func GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetAllCategories()
	if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)
}
