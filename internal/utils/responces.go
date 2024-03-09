package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/say8hi/go-api-test/internal/models"
)

func SendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errorResponse := models.GeneralResponse{
		Status:  "error",
		Message: message,
	}
	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshalling error response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
