package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

type Suit uint8

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

var suits []string = []string{"♣", "♦", "♥", "♠"}
var numbers []string = []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

type Card struct {
	suit   Suit
	number int
}

func (card Card) String() string {
	return fmt.Sprintf("%s%s", suits[card.suit], numbers[card.number])
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
	return hand[i].suit < hand[j].suit || (hand[i].suit == hand[j].suit && hand[i].number < hand[j].number)
}
func (hand Hand) Swap(i, j int) { hand[i], hand[j] = hand[j], hand[i] }

func main() {
	rand.Seed(time.Now().UnixNano())

	hands := randomHands()
	for handIdx, hand := range hands {
		fmt.Printf("Hand %d: %s\n", handIdx+1, hand)
	}
}

func randomHands() [4]Hand {
	var deck []Card
	for suit := Clubs; suit <= Spades; suit++ {
		for number := 0; number <= 12; number++ {
			deck = append(deck, Card{suit: suit, number: number})
		}
	}

	var hands [4]Hand
	for len(deck) > 0 {
		pickedIdx := rand.Intn(len(deck))
		pickedCard := deck[pickedIdx]
		deck = slices.Delete(deck, pickedIdx, pickedIdx+1)

		handIdx := (52 - len(deck) - 1) / 13
		hands[handIdx] = append(hands[handIdx], pickedCard)
	}

	for _, hand := range hands {
		sort.Sort(hand)
	}

	return hands
}

func winner(cards [4]Card) int {
	winningIdx := 0
	for cardIdx := 1; cardIdx < 4; cardIdx++ {
		winnerSoFar := cards[winningIdx]
		if cards[cardIdx].suit == winnerSoFar.suit {
			if cards[cardIdx].number > winnerSoFar.number {
				winningIdx = cardIdx
			}
		} else if cards[cardIdx].suit == Spades {
			winningIdx = cardIdx
		}
	}
	return winningIdx
}
