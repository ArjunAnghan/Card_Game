package services

import (
	"context"
	"errors"
	"my-card-game/internal/api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SuitCount represents the count of remaining cards for a specific suit.
// It includes the suit name and the count of cards remaining.
type SuitCount struct {
	Suit  string `json:"suit"`
	Count int    `json:"count"`
}

// CardCount represents the count of remaining cards for a specific suit and value.
// It includes the suit, value, and the count of cards remaining.
type CardCount struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
	Count int    `json:"count"`
}

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

// AddDeckToGame adds a new deck of cards to an existing game's deck.
// It finds the game by its ID, appends the new deck to the game's deck,
// and updates the game document in the MongoDB collection.
func (s *GameService) AddDeckToGame(gameID string, deck *models.Deck) (*models.Game, error) {
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

	// Append the new deck to the existing game deck
	game.GameDeck = append(game.GameDeck, deck.Cards...)

	// Update the game document in the MongoDB collection with the new deck
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck},
	})
	if err != nil {
		// Return an error if the update operation fails
		return nil, err
	}

	// Return the updated game object
	return &game, nil
}

// Shuffle the Deck
func (s *GameService) ShuffleGameDeck(gameID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gameIDObj, err := primitive.ObjectIDFromHex(gameID)
	if err != nil {
		return errors.New("invalid game ID")
	}

	var game models.Game
	err = s.collection.FindOne(ctx, bson.M{"_id": gameIDObj}).Decode(&game)
	if err != nil {
		return errors.New("game not found")
	}

	// Shuffle the game deck
	game.ShuffleDeck()

	// Update the game state in the database
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck},
	})
	if err != nil {
		return err
	}

	return nil
}

// GetRemainingCardsCountBySuit retrieves the count of remaining cards for each suit in a game.
// The function returns a list of SuitCount objects, each representing the count of remaining cards for a specific suit.
func (s *GameService) GetRemainingCardsCountBySuit(gameID string) ([]SuitCount, error) {
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

	// Initialize counters for each suit
	suitCounts := map[string]int{
		"Hearts":   0,
		"Diamonds": 0,
		"Clubs":    0,
		"Spades":   0,
	}

	// Count the number of cards left for each suit
	for _, card := range game.GameDeck {
		suitCounts[card.Suit]++
	}

	// Convert the map to a slice of SuitCount
	remainingCounts := []SuitCount{}
	for suit, count := range suitCounts {
		remainingCounts = append(remainingCounts, SuitCount{
			Suit:  suit,
			Count: count,
		})
	}

	// Return the list of SuitCount objects
	return remainingCounts, nil
}

// GetRemainingCardsSorted retrieves the count of each card (suit and value) remaining in the game deck,
// sorted by suit (Hearts, Spades, Clubs, Diamonds) and face value from high value to low value (King, Queen, Jack, etc.).
// The function returns a list of CardCount objects representing the sorted remaining cards.
func (s *GameService) GetRemainingCardsSorted(gameID string) ([]CardCount, error) {
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

	// Initialize a map to count the cards
	cardCounts := map[string]map[string]int{
		"Hearts":   {},
		"Diamonds": {},
		"Clubs":    {},
		"Spades":   {},
	}

	// Count the remaining cards in the game deck
	for _, card := range game.GameDeck {
		cardCounts[card.Suit][card.Value]++
	}

	// Convert the map to a slice of CardCount and sort it
	remainingCards := []CardCount{}
	// Define the order of suits and values for sorting
	suitsOrder := []string{"Hearts", "Spades", "Clubs", "Diamonds"}
	valuesOrder := []string{"King", "Queen", "Jack", "10", "9", "8", "7", "6", "5", "4", "3", "2", "Ace"}

	// Iterate over the suits and values in the specified order
	for _, suit := range suitsOrder {
		for _, value := range valuesOrder {
			// Get the count of the current suit and value
			count := cardCounts[suit][value]
			if count > 0 {
				// Add the suit, value, and count to the remainingCards slice
				remainingCards = append(remainingCards, CardCount{
					Suit:  suit,
					Value: value,
					Count: count,
				})
			}
		}
	}

	// Return the sorted list of remaining cards
	return remainingCards, nil
}
