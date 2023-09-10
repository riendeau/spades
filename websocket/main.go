package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	http.ServeFile(w, r, "websocket.html")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error reading message: %v", err)
		return
	}
	log.Printf("Received message: %s", string(msg))

	var reply = ""
	for _, c := range string(msg) {
		reply = string(c) + reply
	}

	if err = conn.WriteMessage(websocket.TextMessage, []byte(reply)); err != nil {
		log.Printf("Error writing message: %v", err)
		return
	}
	log.Printf("Sent message: %s", reply)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/echo", serveWs)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
