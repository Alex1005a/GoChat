package main

import (
	"awesomeProject1/src/app/application"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Text string `json:"text"`
}

var connection *websocket.Conn


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ws", wsHandler)
	router.HandleFunc("/", createMessage).Methods("POST")
	panic(http.ListenAndServe(":8080", router))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	connection = conn
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var mes Message
	_ = json.NewDecoder(r.Body).Decode(&mes)

	application.NewMessageService().CreateMessage(mes.Text, connection)
	w.WriteHeader(200)
}