package main

import (
	"context"
	"fmt"
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

type player struct {
	token        string
	name         string
	eventChannel chan string
}

var players []player

func (s *server) Register(ctx context.Context, in *common.RegisterRequest) (*common.RegisterReply, error) {
	log.Printf("Received register request; name: %s", in.GetName())
	if len(players) >= 4 {
		return nil, common.ErrGameFull
	}

	newPlayer := player{
		name:         in.GetName(),
		token:        uuid.NewString(),
		eventChannel: make(chan string),
	}

	for i := range players {
		players[i].eventChannel <- fmt.Sprintf("%s sat down", newPlayer.name)
	}

	players = append(players, newPlayer)
	return &common.RegisterReply{Token: newPlayer.token}, nil
}

func (s *server) Events(in *common.EventsRequest, srv common.Spades_EventsServer) error {
	for i := range players {
		if players[i].token == in.PlayerToken {
			for nextEvent := range players[i].eventChannel {
				srv.Send(&common.EventReply{Description: nextEvent})
			}
		}
	}
	return nil
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
