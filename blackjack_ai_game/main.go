package main

import (
	"blackjack_ai_game/blackjack"
	"blackjack_ai_game/deck"
	"fmt"
)

type BasicAI struct {
	winnings int
	handsPlayed int
}

func (ai *BasicAI) Bet() int {
	return 1
}

func (ai *BasicAI) Play(hand []deck.Card, dealer deck.Card) string {
	playerScore := blackjack.Score(hand)
	dealerUpCard := int(dealer.Rank)
	if dealer.Rank == deck.Ace {
		dealerUpCard = 11
	} else if dealer.Rank >= deck.Jack {
		dealerUpCard = 10
	}
	
	if blackjack.CanSplit(hand) {
		return ai.shouldSplit(hand, dealerUpCard)
	}
	
	if blackjack.CanDouble(hand) && ai.shouldDouble(hand, dealerUpCard) {
		return "double"
	}
	
	return ai.hitOrStand(playerScore, dealerUpCard, ai.isSoftHand(hand))
}

func (ai *BasicAI) Results(hand []deck.Card, dealer []deck.Card, amount int) {
	ai.winnings += amount
	ai.handsPlayed++
}

func (ai *BasicAI) GetStats() (int, int) {
	return ai.winnings, ai.handsPlayed
}

func (ai *BasicAI) shouldSplit(hand []deck.Card, dealerUpCard int) string {
	cardValue := int(hand[0].Rank)
	if hand[0].Rank == deck.Ace {
		cardValue = 11
	} else if hand[0].Rank >= deck.Jack {
		cardValue = 10
	}
	
	switch cardValue {
	case 11:
		return "split"
	case 8:
		return "split"
	case 9:
		if dealerUpCard == 7 || dealerUpCard == 10 || dealerUpCard == 11 {
			return "stand"
		}
		return "split"
	case 7:
		if dealerUpCard <= 7 {
			return "split"
		}
		return "hit"
	case 6:
		if dealerUpCard <= 6 {
			return "split"
		}
		return "hit"
	case 4:
		if dealerUpCard == 5 || dealerUpCard == 6 {
			return "hit"
		}
		return "hit"
	case 3, 2:
		if dealerUpCard <= 7 {
			return "split"
		}
		return "hit"
	default:
		return "hit"
	}
}

func (ai *BasicAI) shouldDouble(hand []deck.Card, dealerUpCard int) bool {
	playerScore := blackjack.Score(hand)
	
	if ai.isSoftHand(hand) {
		switch playerScore {
		case 13, 14:
			return dealerUpCard == 5 || dealerUpCard == 6
		case 15, 16:
			return dealerUpCard >= 4 && dealerUpCard <= 6
		case 17, 18:
			return dealerUpCard >= 3 && dealerUpCard <= 6
		default:
			return false
		}
	} else {
		switch playerScore {
		case 9:
			return dealerUpCard >= 3 && dealerUpCard <= 6
		case 10:
			return dealerUpCard <= 9
		case 11:
			return dealerUpCard <= 10
		default:
			return false
		}
	}
}

func (ai *BasicAI) hitOrStand(playerScore int, dealerUpCard int, isSoft bool) string {
	if isSoft {
		switch playerScore {
		case 18:
			if dealerUpCard == 9 || dealerUpCard == 10 || dealerUpCard == 11 {
				return "hit"
			}
			return "stand"
		case 19, 20, 21:
			return "stand"
		default:
			return "hit"
		}
	} else {
		if playerScore >= 17 {
			return "stand"
		}
		if playerScore <= 11 {
			return "hit"
		}
		
		switch playerScore {
		case 12:
			if dealerUpCard >= 4 && dealerUpCard <= 6 {
				return "stand"
			}
			return "hit"
		case 13, 14, 15, 16:
			if dealerUpCard <= 6 {
				return "stand"
			}
			return "hit"
		default:
			return "hit"
		}
	}
}

func (ai *BasicAI) isSoftHand(hand []deck.Card) bool {
	score := 0
	hasAce := false
	
	for _, card := range hand {
		switch card.Rank {
		case deck.Ace:
			hasAce = true
		case deck.Jack, deck.Queen, deck.King:
			score += 10
		default:
			score += int(card.Rank)
		}
	}
	
	return hasAce && score <= 10
}

func main() {
	fmt.Println("Starting Blackjack AI Test...")
	
	var ai BasicAI
	
	testConfigs := []blackjack.Options{
		{Hands: 100, Decks: 1},
		{Hands: 1000, Decks: 2},
		{Hands: 10000, Decks: 6},
	}
	
	for _, opts := range testConfigs {
		fmt.Printf("\nTesting with %d hands, %d deck(s):\n", opts.Hands, opts.Decks)
		
		ai = BasicAI{}
		
		game := blackjack.New(opts)
		winnings := game.Play(&ai)
		
		totalWinnings, handsPlayed := ai.GetStats()
		winRate := float64(totalWinnings) / float64(handsPlayed) * 100
		
		fmt.Printf("Total winnings: %d units\n", totalWinnings)
		fmt.Printf("Hands played: %d\n", handsPlayed)
		fmt.Printf("Win rate: %.2f%% per hand\n", winRate)
		
		if winnings != totalWinnings {
			fmt.Printf("Warning: Game.Play() returned %d but AI tracked %d\n", winnings, totalWinnings)
		}
	}
	
	fmt.Println("\nAI testing complete!")
}

type SimpleAI struct{}

func (ai *SimpleAI) Bet() int {
	return 1
}

func (ai *SimpleAI) Play(hand []deck.Card, dealer deck.Card) string {
	score := blackjack.Score(hand)
	if score < 17 {
		return "hit"
	}
	return "stand"
}

func (ai *SimpleAI) Results(hand []deck.Card, dealer []deck.Card, amount int) {
}