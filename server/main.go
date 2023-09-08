package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"

	"github.com/riendeau/spades/common"
)

func main() {
	http.HandleFunc("/hands", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(randomHands(4))
	})

	http.ListenAndServe(":8765", nil)
}

func shuffledDeck() []common.Card {
	var deck []common.Card
	for suit := range common.SuitToRune {
		for number := range common.CardRanks {
			deck = append(deck, common.Card{Suit: suit, RankIndex: number})
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
