// websocket/manager.go
package websocket

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	// Authenticate the user
	err = authenticateUser(conn)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		return
	}

	// Prompt for a temporary username in a separate goroutine
	go promptForUsername(conn)

	// Announce the new user connection to all clients
	username := chatRoom.GetUsername(conn)
	announcement := fmt.Sprintf("%s has joined. Connected Users: %s", username, chatRoom.GetConnectedUserList())
	chatRoom.Broadcast(conn, []byte(announcement))

	// Send the initial list of connected users to the newly connected user
	chatRoom.SendUserList(conn)

	// Handle user input in a separate goroutine
	go handleUserInput(conn, username)

	// Broadcast the updated list of connected users when a user disconnects
	defer func() {
		log.Println("Client Disconnected")
		chatRoom.RemoveClient(conn)
		chatRoom.BroadcastUserList()
	}()

	// Keep the main goroutine alive until the connection is closed
	select {}
}

// handleUserInput reads messages from the user and broadcasts them to others
func handleUserInput(conn *websocket.Conn, username string) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading %v message from client: %v", messageType, err)
			return
		}

		log.Printf("Received message: %s", p)

		// Broadcast the message with the username prefix to all clients
		chatRoom.Broadcast(conn, []byte(username+": "+string(p)))
	}
}

// authenticateUser checks the user's password
func authenticateUser(conn *websocket.Conn) error {
	const correctPassword = "password"

	err := conn.WriteMessage(websocket.TextMessage, []byte("Enter the password: "))
	if err != nil {
		return err
	}

	_, passwordBytes, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// Trim whitespaces from the received password
	receivedPassword := strings.TrimSpace(string(passwordBytes))

	// Check if the received message is not empty
	if len(receivedPassword) == 0 {
		return errors.New("password not provided")
	}

	fmt.Printf("Received password: %s\n", receivedPassword) // Print the received password

	if receivedPassword != correctPassword {
		return errors.New("incorrect password")
	}

	return nil
}

// promptForUsername prompts the user to enter a temporary username
func promptForUsername(conn *websocket.Conn) {
	err := conn.WriteMessage(websocket.TextMessage, []byte("Enter a temporary username: "))
	if err != nil {
		log.Printf("Error writing message to client: %v", err)
		return
	}

	_, username, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error reading message from client: %v", err)
		return
	}

	chatRoom.SetUsername(conn, string(username))
}
