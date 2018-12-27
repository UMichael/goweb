package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

type Message struct {
	User    User
	Message string `json:"message"`
}

func HandleConn(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//handle
	}
	defer ws.Close()
	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			//Assume client is out
			//i'll try and create a log
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}
func SendMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				//sort error
				client.Close()
				delete(clients, client)
			}
		}
		//would send messages to database for users who have created and are not online
		Db.QueryRow("insert into messages(nickname, content) values($1, $2)", msg.User.Nickname, msg.Message)
	}
}
