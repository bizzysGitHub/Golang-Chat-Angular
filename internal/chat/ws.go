package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader turns an HTTP request into a WebSocket connection.
// NOTE: CheckOrigin is wide-open for local dev. Tighten this in prod.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS upgrades the HTTP connection and registers a new Client.
func ServeWS(m *Manager, w http.ResponseWriter, r *http.Request) {
	// 1) Upgrade HTTP -> WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}

	// 2) Identify the client (dev: use remote addr; swap for userID/JWT later)
	id := r.RemoteAddr

	// 3) Create client and register with the manager
	client := NewClient(id, conn)
	m.Register() <- client

	// 4) Start the per-client goroutines (bodies fill in next step)
	go client.read(m)
	go client.write(m)

	log.Printf("ws connected: %s", id)
}
