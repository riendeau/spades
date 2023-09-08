package common

import (
	"errors"
	"fmt"
	"strings"
)

type Suit uint8

const (
	Spades Suit = iota
	Hearts
	Clubs
	Diamonds
)

var SuitToRune = map[Suit]string{Spades: "♠", Hearts: "♥", Clubs: "♣", Diamonds: "♦"}
var CardRanks = [...]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

var ErrGameFull = errors.New("game is already full")

type Card struct {
	Suit      Suit `json:"suit"`
	RankIndex int  `json:"rankIndex"`
}

func (card Card) String() string {
	return fmt.Sprintf("%s%s", SuitToRune[card.Suit], CardRanks[card.RankIndex])
}

type Hand []Card

func (hand Hand) String() string {
	var sb strings.Builder
	for i, card := range hand {
		if i > 0 {
			fmt.Fprint(&sb, " ")
		}
		fmt.Fprintf(&sb, "%s", card)
	}
	return sb.String()
}

func (hand Hand) Len() int { return len(hand) }
func (hand Hand) Less(i, j int) bool {
	return hand[i].Suit < hand[j].Suit || (hand[i].Suit == hand[j].Suit && hand[i].RankIndex < hand[j].RankIndex)
}
func (hand Hand) Swap(i, j int) { hand[i], hand[j] = hand[j], hand[i] }

func winner(cards []Card) int {
	winningIdx := 0
	for cardIdx := 1; cardIdx < len(cards); cardIdx++ {
		winnerSoFar := cards[winningIdx]
		if cards[cardIdx].Suit == winnerSoFar.Suit {
			if cards[cardIdx].RankIndex > winnerSoFar.RankIndex {
				winningIdx = cardIdx
			}
		} else if cards[cardIdx].Suit == Spades {
			winningIdx = cardIdx
		}
	}
	return winningIdx
}
