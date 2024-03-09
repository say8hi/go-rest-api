package middlewares

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		receivedHash := r.Header.Get("Authorization")
		if !strings.HasPrefix(receivedHash, "Bearer ") {
			utils.SendJSONError(w, "Bad Authorization header format.", http.StatusBadRequest)
			return
		}

		receivedHash = strings.TrimPrefix(receivedHash, "Bearer ")
		_, err := database.GetUserByPasswordHash(receivedHash)
		if err == sql.ErrNoRows {
			utils.SendJSONError(w, "Unauthorized", http.StatusConflict)
			return
		} else if err != nil {
			utils.SendJSONError(w, "Database error", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
