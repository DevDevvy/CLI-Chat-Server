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

// AuthenticateUserResponse represents the result of an authentication attempt
type AuthenticateUserResponse struct {
	Success  bool
	Username string
}

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
	authResponse, err := authenticateUser(conn)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		return
	}

	// If authentication is unsuccessful, close the connection
	if !authResponse.Success {
		log.Println("Authentication failed. Closing connection.")
		return
	}

	// Set the username for the connection and add the client to the chat room
	chatRoom.SetUsername(conn, authResponse.Username)
	chatRoom.AddClient(conn, authResponse.Username)

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

	// Start broadcasting if it's not already running
	chatRoom.StartBroadcastingOnce()

	// Keep the main goroutine alive until the connection is closed
	select {}
}

// promptForUsername prompts the user for a temporary username
func promptForUsername(conn *websocket.Conn, authSuccess bool) {
	if !authSuccess {
		log.Println("Authentication failed. Not prompting for username.")
		return
	}

	log.Println("Prompting for username")

	// Prompt for a temporary username
	err := conn.WriteMessage(websocket.TextMessage, []byte("Enter a temporary username: "))
	if err != nil {
		log.Printf("Error prompting for username: %v", err)
		return
	}

	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Error reading username: %v", err)
		return
	}

	// Trim whitespaces from the received username
	username := strings.TrimSpace(string(usernameBytes))

	// Check if the received message is not empty
	if len(username) == 0 {
		log.Println("Username not provided")
		return
	}

	// Set the username for the connection
	chatRoom.SetUsername(conn, username)

	// Log that the client is added (you can remove this if it becomes noisy)
	chatRoom.AddClient(conn, username)
	log.Println("Username set and client added")
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
func authenticateUser(conn *websocket.Conn) (*AuthenticateUserResponse, error) {
	const correctPassword = "password"

	err := conn.WriteMessage(websocket.TextMessage, []byte("Enter the password: "))
	if err != nil {
		return nil, err
	}

	_, passwordBytes, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Trim whitespaces from the received password
	receivedPassword := strings.TrimSpace(string(passwordBytes))

	// Check if the received message is not empty
	if len(receivedPassword) == 0 {
		return nil, errors.New("password not provided")
	}

	fmt.Printf("Received password: %s\n", receivedPassword) // Print the received password

	if receivedPassword != correctPassword {
		return nil, errors.New("incorrect password")
	}

	// Authentication successful, prompt for username
	err = conn.WriteMessage(websocket.TextMessage, []byte("Enter a temporary username: "))
	if err != nil {
		return nil, err
	}

	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	// Trim whitespaces from the received username
	username := strings.TrimSpace(string(usernameBytes))

	// Check if the received message is not empty
	if len(username) == 0 {
		return nil, errors.New("username not provided")
	}

	fmt.Printf("Received username: %s\n", username) // Print the received username

	return &AuthenticateUserResponse{
		Success:  true,
		Username: username,
	}, nil
}

// Start broadcasting messages in a separate goroutine
func (cr *ChatRoom) StartBroadcastingOnce() {
	cr.startOnce.Do(func() {
		go cr.StartBroadcasting()
	})
}
