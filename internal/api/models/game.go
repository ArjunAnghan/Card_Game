package models

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Game represents a card game.
// It includes an ID, a name, a list of players, the game deck (cards available in the game),
// and a map to track the cards held by each player.
type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Players     []string           `bson:"players" json:"players"` // This can be a slice of player IDs
	GameDeck    []Card             `bson:"game_deck" json:"game_deck"`
	PlayerHands map[string][]Card  `bson:"player_hands" json:"player_hands"`
}

// Card represents an individual playing card.
// It includes the suit and value of the card.
type Card struct {
	Suit  string `bson:"suit" json:"suit"`
	Value string `bson:"value" json:"value"`
}

// AddDeckToGame adds a deck of cards to the game's deck.
// The new deck is appended to the existing game deck.
func (g *Game) AddDeckToGame(deck *Deck) {
	g.GameDeck = append(g.GameDeck, deck.Cards...)
}

// ShuffleDeck shuffles the cards in the game deck using a custom shuffle algorithm.
// The cards are shuffled in place using a random number generator.
func (g *Game) ShuffleDeck() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator with the current time
	n := len(g.GameDeck)
	for i := range g.GameDeck {
		j := rand.Intn(n)                                           // Generate a random index between 0 and n-1
		g.GameDeck[i], g.GameDeck[j] = g.GameDeck[j], g.GameDeck[i] // Swap the card at index i with the card at index j
	}
}
