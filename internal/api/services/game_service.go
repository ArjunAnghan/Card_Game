package services

import (
	"context"
	"errors"
	"my-card-game/internal/api/models"
	"my-card-game/internal/db"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerHandValue struct {
	PlayerName string `json:"player_name"`
	HandValue  int    `json:"hand_value"`
}

type GameService struct {
	collection *mongo.Collection
}

type SuitCount struct {
	Suit  string `json:"suit"`
	Count int    `json:"count"`
}

type CardCount struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
	Count int    `json:"count"`
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

func (s *GameService) ShuffleGameDeck(gameID string) (*models.Game, error) {
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

	game.ShuffleDeck()

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck},
	})
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *GameService) DealCardToPlayer(gameID, playerName string) (*models.Card, error) {
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

	if len(game.GameDeck) == 0 {
		return nil, errors.New("no cards left to deal")
	}

	// Deal the top card from the deck
	dealtCard := game.GameDeck[0]
	game.GameDeck = game.GameDeck[1:]

	// Add the card to the player's hand
	if game.PlayerHands == nil {
		game.PlayerHands = make(map[string][]models.Card)
	}
	game.PlayerHands[playerName] = append(game.PlayerHands[playerName], dealtCard)

	// Update the game state in the database
	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": gameIDObj}, bson.M{
		"$set": bson.M{"game_deck": game.GameDeck, "player_hands": game.PlayerHands},
	})
	if err != nil {
		return nil, err
	}

	return &dealtCard, nil
}

func (s *GameService) GetPlayerHand(gameID, playerName string) ([]models.Card, error) {
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

	hand, exists := game.PlayerHands[playerName]
	if !exists {
		return nil, errors.New("player not found or no cards dealt to this player")
	}

	return hand, nil
}

func (s *GameService) GetPlayersWithHandValues(gameID string) ([]PlayerHandValue, error) {
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

	// Calculate the hand value for each player
	playerHandValues := []PlayerHandValue{}
	for player, hand := range game.PlayerHands {
		totalValue := 0
		for _, card := range hand {
			totalValue += s.getCardValue(card)
		}
		playerHandValues = append(playerHandValues, PlayerHandValue{
			PlayerName: player,
			HandValue:  totalValue,
		})
	}

	// Sort the players by hand value in descending order
	sort.Slice(playerHandValues, func(i, j int) bool {
		return playerHandValues[i].HandValue > playerHandValues[j].HandValue
	})

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

func (s *GameService) GetRemainingCardsCountBySuit(gameID string) ([]SuitCount, error) {
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

	return remainingCounts, nil
}

func (s *GameService) GetRemainingCardsSorted(gameID string) ([]CardCount, error) {
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
	suitsOrder := []string{"Hearts", "Spades", "Clubs", "Diamonds"}
	valuesOrder := []string{"King", "Queen", "Jack", "10", "9", "8", "7", "6", "5", "4", "3", "2", "Ace"}

	for _, suit := range suitsOrder {
		for _, value := range valuesOrder {
			count := cardCounts[suit][value]
			if count > 0 {
				remainingCards = append(remainingCards, CardCount{
					Suit:  suit,
					Value: value,
					Count: count,
				})
			}
		}
	}

	return remainingCards, nil
}
