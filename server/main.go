package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type client struct {
	eventChannel chan string
	name         string
	seat         int
}

var clients []*client

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func getPlayers() [4]*client {
	var playerClients [4]*client
	for i := range clients {
		if clients[i].seat > 0 {
			playerClients[clients[i].seat-1] = clients[i]
		}
	}
	return [4]*client(playerClients)
}

func freeIdx() int {
	curPlayers := getPlayers()
	var freeIdx int
	for freeIdx = 0; freeIdx < 4; freeIdx++ {
		if curPlayers[freeIdx] == nil {
			break
		}
	}
	return freeIdx
}

func sit(client *client, name string) {
	log.Printf("Received sit request; name: %s", name)
	client.name = name

	if client.seat == 0 {
		freeIdx := freeIdx()
		if freeIdx > 3 {
			client.eventChannel <- "joinreply full"
			return
		}

		client.seat = freeIdx + 1
		log.Printf("Seated %s in seat %d", name, client.seat)
	}

	newPlayerList := buildPlayerList()
	for i := range clients {
		clients[i].eventChannel <- newPlayerList
	}

	if freeIdx := freeIdx(); freeIdx > 3 {
		for i := range clients {
			clients[i].eventChannel <- "newhand"
		}
	}
}

func buildPlayerList() string {
	playerListReply := "players"
	players := getPlayers()
	for i := range players {
		name := "(empty)"
		if players[i] != nil {
			name = players[i].name
		}
		playerListReply = fmt.Sprintf("%s %s", playerListReply, name)
	}
	return playerListReply
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "spades.html")
		return
	}
	if strings.HasPrefix(r.URL.Path, "/images/") {
		http.ServeFile(w, r, r.URL.Path[1:])
		return
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	sendChannel := make(chan string)
	client := &client{eventChannel: sendChannel}
	clients = append(clients, client)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			processMessage(client, string(msg))
		}
	}()

	go func() {
		for message := range sendChannel {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Printf("Error sending message %s: %v", message, err)
			} else {
				log.Printf("Sent message %s", message)
			}
		}
	}()
	sendChannel <- buildPlayerList()
}

func processMessage(client *client, message string) {
	log.Printf("Received message: %s", message)

	segments := strings.Split(message, " ")
	switch segments[0] {
	case "sit":
		sit(client, segments[1])
	}
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
