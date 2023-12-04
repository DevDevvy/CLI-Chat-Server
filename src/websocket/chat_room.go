// websocket/chat_room.go
package websocket

import (
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// ChatRoom represents a chat room with connected clients
type ChatRoom struct {
	clients   map[*websocket.Conn]bool
	usernames map[*websocket.Conn]string
	mutex     sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:   make(map[*websocket.Conn]bool),
		usernames: make(map[*websocket.Conn]string),
	}
}

// AddClient adds a new client to the chat room
func (cr *ChatRoom) AddClient(client *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.clients[client] = true

	// Retrieve the username associated with the client
	username, exists := cr.usernames[client]
	if !exists {
		log.Printf("Error: No username found for client %v", client)
	}

	log.Printf("Client %v (%s) added. Total clients: %d", client, username, len(cr.clients))

}

// RemoveClient removes a client from the chat room
func (cr *ChatRoom) RemoveClient(client *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	delete(cr.clients, client)
}

// SetUsername sets the temporary username for a connection
func (cr *ChatRoom) SetUsername(conn *websocket.Conn, username string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.usernames[conn] = username
}

// GetUserame gets the temporary username for a connection
func (cr *ChatRoom) GetUsername(conn *websocket.Conn) string {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	return cr.usernames[conn]
}

// SendUserList sends the list of connected users to a specific user
func (cr *ChatRoom) SendUserList(conn *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	var userList []string
	for _, username := range cr.usernames {
		userList = append(userList, username)
	}

	message := "Connected Users: " + strings.Join(userList, ", ")
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error sending user list to client: %v", err)
	}
}

// BroadcastUserList sends the updated list of connected users to all clients
func (cr *ChatRoom) BroadcastUserList() {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	var userList []string
	for _, username := range cr.usernames {
		userList = append(userList, username)
	}

	message := "Connected Users: " + strings.Join(userList, ", ")
	for client := range cr.clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			// Handle error (e.g., client disconnected)
		}
	}
}

// GetConnectedUserList returns a formatted list of connected users
func (cr *ChatRoom) GetConnectedUserList() string {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	var userList []string
	for _, username := range cr.usernames {
		userList = append(userList, username)
	}

	return strings.Join(userList, ", ")
}

// Broadcast sends a message to all connected clients except the sender
func (cr *ChatRoom) Broadcast(sender *websocket.Conn, message []byte) {
	log.Println("Broadcast")
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	for client := range cr.clients {
		log.Printf("client: %v", client)
		// if client != sender {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			// }
		}
	}
}
