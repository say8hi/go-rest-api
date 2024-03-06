package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
)
// Users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  var request_user models.CreateUserRequest
  err := json.NewDecoder(r.Body).Decode(&request_user)
  if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
  
  _, err = database.GetUserByUsername(request_user.Username)
  if err == nil { 
      http.Error(w, "This username is already taken.", http.StatusConflict)
      return
  } else if err != sql.ErrNoRows { 
      http.Error(w, "Database error.", http.StatusInternalServerError)
      return
  }

  hash := sha256.Sum256([]byte(request_user.Password + request_user.Username))
	hashedPassword := hex.EncodeToString(hash[:])
  
  request_user.Password = string(hashedPassword)
  user, err := database.CreateUser(request_user)
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
  }
    
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

