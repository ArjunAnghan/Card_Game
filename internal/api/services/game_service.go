package services

import (
	"context"
	"errors"
	"my-card-game/internal/api/models"
	"my-card-game/internal/db"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GameService provides services related to game operations.
// It interacts with the MongoDB collection where game data is stored.
type GameService struct {
	collection *mongo.Collection
}

// NewGameService creates and returns a new instance of GameService.
// It initializes the service with a reference to the MongoDB collection where game data is stored.
func NewGameService() *GameService {
	return &GameService{
		collection: db.GetCollection("games"),
	}
}

// CreateGame creates a new game with the given name.
// It initializes the game with a unique ID, an empty list of players, and an empty game deck.
// The game is then inserted into the MongoDB collection, and the created game is returned.
func (s *GameService) CreateGame(name string) (*models.Game, error) {
	// Create a context with a timeout of 5 seconds to manage the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize a new game with a unique ID, the provided name, no players, and an empty deck
	game := &models.Game{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Players:  []string{},
		GameDeck: []models.Card{}, // Initialize with an empty deck
	}

	// Insert the new game into the MongoDB collection
	_, err := s.collection.InsertOne(ctx, game)
	if err != nil {
		// Return an error if the insertion fails
		return nil, err
	}

	// Return the created game
	return game, nil
}

// DeleteGame deletes an existing game by its ID.
// The game ID is converted from a hex string to an ObjectID, and the corresponding game is deleted from the collection.
// If the game is not found or the ID is invalid, an error is returned.
func (s *GameService) DeleteGame(id string) error {
	// Create a context with a timeout of 5 seconds to manage the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the game ID from a hex string to an ObjectID
	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Return an error if the game ID is invalid
		return errors.New("invalid game ID")
	}

	// Attempt to delete the game from the MongoDB collection
	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": gameID})
	if err != nil {
		// Return an error if the deletion fails
		return err
	}

	// Check if any document was deleted; if not, return an error indicating the game was not found
	if result.DeletedCount == 0 {
		return errors.New("game not found")
	}

	// Return nil if the deletion was successful
	return nil
}
