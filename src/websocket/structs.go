package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Sender  *websocket.Conn
	Content []byte
}

// ChatRoom represents a chat room with connected clients
type ChatRoom struct {
	clients   map[*websocket.Conn]bool
	usernames map[*websocket.Conn]string
	messages  chan Message
	mutex     sync.Mutex
	startOnce sync.Once // Ensures StartBroadcasting is started only once
}

// ANSI color codes for usernames
var colors = []string{
	"\033[31m", // Red
	"\033[32m", // Green
	"\033[33m", // Yellow
	"\033[34m", // Blue
	"\033[35m", // Magenta
	"\033[36m", // Cyan
}
