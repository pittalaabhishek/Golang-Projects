package main

import (
	"fmt"

	"card_deck_generator/deck"
)

func main() {
	fmt.Println("Testing Deck Package\n")

	// Test 1: Standard deck
	fmt.Println("1. Standard Deck:")
	cards := deck.New()
	fmt.Printf("   Cards in deck: %d\n", len(cards))
	fmt.Printf("   First card: %s\n", cards[0])
	fmt.Printf("   Last card: %s\n", cards[len(cards)-1])
	fmt.Printf("   Sample cards: %s, %s, %s\n\n", cards[0], cards[13], cards[26])

	// Test 2: Shuffled deck
	fmt.Println("2. Shuffled Deck:")
	shuffledCards := deck.New(deck.Shuffle)
	fmt.Printf("   First 5 cards: ")
	for i := 0; i < 5; i++ {
		fmt.Printf("%s, ", shuffledCards[i])
	}
	fmt.Println("\n")

	// Test 3: Deck with jokers
	fmt.Println("3. Deck with 2 Jokers:")
	cardsWithJokers := deck.New(deck.WithJokers(2))
	fmt.Printf("   Total cards: %d\n", len(cardsWithJokers))

	// Count and show jokers
	jokerCount := 0
	for _, card := range cardsWithJokers {
		if card.IsJoker() {
			jokerCount++
			fmt.Printf("   Found joker: %s\n", card)
		}
	}
	fmt.Printf("   Total jokers: %d\n\n", jokerCount)

	// Test 4: Filtered deck (no 2s and 3s)
	fmt.Println("4. Deck without 2s and 3s:")
	filteredCards := deck.New(deck.WithoutRanks(deck.Two, deck.Three))
	fmt.Printf("   Cards in filtered deck: %d\n", len(filteredCards))

	has2or3 := false
	for _, card := range filteredCards {
		if card.Rank == deck.Two || card.Rank == deck.Three {
			has2or3 = true
			break
		}
	}
	fmt.Printf("   Contains 2s or 3s: %t\n\n", has2or3)

	fmt.Println("5. Triple Deck (3 standard decks):")
	tripleCards := deck.New(deck.WithDecks(3))
	fmt.Printf("   Total cards: %d\n", len(tripleCards))
	fmt.Printf("   First Ace of Spades at position: 0\n")
	fmt.Printf("   Second Ace of Spades at position: 52\n")
	fmt.Printf("   Third Ace of Spades at position: 104\n")
	fmt.Printf("   All three are identical: %t\n\n", 
		tripleCards[0].String() == tripleCards[52].String() && 
		tripleCards[52].String() == tripleCards[104].String())

	// Test 6: Sorted deck
	fmt.Println("6. Default Sorted Deck:")
	sortedCards := deck.New(deck.DefaultSort)
	fmt.Printf("   First 4 cards (should be Spades A-4): ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%s, ", sortedCards[i])
	}
	fmt.Println("\n")

	// Test 7: Custom sorted deck (by rank first)
	fmt.Println("7. Custom Sorted Deck (by rank, then suit):")
	customSorted := deck.New(
		deck.Sort(func(cards []deck.Card) func(i, j int) bool {
			return func(i, j int) bool {
				if cards[i].Rank != cards[j].Rank {
					return cards[i].Rank < cards[j].Rank
				}
				return cards[i].Suit < cards[j].Suit
			}
		}),
	)
	fmt.Printf("   First 4 cards (should be all Aces): ")
	for i := 0; i < 4; i++ {
		fmt.Printf("%s, ", customSorted[i])
	}
	fmt.Println("\n")
}
