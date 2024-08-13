package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateGameHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, err := gameService.CreateGame(req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func DeleteGameHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		if err := gameService.DeleteGame(gameID); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func AddDeckToGameHandler(gameService *services.GameService, deckService *services.DeckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		// Create a new deck to be added to the game
		deck := deckService.CreateDeck()

		game, err := gameService.AddDeckToGame(gameID, deck)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}
func AddPlayerHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		var req struct {
			PlayerName string `json:"player_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, err := gameService.AddPlayer(gameID, req.PlayerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func RemovePlayerHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		var req struct {
			PlayerName string `json:"player_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		game, err := gameService.RemovePlayer(gameID, req.PlayerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func ShuffleGameDeckHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		game, err := gameService.ShuffleGameDeck(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(game)
	}
}

func DealCardToPlayerHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		var req struct {
			PlayerName string `json:"player_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		card, err := gameService.DealCardToPlayer(gameID, req.PlayerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}
}
func GetPlayerHandHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]
		playerName := r.URL.Query().Get("player_name")

		if playerName == "" {
			http.Error(w, "player_name is required", http.StatusBadRequest)
			return
		}

		hand, err := gameService.GetPlayerHand(gameID, playerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hand)
	}
}

func GetPlayersWithHandValuesHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		playerHandValues, err := gameService.GetPlayersWithHandValues(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(playerHandValues)
	}
}

func GetRemainingCardsCountBySuitHandler(gameService *services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		gameID := vars["id"]

		suitCounts, err := gameService.GetRemainingCardsCountBySuit(gameID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(suitCounts)
	}
}
