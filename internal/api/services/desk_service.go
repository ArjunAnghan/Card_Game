package services

import "my-card-game/internal/api/models"

// DeckService provides services related to deck operations.
// It serves as a layer between the application and the deck model.
type DeckService struct{}

// NewDeckService creates and returns a new instance of DeckService.
func NewDeckService() *DeckService {
	return &DeckService{}
}

// CreateDeck creates a new deck of 52 cards using the Deck model.
// It returns a pointer to the newly created deck.
func (ds *DeckService) CreateDeck() *models.Deck {
	return models.NewDeck()
}
