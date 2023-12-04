// websocket/chat_room.go
package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/gorilla/websocket"
)

// Available colors starts as a copy of the colors slice
var availableColors = append([]string(nil), colors...)

// ChatRoom represents a chat room with connected clients
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:   make(map[*websocket.Conn]bool),
		usernames: make(map[*websocket.Conn]string),
		messages:  make(chan Message, 256), // Buffered channel for messages
	}
}

// AddClient adds a new client to the chat room
func (cr *ChatRoom) AddClient(client *websocket.Conn, username string) {

	// If there are no more available colors, reset the list
	if len(availableColors) == 0 {
		availableColors = append([]string(nil), colors...)
	}

	// Choose a random color from the available colors
	index := rand.Intn(len(availableColors))
	color := availableColors[index]

	// Remove the chosen color from the available colors
	availableColors = append(availableColors[:index], availableColors[index+1:]...)

	cr.usernames[client] = color + username + "\033[0m" // Reset color after username

	cr.mutex.Lock()
	cr.clients[client] = true
	cr.mutex.Unlock()

	// Broadcast that the new user has joined, along with the updated list of connected users
	message := username + " has joined. Connected Users: " + cr.GetConnectedUserList()
	cr.Broadcast(client, []byte(message))
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

// Use the color when broadcasting messages
func (cr *ChatRoom) Broadcast(sender *websocket.Conn, message []byte) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	for client := range cr.clients {
		if client != sender {
			formattedMessage := fmt.Sprintf("%s", string(message))
			err := client.WriteMessage(websocket.TextMessage, []byte(formattedMessage))
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		}
	}
}

// Start broadcasting messages in a separate goroutine
func (cr *ChatRoom) StartBroadcasting() {
	log.Printf("StartBroadcasting, chatroom: %v", cr)
	for msg := range cr.messages {
		log.Printf("Received message from channel: %s\n", string(msg.Content))

		cr.mutex.Lock()
		log.Printf("Number of connected clients: %d\n", len(cr.clients)) // Add this line
		for client := range cr.clients {
			if client != nil && client != msg.Sender {
				err := client.WriteMessage(websocket.TextMessage, msg.Content)
				if err != nil {
					log.Printf("Error broadcasting message: %v", err)
				}
			}
		}
		cr.mutex.Unlock()
	}
}
