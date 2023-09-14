package spades

import (
	"math/rand"
	"sort"
)

func shuffledDeck() []Card {
	var deck []Card
	for suit := range SuitToRune {
		for rank := range CardRanks {
			deck = append(deck, Card{Suit: suit, RankIndex: rank})
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func RandomHands(numHands int) []Hand {
	deck := shuffledDeck()
	hands := make([]Hand, numHands)
	cardsPerHand := len(deck) / numHands

	for i := 0; i < numHands; i++ {
		startIdx := i * cardsPerHand
		endIdx := startIdx + cardsPerHand
		hands[i] = deck[startIdx:endIdx]
		sort.Sort(hands[i])
	}

	return hands
}
