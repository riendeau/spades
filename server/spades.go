package main

import (
	"math/rand"
	"sort"

	"github.com/riendeau/spades/common"
)

func shuffledDeck() []common.Card {
	var deck []common.Card
	for suit := range common.SuitToRune {
		for rank := range common.CardRanks {
			deck = append(deck, common.Card{Suit: suit, RankIndex: rank})
		}
	}
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
	return deck
}

func randomHands(numHands int) []common.Hand {
	deck := shuffledDeck()
	hands := make([]common.Hand, numHands)
	cardsPerHand := len(deck) / numHands

	for i := 0; i < numHands; i++ {
		startIdx := i * cardsPerHand
		endIdx := startIdx + cardsPerHand
		hands[i] = deck[startIdx:endIdx]
		sort.Sort(hands[i])
	}

	return hands
}
