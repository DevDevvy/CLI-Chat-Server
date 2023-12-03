// websocket/chat_room.go
package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ChatRoom represents a chat room with connected clients
type ChatRoom struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

// NewChatRoom creates a new ChatRoom instance
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient adds a new client to the chat room
func (cr *ChatRoom) AddClient(client *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.clients[client] = true
}

// RemoveClient removes a client from the chat room
func (cr *ChatRoom) RemoveClient(client *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	delete(cr.clients, client)
}

// Broadcast sends a message to all clients in the chat room, excluding the sender
func (cr *ChatRoom) Broadcast(sender *websocket.Conn, message []byte) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	for client := range cr.clients {
		if client != sender {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Handle error (e.g., client disconnected)
			}
		}
	}
}
