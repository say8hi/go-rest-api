package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/models"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  var request_user models.CreateUserRequest
  json.NewDecoder(r.Body).Decode(&request_user)
  
hash := sha256.Sum256([]byte(*&request_user.Password))
	hashedPassword := hex.EncodeToString(hash[:])
  
  fmt.Println(string(hashedPassword))
  request_user.Password = string(hashedPassword)
  user, err := database.CreateUser(request_user)
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
  }
    
	json.NewEncoder(w).Encode(user)
}
