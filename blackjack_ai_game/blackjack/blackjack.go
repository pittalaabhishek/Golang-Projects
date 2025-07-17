package blackjack

import (
	"blackjack_ai_game/deck"
)

type AI interface {
	Bet() int
	
	Play(hand []deck.Card, dealer deck.Card) string
	
	Results(hand []deck.Card, dealer []deck.Card, amount int)
}

type Options struct {
	Hands int
	Decks int
}

type Game struct {
	deck     []deck.Card
	options  Options
	handsLeft int
}

type hand struct {
	cards    []deck.Card
	bet      int
	doubled  bool
	split    bool
}

type gameState struct {
	playerHands []hand
	dealerHand  []deck.Card
	currentHand int
}

func New(opts Options) *Game {
	cards := deck.New(deck.WithDecks(opts.Decks), deck.Shuffle)
	
	return &Game{
		deck:      cards,
		options:   opts,
		handsLeft: opts.Hands,
	}
}

func (g *Game) Play(ai AI) int {
	totalWinnings := 0
	
	for g.handsLeft > 0 {
		winnings := g.playHand(ai)
		totalWinnings += winnings
		g.handsLeft--
	}
	
	return totalWinnings
}

func (g *Game) playHand(ai AI) int {
	bet := ai.Bet()
	if bet <= 0 {
		bet = 1
	}
	
	playerHand := hand{
		cards: []deck.Card{g.deal(), g.deal()},
		bet:   bet,
	}
	
	dealerHand := []deck.Card{g.deal(), g.deal()}
	
	state := &gameState{
		playerHands: []hand{playerHand},
		dealerHand:  dealerHand,
		currentHand: 0,
	}
	
	if g.calculateScore(playerHand.cards) == 21 {
		if g.calculateScore(dealerHand) == 21 {
			ai.Results(playerHand.cards, dealerHand, 0)
			return 0
		} else {
			winnings := bet
			ai.Results(playerHand.cards, dealerHand, winnings)
			return winnings
		}
	}
	
	for state.currentHand < len(state.playerHands) {
		g.playPlayerHand(ai, state)
		state.currentHand++
	}
	
	g.playDealerHand(state)
	
	totalWinnings := 0
	for _, hand := range state.playerHands {
		winnings := g.calculateHandResult(hand, state.dealerHand)
		totalWinnings += winnings
		ai.Results(hand.cards, state.dealerHand, winnings)
	}
	
	return totalWinnings
}

func (g *Game) playPlayerHand(ai AI, state *gameState) {
	currentHand := &state.playerHands[state.currentHand]
	
	for {
		score := g.calculateScore(currentHand.cards)
		if score > 21 {
			break
		}
		
		if score == 21 {
			break
		}
		
		action := ai.Play(currentHand.cards, state.dealerHand[0])
		
		switch action {
		case "hit":
			currentHand.cards = append(currentHand.cards, g.deal())
			
		case "stand":
			return
			
		case "double":
			if len(currentHand.cards) == 2 && !currentHand.doubled {
				currentHand.doubled = true
				currentHand.bet *= 2
				currentHand.cards = append(currentHand.cards, g.deal())
					return
			}
			currentHand.cards = append(currentHand.cards, g.deal())
			
		case "split":
			if len(currentHand.cards) == 2 && 
			   currentHand.cards[0].Rank == currentHand.cards[1].Rank &&
			   !currentHand.split {
				g.splitHand(state)
				return
			}
			currentHand.cards = append(currentHand.cards, g.deal())
			
		default:
			currentHand.cards = append(currentHand.cards, g.deal())
		}
	}
}

func (g *Game) splitHand(state *gameState) {
	currentHand := &state.playerHands[state.currentHand]
	
	newHand := hand{
		cards: []deck.Card{currentHand.cards[1], g.deal()},
		bet:   currentHand.bet,
		split: true,
	}
	
	currentHand.cards = []deck.Card{currentHand.cards[0], g.deal()}
	currentHand.split = true
	
	state.playerHands = append(state.playerHands[:state.currentHand+1], 
		append([]hand{newHand}, state.playerHands[state.currentHand+1:]...)...)
}

func (g *Game) playDealerHand(state *gameState) {
	allBusted := true
	for _, hand := range state.playerHands {
		if g.calculateScore(hand.cards) <= 21 {
			allBusted = false
			break
		}
	}
	
	if allBusted {
		return
	}
	
	for g.calculateScore(state.dealerHand) < 17 || g.isSoft17(state.dealerHand) {
		state.dealerHand = append(state.dealerHand, g.deal())
	}
}

func (g *Game) calculateHandResult(playerHand hand, dealerHand []deck.Card) int {
	playerScore := g.calculateScore(playerHand.cards)
	dealerScore := g.calculateScore(dealerHand)
	
	if playerScore > 21 {
		return -playerHand.bet
	}
	
	if dealerScore > 21 {
		return playerHand.bet
	}
	
	if playerScore > dealerScore {
		return playerHand.bet
	} else if dealerScore > playerScore {
		return -playerHand.bet
	} else {
		return 0
	}
}

func (g *Game) deal() deck.Card {
	if len(g.deck) == 0 {
		g.deck = deck.New(deck.WithDecks(g.options.Decks), deck.Shuffle)
	}
	
	card := g.deck[0]
	g.deck = g.deck[1:]
	return card
}

func (g *Game) calculateScore(hand []deck.Card) int {
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

func (g *Game) isSoft17(hand []deck.Card) bool {
	if g.calculateScore(hand) != 17 {
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

func Score(hand []deck.Card) int {
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

func CanSplit(hand []deck.Card) bool {
	return len(hand) == 2 && hand[0].Rank == hand[1].Rank
}

func CanDouble(hand []deck.Card) bool {
	return len(hand) == 2
}