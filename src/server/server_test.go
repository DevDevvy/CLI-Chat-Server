// server/server_test.go
package server

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestServerFunction(t *testing.T) {
	go StartServer() // Run the server function in a separate goroutine

	// Give the server a second to start
	time.Sleep(3 * time.Second)

	// Send a WebSocket upgrade request to the server, The Dial function returns a WebSocket connection
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Fatalf("Could not send WebSocket upgrade request: %v", err)
	}
	defer c.Close()
}
