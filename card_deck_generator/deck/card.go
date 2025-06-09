// card.go
package deck

import "fmt"

//go:generate stringer -type=Suit,Rank

// Suit represents the suit of a playing card
type Suit int

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

// Rank represents the rank of a playing card
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return "Joker"
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

func (c Card) IsJoker() bool {
	return c.Suit == Joker
}

func NewJoker() Card {
	return Card{Suit: Joker, Rank: Ace} // Rank doesn't matter for jokers
}