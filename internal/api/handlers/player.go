package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"

	"github.com/gorilla/mux"
)

// AddPlayerHandler handles the HTTP request to add a player to a game.
// It decodes the request payload to get the player's name and uses the GameService
// to add the player to the specified game. The updated game is returned as a JSON response.
func AddPlayerHandler(gameService *services.GameService) http.HandlerFunc {
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

		// Add the player to the specified game using the game service
		game, err := gameService.AddPlayer(gameID, req.PlayerName)
		if err != nil {
			// Return a 500 Internal Server Error status if adding the player fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the updated game as JSON and write it to the response
		json.NewEncoder(w).Encode(game)
	}
}

// RemovePlayerHandler handles the HTTP request to remove a player from a game.
// It decodes the request payload to get the player's name and uses the GameService
// to remove the player from the specified game. The updated game is returned as a JSON response.
func RemovePlayerHandler(gameService *services.GameService) http.HandlerFunc {
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

		// Remove the player from the specified game using the game service
		game, err := gameService.RemovePlayer(gameID, req.PlayerName)
		if err != nil {
			// Return a 500 Internal Server Error status if removing the player fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the updated game as JSON and write it to the response
		json.NewEncoder(w).Encode(game)
	}
}

// GetPlayerHandHandler handles the HTTP request to get the list of cards held by a specific player in a game.
// It extracts the player's name from the query parameters, uses the GameService to retrieve the player's hand,
// and returns the list of cards as a JSON response.
func GetPlayerHandHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Get the player's name from the query parameters
		playerName := r.URL.Query().Get("player_name")

		// Check if the player's name is provided in the query parameters
		if playerName == "" {
			// Return a 400 Bad Request status if the player name is not provided
			http.Error(w, "player_name is required", http.StatusBadRequest)
			return
		}

		// Get the player's hand using the game service
		hand, err := gameService.GetPlayerHand(gameID, playerName)
		if err != nil {
			// Return a 500 Internal Server Error status if retrieving the hand fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the player's hand as JSON and write it to the response
		json.NewEncoder(w).Encode(hand)
	}
}

// GetPlayersWithHandValuesHandler handles the HTTP request to get the list of players in a game
// along with the total value of all the cards each player holds. The list is sorted in descending order
// based on the hand values. The sorted list is returned as a JSON response.
func GetPlayersWithHandValuesHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the game ID from the URL path variables
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Retrieve the list of players with their hand values, sorted in descending order
		playerHandValues, err := gameService.GetPlayersWithHandValues(gameID)
		if err != nil {
			// Return a 500 Internal Server Error status if retrieving the hand values fails
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the response header to indicate JSON content
		w.Header().Set("Content-Type", "application/json")

		// Encode the list of players with hand values as JSON and write it to the response
		json.NewEncoder(w).Encode(playerHandValues)
	}
}
