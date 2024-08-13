package models

type Deck struct {
	Cards []Card `json:"cards"`
}

// NewDeck initializes a new deck of 52 cards
func NewDeck() *Deck {
	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
	values := []string{"Ace", "2", "3", "4", "5", "6", "7", "8", "9", "10", "Jack", "Queen", "King"}

	var cards []Card

	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, Card{Suit: suit, Value: value})
		}
	}

	return &Deck{Cards: cards}
}
