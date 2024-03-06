package main

import (
	"log"
	"net/http"

	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/handlers"
	"github.com/say8hi/go-api-test/internal/middlewares"
)

func main() {
	database.Init()
  database.CreateTables()

	mux := http.NewServeMux()
  
  // Users
  mux.HandleFunc("POST /users/create", handlers.CreateUserHandler)
	
  // Categories
  mux.HandleFunc("POST /category/create",
    func(w http.ResponseWriter, r *http.Request) {
    middleware := middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateCategoryHandler))
    middleware.ServeHTTP(w, r)
})

  // Products
  mux.HandleFunc("POST /product/create",
    func(w http.ResponseWriter, r *http.Request) {
    middleware := middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateProductHandler))
    middleware.ServeHTTP(w, r)
})
  
  defer database.CloseConnection()	
  log.Fatal(http.ListenAndServe(":8080", mux))
}
