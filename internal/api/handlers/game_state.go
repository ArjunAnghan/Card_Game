package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"

	"github.com/gorilla/mux"
)

// GetRemainingCardsCountBySuitHandler handles the HTTP request to get the count of how many cards
// per suit are left undealt in the game deck. The counts for each suit are returned as a JSON response.
func GetRemainingCardsCountBySuitHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Retrieve the count of remaining cards per suit
		suitCounts, err := gameService.GetRemainingCardsCountBySuit(gameID)
		if err != nil {
			// Return a 500 Internal Server Error status if retrieving the counts fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the suit counts as JSON and write it to the response
		json.NewEncoder(w).Encode(suitCounts)
	}
}

// GetRemainingCardsSortedHandler handles the HTTP request to get the count of each card (suit and value)
// remaining in the game deck, sorted by suit (hearts, spades, clubs, diamonds) and face value from high
// value to low value (King, Queen, Jack, 10….2, Ace with value of 1). The sorted counts are returned as a JSON response.
func GetRemainingCardsSortedHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Retrieve the remaining cards sorted by suit and value
		remainingCards, err := gameService.GetRemainingCardsSorted(gameID)
		if err != nil {
			// Return a 500 Internal Server Error status if retrieving the sorted cards fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the sorted remaining cards as JSON and write it to the response
		json.NewEncoder(w).Encode(remainingCards)
	}
}
