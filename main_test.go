package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinner(t *testing.T) {
	tests := []struct {
		cards              []Card
		expectedWinningIdx int
	}{
		{
			cards:              []Card{{suit: Clubs, number: 2}, {suit: Clubs, number: 3}, {suit: Clubs, number: 4}, {suit: Clubs, number: 5}},
			expectedWinningIdx: 3,
		},
		{
			cards:              []Card{{suit: Clubs, number: 2}, {suit: Clubs, number: 3}, {suit: Clubs, number: 4}, {suit: Hearts, number: 5}},
			expectedWinningIdx: 2,
		},
		{
			cards:              []Card{{suit: Hearts, number: 2}, {suit: Clubs, number: 3}, {suit: Clubs, number: 4}, {suit: Clubs, number: 5}},
			expectedWinningIdx: 0,
		},
		{
			cards:              []Card{{suit: Clubs, number: 2}, {suit: Spades, number: 3}, {suit: Clubs, number: 4}, {suit: Clubs, number: 5}},
			expectedWinningIdx: 1,
		},
		{
			cards:              []Card{{suit: Clubs, number: 2}, {suit: Spades, number: 3}, {suit: Spades, number: 4}, {suit: Clubs, number: 5}},
			expectedWinningIdx: 2,
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			assert.Equal(t, tt.expectedWinningIdx, winner(tt.cards))
		})
	}
}
