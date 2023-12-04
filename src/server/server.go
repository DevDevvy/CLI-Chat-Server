// src/server/server.go
package server

import (
	"chat-server/src/websocket"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/ws", websocket.Handler)

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
