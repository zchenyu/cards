package cards

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

var fullDeck *Deck

func init() {
	fullDeck = &Deck{initializeFullCards()}
	rand.Seed(time.Now().UnixNano())
}

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{}
	deck.Shuffle()
	return deck
}

// Copy returns a copy of the deck.
func (deck *Deck) Copy() *Deck {
	if deck == nil {
		return nil
	}
	deck2 := &Deck{
		cards: make([]Card, len(deck.cards)),
	}
	copy(deck2.cards, deck.cards)
	return deck2
}

// Remove removes a given card from the deck.
func (deck *Deck) Remove(c Card) {
	for i, c2 := range deck.cards {
		if c == c2 {
			deck.cards = append(deck.cards[:i], deck.cards[i+1:]...)
		}
	}
}

func (deck *Deck) Shuffle() {
	deck.cards = make([]Card, len(fullDeck.cards))
	copy(deck.cards, fullDeck.cards)
	rand.Shuffle(len(deck.cards), func(i, j int) {
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	})
}

func (deck *Deck) Draw(n int) []Card {
	cards := make([]Card, n)
	copy(cards, deck.cards[:n])
	deck.cards = deck.cards[n:]
	return cards
}

func (deck *Deck) Empty() bool {
	return len(deck.cards) == 0
}

// Scan implements the Scanner interface.
func (deck *Deck) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		if len(v)%2 != 0 {
			return errors.New("string must have even length")
		}
		deck.cards = make([]Card, len(v)/2)
		for i := 0; i < len(v); i += 2 {
			b := []byte(v[i:i+2])
			if err := deck.cards[i/2].UnmarshalText(b); err != nil {
				return fmt.Errorf("Could not unmarshal card '%s': %v", b, err)
			}
		}
		return nil
	default:
		return errors.New("must be `string` type")
	}
}

// Value implements the Valuer interface. 
func (deck *Deck) Value() (driver.Value, error) {
	return deck.String(), nil
}

func (deck *Deck) String() string {
	s := ""
	for _, c := range deck.cards {
		s += c.String()
	}
	return s
}

func initializeFullCards() []Card {
	var cards []Card

	for _, rank := range strRanks {
		for suit := range charSuitToIntSuit {
			cards = append(cards, NewCard(string(rank)+string(suit)))
		}
	}

	return cards
}
