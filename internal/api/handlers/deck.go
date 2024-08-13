package handlers

import (
	"encoding/json"
	"my-card-game/internal/api/services"
	"net/http"
)

func CreateDeckHandler(deckService *services.DeckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deck := deckService.CreateDeck()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(deck)
	}
}
