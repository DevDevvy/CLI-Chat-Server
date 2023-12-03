// websocket/manager.go
package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chatRoom = NewChatRoom()

func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("Client Connected")
	chatRoom.AddClient(conn)
	defer func() {
		log.Println("Client Disconnected")
		chatRoom.RemoveClient(conn)
	}()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message %v from client: %v", messageType, err)
			return
		}

		log.Printf("Received message: %s", p)

		// Broadcast the message to all clients in the chat room, excluding the sender
		chatRoom.Broadcast(conn, p)
	}
}
