package models

// Deck represents a deck of playing cards.
// It contains a slice of Card structs, representing the cards in the deck.
type Deck struct {
	Cards []Card `json:"cards"`
}

// NewDeck initializes a new deck of 52 cards.
// The deck contains cards from all four suits (Hearts, Diamonds, Clubs, Spades)
// and thirteen face values (Ace, 2-10, Jack, Queen, King).
func NewDeck() *Deck {
	// Define the suits and values for a standard deck of cards
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	values := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	var cards []Card

	// Loop through each suit
	for _, suit := range suits {
		// Loop through each value
		for _, value := range values {
			// Create a new card with the current suit and value, and add it to the deck
			cards = append(cards, Card{Suit: suit, Value: value})
		}
	}

	// Return a pointer to a new Deck containing the initialized cards
	return &Deck{Cards: cards}
}
