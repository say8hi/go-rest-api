package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
	"github.com/say8hi/go-api-test/internal/utils"
)

// Users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  var request_user models.CreateUserRequest
  err := json.NewDecoder(r.Body).Decode(&request_user)
  if err != nil {
		utils.SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
  
  hash := sha256.Sum256([]byte(request_user.Password + request_user.Username))
	hashedPassword := hex.EncodeToString(hash[:])
  
  request_user.Password = string(hashedPassword)
  user, err := database.CreateUser(request_user)
  if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
      utils.SendJSONError(w, "This username is already taken.", http.StatusBadRequest)
      return
  } else if err != nil{
      utils.SendJSONError(w, "Database error.", http.StatusInternalServerError)
      return
  }
    
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

