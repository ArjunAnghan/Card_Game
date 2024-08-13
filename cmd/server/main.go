package main

import (
	"log"
	"my-card-game/internal/api"
	"my-card-game/internal/config"
	"my-card-game/internal/db"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	db.ConnectDB(cfg) // Ensure this is called first
	//defer db.DisconnectDB()

	//Initialize the router
	r := mux.NewRouter()

	// Register routes
	api.RegisterRoutes(r)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
