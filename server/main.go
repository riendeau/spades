package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type player struct {
	token        string
	name         string
	eventChannel chan string
}

var players []player

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func join(conn *websocket.Conn, name string) {
	log.Printf("Received register request; name: %s", name)
	if len(players) >= 4 {
		sendMessage(conn, "joinreply full")
		return
	}

	newPlayer := player{
		name:         name,
		token:        uuid.NewString(),
		eventChannel: make(chan string),
	}

	go func() {
		for nextEvent := range newPlayer.eventChannel {
			sendMessage(conn, nextEvent)
		}
	}()

	players = append(players, newPlayer)

	joinReply := "players"
	for _, player := range players {
		joinReply = fmt.Sprintf("%s %s", joinReply, player.name)
	}

	for i := range players {
		players[i].eventChannel <- joinReply
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "spades.html")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			processMessage(conn, string(msg))
		}
	}()
}

func sendMessage(conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Printf("Error sending message %s: %v", message, err)
	} else {
		log.Printf("Sent message %s", message)
	}
}

func processMessage(conn *websocket.Conn, message string) {
	log.Printf("Received message: %s", message)

	segments := strings.Split(message, " ")
	switch segments[0] {
	case "join":
		join(conn, segments[1])
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
