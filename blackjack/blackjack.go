package main

import (
	"fmt"
	"strings"
	"blackjack/deck"
)

type Player struct {
	Name  string
	Hand  []deck.Card
	Score int
}

type Game struct {
	Deck   []deck.Card
	Player *Player
	Dealer *Player
}

func NewGame() *Game {
	cards := deck.New()
	cards = deck.Shuffle(cards)
	
	return &Game{
		Deck:   cards,
		Player: &Player{Name: "Player", Hand: []deck.Card{}},
		Dealer: &Player{Name: "Dealer", Hand: []deck.Card{}},
	}
}

func (g *Game) Deal() deck.Card {
	card := g.Deck[0]
	g.Deck = g.Deck[1:]
	return card
}

func calculateScore(hand []deck.Card) int {
	score := 0
	aces := 0
	
	for _, card := range hand {
		switch card.Rank {
		case deck.Ace:
			aces++
		case deck.Jack, deck.Queen, deck.King:
			score += 10
		default:
			score += int(card.Rank)
		}
	}
	
	for aces > 0 {
		if score+11 <= 21 {
			score += 11
		} else {
			score += 1
		}
		aces--
	}
	
	return score
}

func isSoft17(hand []deck.Card) bool {
	if calculateScore(hand) != 17 {
		return false
	}
	
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
	
	return hasAce && score == 6
}

func (g *Game) dealInitialCards() {
	g.Player.Hand = append(g.Player.Hand, g.Deal())
	g.Dealer.Hand = append(g.Dealer.Hand, g.Deal())
	g.Player.Hand = append(g.Player.Hand, g.Deal())
	g.Dealer.Hand = append(g.Dealer.Hand, g.Deal())
	
	g.Player.Score = calculateScore(g.Player.Hand)
	g.Dealer.Score = calculateScore(g.Dealer.Hand)
}

func (g *Game) displayHands(showDealerCard bool) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	
	fmt.Printf("%s's Hand (Score: %d):\n", g.Player.Name, g.Player.Score)
	for i, card := range g.Player.Hand {
		fmt.Printf("  %d. %s\n", i+1, card)
	}
	
	fmt.Println()
	
	fmt.Printf("%s's Hand:\n", g.Dealer.Name)
	if showDealerCard {
		fmt.Printf("  Score: %d\n", g.Dealer.Score)
		for i, card := range g.Dealer.Hand {
			fmt.Printf("  %d. %s\n", i+1, card)
		}
	} else {
		fmt.Printf("  1. %s\n", g.Dealer.Hand[0])
		fmt.Printf("  2. [Hidden Card]\n")
	}
	
	fmt.Println(strings.Repeat("=", 50))
}

func (g *Game) playerTurn() bool {
	for {
		g.displayHands(false)
		
		if g.Player.Score > 21 {
			fmt.Printf("\n%s busts with %d! Dealer wins!\n", g.Player.Name, g.Player.Score)
			return false
		}
		
		if g.Player.Score == 21 {
			fmt.Printf("\n%s has 21!\n", g.Player.Name)
			return true
		}
		
		fmt.Print("\nChoose an action:\n1. Hit\n2. Stand\nEnter choice (1 or 2): ")
		var choice string
		fmt.Scanln(&choice)
		
		switch choice {
		case "1":
			card := g.Deal()
			g.Player.Hand = append(g.Player.Hand, card)
			g.Player.Score = calculateScore(g.Player.Hand)
			fmt.Printf("\n%s drew: %s\n", g.Player.Name, card)
		case "2":
			fmt.Printf("\n%s stands with %d\n", g.Player.Name, g.Player.Score)
			return true
		default:
			fmt.Println("Invalid choice. Please enter 1 or 2.")
		}
	}
}

func (g *Game) dealerTurn() {
	fmt.Println("\nDealer's turn:")
	g.displayHands(true)
	
	for g.Dealer.Score < 17 || isSoft17(g.Dealer.Hand) {
		fmt.Printf("\nDealer hits (Score: %d)\n", g.Dealer.Score)
		card := g.Deal()
		g.Dealer.Hand = append(g.Dealer.Hand, card)
		g.Dealer.Score = calculateScore(g.Dealer.Hand)
		fmt.Printf("Dealer drew: %s\n", card)
		
		g.displayHands(true)
		
		if g.Dealer.Score > 21 {
			fmt.Printf("\nDealer busts with %d!\n", g.Dealer.Score)
			return
		}
	}
	
	fmt.Printf("\nDealer stands with %d\n", g.Dealer.Score)
}

func (g *Game) determineWinner() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("FINAL RESULTS:")
	fmt.Printf("Player Score: %d\n", g.Player.Score)
	fmt.Printf("Dealer Score: %d\n", g.Dealer.Score)
	
	if g.Player.Score > 21 {
		fmt.Println("Player busted! Dealer wins!")
		return
	}
	
	if g.Dealer.Score > 21 {
		fmt.Println("Dealer busted! Player wins!")
		return
	}
	
	if g.Player.Score > g.Dealer.Score {
		fmt.Println("Player wins!")
	} else if g.Dealer.Score > g.Player.Score {
		fmt.Println("Dealer wins!")
	} else {
		fmt.Println("It's a tie!")
	}
	fmt.Println(strings.Repeat("=", 50))
}

func (g *Game) Play() {
	fmt.Println("Welcome to Blackjack!")
	
	g.dealInitialCards()
	
	playerContinues := g.playerTurn()
	
	if playerContinues && g.Player.Score <= 21 {
		g.dealerTurn()
		g.determineWinner()
	}
}