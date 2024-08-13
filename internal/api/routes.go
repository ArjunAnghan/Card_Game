package api

import (
	"my-card-game/internal/api/handlers"
	"my-card-game/internal/api/services"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	// Initialize services here instead of as global variables
	gameService := services.NewGameService()
	deckService := services.NewDeckService()

	r.HandleFunc("/games", handlers.CreateGameHandler(gameService)).Methods("POST")
	r.HandleFunc("/games/{id}", handlers.DeleteGameHandler(gameService)).Methods("DELETE")
	r.HandleFunc("/decks", handlers.CreateDeckHandler(deckService)).Methods("POST")

	// Add other routes here...
}
