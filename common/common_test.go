package common

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
			cards:              []Card{{Suit: Clubs, RankIndex: 2}, {Suit: Clubs, RankIndex: 3}, {Suit: Clubs, RankIndex: 4}, {Suit: Clubs, RankIndex: 5}},
			expectedWinningIdx: 3,
		},
		{
			cards:              []Card{{Suit: Clubs, RankIndex: 2}, {Suit: Clubs, RankIndex: 3}, {Suit: Clubs, RankIndex: 4}, {Suit: Hearts, RankIndex: 5}},
			expectedWinningIdx: 2,
		},
		{
			cards:              []Card{{Suit: Hearts, RankIndex: 2}, {Suit: Clubs, RankIndex: 3}, {Suit: Clubs, RankIndex: 4}, {Suit: Clubs, RankIndex: 5}},
			expectedWinningIdx: 0,
		},
		{
			cards:              []Card{{Suit: Clubs, RankIndex: 2}, {Suit: Spades, RankIndex: 3}, {Suit: Clubs, RankIndex: 4}, {Suit: Clubs, RankIndex: 5}},
			expectedWinningIdx: 1,
		},
		{
			cards:              []Card{{Suit: Clubs, RankIndex: 2}, {Suit: Spades, RankIndex: 3}, {Suit: Spades, RankIndex: 4}, {Suit: Clubs, RankIndex: 5}},
			expectedWinningIdx: 2,
		},
	}

	for idx, tt := range tests {
		t.Run(strconv.Itoa(idx), func(t *testing.T) {
			assert.Equal(t, tt.expectedWinningIdx, winner(tt.cards))
		})
	}
}
