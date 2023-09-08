package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sort"

	"github.com/google/uuid"
	"github.com/riendeau/spades/common"
	"google.golang.org/grpc"
)

type server struct {
	common.UnimplementedSpadesServer
}

var uuidToName = make(map[string]string, 4)

func (s *server) Register(ctx context.Context, in *common.RegisterRequest) (*common.RegisterReply, error) {
	log.Printf("Received register request; name: %s", in.GetName())
	if len(uuidToName) >= 4 {
		return nil, common.ErrGameFull
	}

	uuid := uuid.NewString()
	uuidToName[uuid] = in.GetName()

	return &common.RegisterReply{Token: uuid}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	common.RegisterSpadesServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

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
