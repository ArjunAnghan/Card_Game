package services

import (
	"context"
	"errors"
	"my-card-game/internal/api/models"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlayerHandValue represents the total value of a player's hand.
// It includes the player's name and the total hand value.
type PlayerHandValue struct {
	PlayerName string `json:"player_name"`
	HandValue  int    `json:"hand_value"`
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

// DealCardToPlayer deals a card from the game's deck to the specified player.
// The top card from the game deck is removed and added to the player's hand.
// The updated game state is then saved to the database.
func (s *GameService) DealCardToPlayer(gameID, playerName string) (*models.Card, error) {
	// Create a context with a timeout of 5 seconds to manage the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the game ID from a hex string to an ObjectID
	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		// Return an error if the game ID is invalid
		return nil, errors.New("invalid game ID")
	}

	// Find the game in the MongoDB collection using the provided game ID
	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		// Return an error if the game is not found
		return nil, errors.New("game not found")
	}

	// Check if there are any cards left to deal
	if len(game.GameDeck) == 0 {
		// Return an error if there are no cards left in the deck
		return nil, errors.New("no cards left to deal")
	}

	// Deal the top card from the deck
	dealtCard := game.GameDeck[0]
	// Remove the dealt card from the game deck
	game.GameDeck = game.GameDeck[1:]

	// Initialize the player hands map if it hasn't been already
	if game.PlayerHands == nil {
		game.PlayerHands = make(map[string][]models.Card)
	}
	// Add the dealt card to the player's hand
	game.PlayerHands[playerName] = append(game.PlayerHands[playerName], dealtCard)

	// Update the game state in the database
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck, "player_hands": game.PlayerHands},
	})
	if err != nil {
		// Return an error if the update operation fails
		return nil, err
	}

	// Return the dealt card
	return &dealtCard, nil
}

// GetPlayerHand retrieves the list of cards held by a specific player in a game.
// It finds the game by its ID, checks if the player has any cards dealt,
// and returns the player's hand or an error if the game or player is not found.
func (s *GameService) GetPlayerHand(gameID, playerName string) ([]models.Card, error) {
	// Create a context with a timeout of 5 seconds to manage the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the game ID from a hex string to an ObjectID
	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		// Return an error if the game ID is invalid
		return nil, errors.New("invalid game ID")
	}

	// Find the game in the MongoDB collection using the provided game ID
	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		// Return an error if the game is not found
		return nil, errors.New("game not found")
	}

	// Retrieve the player's hand from the game's PlayerHands map
	hand, exists := game.PlayerHands[playerName]
	if !exists {
		// Return an error if the player is not found or has no cards dealt
		return nil, errors.New("player not found or no cards dealt to this player")
	}

	// Return the player's hand
	return hand, nil
}

// GetPlayersWithHandValues retrieves the list of players in a game along with the total value of their hands.
// The players are sorted in descending order based on the value of their hands, and the sorted list is returned.
func (s *GameService) GetPlayersWithHandValues(gameID string) ([]PlayerHandValue, error) {
	// Create a context with a timeout of 5 seconds to manage the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert the game ID from a hex string to an ObjectID
	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		// Return an error if the game ID is invalid
		return nil, errors.New("invalid game ID")
	}

	// Find the game in the MongoDB collection using the provided game ID
	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		// Return an error if the game is not found
		return nil, errors.New("game not found")
	}

	// Calculate the hand value for each player
	playerHandValues := []PlayerHandValue{}
	for player, hand := range game.PlayerHands {
		totalValue := 0
		for _, card := range hand {
			// Add the value of each card to the player's total hand value
			totalValue += s.getCardValue(card)
		}
		// Append the player's name and hand value to the playerHandValues slice
		playerHandValues = append(playerHandValues, PlayerHandValue{
			PlayerName: player,
			HandValue:  totalValue,
		})
	}

	// Sort the players by hand value in descending order
	sort.Slice(playerHandValues, func(i, j int) bool {
		return playerHandValues[i].HandValue > playerHandValues[j].HandValue
	})

	// Return the sorted list of players with their hand values
	return playerHandValues, nil
}

// Helper function to get the value of a card
func (s *GameService) getCardValue(card models.Card) int {
	switch card.Value {
	case "Ace":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "Jack":
		return 11
	case "Queen":
		return 12
	case "King":
		return 13
	default:
		return 0
	}
}
