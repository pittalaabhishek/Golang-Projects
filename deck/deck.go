// deck.go
package deck

import (
	"math/rand"
	"sort"
	"time"
)

type Deck []Card

func New(options ...func([]Card) []Card) []Card {
	var cards []Card
	
	suits := []Suit{Spade, Diamond, Club, Heart}
	
	for _, suit := range suits {
		for rank := Ace; rank <= King; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	
	// Apply all options
	for _, option := range options {
		cards = option(cards)
	}
	
	return cards
}

func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, func(i, j int) bool {
		if cards[i].Suit != cards[j].Suit {
			return cards[i].Suit < cards[j].Suit
		}
		return cards[i].Rank < cards[j].Rank
	})
	return cards
}

func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

func Shuffle(cards []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
	return cards
}

func WithJokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, NewJoker())
		}
		return cards
	}
}

func Filter(keep func(Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var filtered []Card
		for _, card := range cards {
			if keep(card) {
				filtered = append(filtered, card)
			}
		}
		return filtered
	}
}

func WithoutRanks(ranks ...Rank) func([]Card) []Card {
	rankMap := make(map[Rank]bool)
	for _, rank := range ranks {
		rankMap[rank] = true
	}
	
	return Filter(func(card Card) bool {
		return !rankMap[card.Rank]
	})
}

func WithoutSuits(suits ...Suit) func([]Card) []Card {
	suitMap := make(map[Suit]bool)
	for _, suit := range suits {
		suitMap[suit] = true
	}
	
	return Filter(func(card Card) bool {
		return !suitMap[card.Suit]
	})
}

func WithDecks(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var result []Card
		for i := 0; i < n; i++ {
			result = append(result, cards...)
		}
		return result
	}
}