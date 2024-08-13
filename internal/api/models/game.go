package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Players  []string           `bson:"players" json:"players"` // This can be a slice of player IDs
	GameDeck []Card             `bson:"game_deck" json:"game_deck"`
}

type Card struct {
	Suit  string `bson:"suit" json:"suit"`
	Value string `bson:"value" json:"value"`
}
