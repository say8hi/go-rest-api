package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/say8hi/go-api-test/internal/database"
	"github.com/say8hi/go-api-test/internal/handlers"
	"github.com/say8hi/go-api-test/internal/middlewares"
	"github.com/say8hi/go-api-test/internal/rabbitmq"
)

func main() {
	database.Init()
	database.CreateTables()
  defer database.CloseConnection()
  
  rabbitMQChannel := rabbitmq.InitRabbitMQ()
  defer rabbitMQChannel.Close()

	r := mux.NewRouter()
	r.Use(middlewares.LoggingMiddleware)

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(middlewares.AuthMiddleware)

	// Unauthorized endpoints
	// Users
	r.HandleFunc("/users/create", handlers.CreateUserHandler).Methods("POST")

	// Categories
	r.HandleFunc("/category/{id:[0-9]+}", handlers.GetCategoryByIDHandler).Methods("GET")
	r.HandleFunc("/category/", handlers.GetAllCategoriesHandler).Methods("GET")

	// Products
	r.HandleFunc("/product/{id:[0-9]+}", handlers.GetProductByIDHandler).Methods("GET")
	r.HandleFunc("/category/{id:[0-9]+}/products", handlers.GetAllProductsInCategoryHandler).Methods("GET")

	// Authorized endpoints
	// Categories
	authRouter.HandleFunc("/category/create", handlers.CreateCategoryHandler).Methods("POST")
	authRouter.HandleFunc("/category/{id:[0-9]+}", handlers.UpdateCategoryHandler).Methods("PATCH")
	authRouter.HandleFunc("/category/{id:[0-9]+}", handlers.DeleteCategoryHandler).Methods("DELETE")

	// Products
	authRouter.HandleFunc("/product/create", handlers.CreateProductHandler).Methods("POST")
	authRouter.HandleFunc("/product/{id:[0-9]+}", handlers.UpdateProductHandler).Methods("PATCH")
	authRouter.HandleFunc("/product/{id:[0-9]+}", handlers.DeleteProductHandler).Methods("DELETE")
  
  go rabbitmq.ConsumeMessages(rabbitMQChannel, "queue_from_datacollector")
	log.Fatal(http.ListenAndServe(":8080", r))
}
