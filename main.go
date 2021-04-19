package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ChatChannel struct {
	Id   int
	Name string
	//Clients   []*websocket.Conn
	Clients   map[*websocket.Conn]bool
	Broadcast chan string
}

var upgrader = websocket.Upgrader{
	// todo remove only for testing
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// create a sample channel
var channel1 = ChatChannel{
	Id:        0,
	Name:      "My Room",
	Clients:   make(map[*websocket.Conn]bool),  //map instead of slice because easier
	Broadcast: make(chan string),
}

func main() {
	http.HandleFunc("/gateway", handleConnections)

	go handleMessages()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func handleConnections(w http.ResponseWriter, r *http.Request) {

	// todo implement separate channels

	//channelId, present := r.URL.Query()["channel"]

	// if no channel ID specified
	//if !present {
	//	w.WriteHeader(404)
	//	return
	//}

	// upgrade to WS
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//todo testing only
	//print(channelId)

	// we close the http connection when the function returns
	defer conn.Close()

	//clients = append(clients, conn)
	channel1.Clients[conn] = true

	for {
		var _, msg, err = conn.ReadMessage()

		if err != nil {
			delete(channel1.Clients, conn)
			break
		}

		// Send the newly received message to the broadcast channel
		channel1.Broadcast <- string(msg)
	}
}

func handleMessages() {
	for {
		// when message arrives in chan
		msg := <-channel1.Broadcast
		println("howdy, received message")
		// send to clients
		for client := range channel1.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				delete(channel1.Clients, client)
				client.Close()
			}
		}
	}
}
