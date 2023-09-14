package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/riendeau/spades/internal/server"
)

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

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", server.ServeWs)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
