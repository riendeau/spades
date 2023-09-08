package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/riendeau/spades/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := common.NewSpadesClient(conn)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		registerReply, err := c.Register(ctx, &common.RegisterRequest{Name: name})
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(registerReply.GetToken()))
			return
		}

		log.Printf("error when registering %s: %v", name, err)
		if err == common.ErrGameFull {
			w.WriteHeader(http.StatusTeapot)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
	})

	err = http.ListenAndServe(":8765", nil)
	log.Printf("Server terminted with %v", err)
}
