package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/riendeau/spades/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var c common.SpadesClient

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c = common.NewSpadesClient(conn)

	http.HandleFunc("/register", handleRegister)

	port := 8600 + rand.Intn(100)
	if len(os.Args) > 1 {
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalf("Could not parse specified port to integer: %v", err)
		}
	}

	log.Printf("Listening on port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Printf("Server terminated with %v", err)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	registerReply, err := c.Register(ctx, &common.RegisterRequest{Name: name})
	if err == nil {
		startEventStream(registerReply.Token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(registerReply.Token))
		return
	}

	log.Printf("error when registering %s: %v", name, err)
	statusCode := http.StatusInternalServerError
	if err == common.ErrGameFull {
		statusCode = http.StatusServiceUnavailable
	}
	w.WriteHeader(statusCode)
}

func startEventStream(token string) {
	eventsClient, err := c.Events(context.Background(), &common.EventsRequest{PlayerToken: token})
	if err != nil {
		log.Fatalf("could not establish events stream: %v", err)
	}
	defer eventsClient.CloseSend()
	go func() {
		for {
			event, err := eventsClient.Recv()
			if err != nil {
				log.Printf("Error receiving event: %v", err)
				return
			}
			handleEvent(event)
		}
	}()
}

func handleEvent(event *common.EventReply) {
	log.Printf("Received event: %s", event.Description)
}
