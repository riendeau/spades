package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Args[1]), nil)
	log.Printf("Server terminted with %v", err)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	registerReply, err := c.Register(ctx, &common.RegisterRequest{Name: name})
	if err == nil {
		receiveEvents(registerReply.Token)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(registerReply.Token))
		return
	}

	log.Printf("error when registering %s: %v", name, err)
	statusCode := http.StatusInternalServerError
	if err == common.ErrGameFull {
		statusCode = http.StatusTeapot
	}
	w.WriteHeader(statusCode)
}

func receiveEvents(token string) {
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
				continue
			}

			log.Printf("Received event: %s", event.Description)
		}
	}()
}
