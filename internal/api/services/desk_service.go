package services

import "my-card-game/internal/api/models"

type DeckService struct{}

func NewDeckService() *DeckService {
	return &DeckService{}
}

func (ds *DeckService) CreateDeck() *models.Deck {
	return models.NewDeck()
}
