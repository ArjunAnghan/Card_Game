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

type GameService struct {
	collection *mongo.Collection
}

func NewGameService() *GameService {
	return &GameService{
		collection: db.GetCollection("games"),
	}
}

func (s *GameService) CreateGame(name string) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	game := &models.Game{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Players:  []string{},
		GameDeck: []models.Card{}, // Initialize with an empty deck
	}

	_, err := s.collection.InsertOne(ctx, game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func (s *GameService) DeleteGame(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid game ID")
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": gameID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("game not found")
	}

	return nil
}
func (s *GameService) AddDeckToGame(gameID string, deck *models.Deck) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return nil, errors.New("invalid game ID")
	}

	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		return nil, errors.New("game not found")
	}

	// Append the new deck to the existing game deck
	game.GameDeck = append(game.GameDeck, deck.Cards...)

	// Update the game document in the database
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck},
	})
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// AddPlayer adds a player to a game
func (s *GameService) AddPlayer(gameID, playerName string) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return nil, errors.New("invalid game ID")
	}

	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		return nil, errors.New("game not found")
	}

	// Add the player to the game if they are not already in it
	for _, player := range game.Players {
		if player == playerName {
			return nil, errors.New("player already in the game")
		}
	}
	game.Players = append(game.Players, playerName)

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"players": game.Players},
	})
	if err != nil {
		return nil, err
	}

	return &game, nil
}

// RemovePlayer removes a player from a game
func (s *GameService) RemovePlayer(gameID, playerName string) (*models.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return nil, errors.New("invalid game ID")
	}

	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		return nil, errors.New("game not found")
	}

	// Remove the player from the game
	newPlayers := []string{}
	for _, player := range game.Players {
		if player != playerName {
			newPlayers = append(newPlayers, player)
		}
	}

	// If the player was not found, return an error
	if len(newPlayers) == len(game.Players) {
		return nil, errors.New("player not found in the game")
	}

	game.Players = newPlayers

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"players": game.Players},
	})
	if err != nil {
		return nil, err
	}

	return &game, nil
}
