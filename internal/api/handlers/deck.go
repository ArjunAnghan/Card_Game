package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"
)

// CreateDeckHandler handles the HTTP request to create a new deck of cards.
// It uses the DeckService to generate a new deck and returns it as a JSON response.
func CreateDeckHandler(deckService *services.DeckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new deck using the deck service
		deck := deckService.CreateDeck()

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the deck as JSON and write it to the response
		json.NewEncoder(w).Encode(deck)
	}
}
