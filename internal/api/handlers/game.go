package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateGameHandler handles the HTTP request to create a new game.
// It decodes the request payload, uses the GameService to create the game,
// and returns the newly created game as a JSON response.
func CreateGameHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define a struct to capture the incoming request payload
		var req struct {
			Name string `json:"name"`
		}

		// Decode the JSON request body into the req struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Return a 400 Bad Request status if the payload is invalid
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Create a new game using the game service
		game, err := gameService.CreateGame(req.Name)
		if err != nil {
			// Return a 500 Internal Server Error status if game creation fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the created game as JSON and write it to the response
		json.NewEncoder(w).Encode(game)
	}
}

// DeleteGameHandler handles the HTTP request to delete an existing game.
// It extracts the game ID from the URL, uses the GameService to delete the game,
// and returns an appropriate HTTP status code based on the outcome.
func DeleteGameHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Attempt to delete the game using the game service
		if err := gameService.DeleteGame(gameID); err != nil {
			// Return a 404 Not Found status if the game does not exist
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Return a 204 No Content status to indicate successful deletion
		w.WriteHeader(http.StatusNoContent)
	}
}

// AddDeckToGameHandler handles the HTTP request to add a new deck of cards to an existing game.
// It uses the DeckService to create a new deck, then adds this deck to the specified game using the GameService.
// The updated game is returned as a JSON response.
func AddDeckToGameHandler(gameService *services.GameService, deckService *services.DeckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Create a new deck using the deck service
		deck := deckService.CreateDeck()

		// Add the new deck to the specified game using the game service
		game, err := gameService.AddDeckToGame(gameID, deck)
		if err != nil {
			// Return a 500 Internal Server Error status if adding the deck to the game fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the updated game as JSON and write it to the response
		json.NewEncoder(w).Encode(game)
	}
}

// ShuffleGameDeckHandler handles the HTTP request to shuffle the game deck.
// It extracts the game ID from the URL, uses the GameService to shuffle the deck,
// and returns an appropriate HTTP status code.
func ShuffleGameDeckHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Attempt to shuffle the game deck using the game service
		err := gameService.ShuffleGameDeck(gameID)
		if err != nil {
			// Return a 500 Internal Server Error status if shuffling fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Return a 200 OK status to indicate successful shuffling
		w.WriteHeader(http.StatusOK)
	}
}

// DealCardToPlayerHandler handles the HTTP request to deal a card to a specific player in a game.
// It decodes the request payload to get the player's name, uses the GameService to deal a card,
// and returns the dealt card as a JSON response.
func DealCardToPlayerHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Define a struct to capture the incoming request payload
		var req struct {
			PlayerName string `json:"player_name"`
		}

		// Decode the JSON request body into the req struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Return a 400 Bad Request status if the payload is invalid
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Deal a card to the specified player using the game service
		card, err := gameService.DealCardToPlayer(gameID, req.PlayerName)
		if err != nil {
			// Return a 500 Internal Server Error status if dealing the card fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the dealt card as JSON and write it to the response
		json.NewEncoder(w).Encode(card)
	}
}
