package main

import (
	"log"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/handlers"
)

func main() {
	database.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users/create", handlers.CreateUserHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
