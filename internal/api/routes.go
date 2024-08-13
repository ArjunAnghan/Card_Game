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
	r.HandleFunc("/games/{id}/add-deck", handlers.AddDeckToGameHandler(gameService, deckService)).Methods("POST")
	r.HandleFunc("/games/{id}/add-player", handlers.AddPlayerHandler(gameService)).Methods("POST")
	r.HandleFunc("/games/{id}/remove-player", handlers.RemovePlayerHandler(gameService)).Methods("POST")
	r.HandleFunc("/games/{id}/shuffle", handlers.ShuffleGameDeckHandler(gameService)).Methods("POST")
	r.HandleFunc("/games/{id}/deal-card", handlers.DealCardToPlayerHandler(gameService)).Methods("POST")
	r.HandleFunc("/games/{id}/player-hand", handlers.GetPlayerHandHandler(gameService)).Methods("GET")
	r.HandleFunc("/games/{id}/player-hand-values", handlers.GetPlayersWithHandValuesHandler(gameService)).Methods("GET")
	r.HandleFunc("/games/{id}/remaining-cards-suit-count", handlers.GetRemainingCardsCountBySuitHandler(gameService)).Methods("GET")

	// Add other routes here...
}
