package models

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Players     []string           `bson:"players" json:"players"` // This can be a slice of player IDs
	GameDeck    []Card             `bson:"game_deck" json:"game_deck"`
	PlayerHands map[string][]Card  `bson:"player_hands" json:"player_hands"`
}

type Card struct {
	Suit  string `bson:"suit" json:"suit"`
	Value string `bson:"value" json:"value"`
}

// AddDeckToGame adds a deck of cards to the game deck
func (g *Game) AddDeckToGame(deck *Deck) {
	g.GameDeck = append(g.GameDeck, deck.Cards...)
}

func (g *Game) ShuffleDeck() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	n := len(g.GameDeck)
	for i := range g.GameDeck {
		j := rand.Intn(n)                                           // Generate a random index
		g.GameDeck[i], g.GameDeck[j] = g.GameDeck[j], g.GameDeck[i] // Swap the cards
	}
}
