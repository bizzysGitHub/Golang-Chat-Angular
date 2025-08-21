package chat

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// Upgrader turns an HTTP request into a WebSocket connection.
// NOTE: CheckOrigin is wide-open for local dev. Tighten this in prod.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:4200" // allow only Angular dev server
	},
}

// ServeWS upgrades the HTTP connection and registers a new Client.
func ServeWS(m *Manager, w http.ResponseWriter, r *http.Request) {
	// 1) Extract token (cookie or header)
	var token string
	if c, err := r.Cookie("access_token"); err == nil {
		token = c.Value
	} else if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		token = strings.TrimPrefix(h, "Bearer ")
	} else {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	// 2) Verify token → get userID
	userID, err := VerifyJWT(token) // implement with your signing key
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// 3) Check origin
	origin := r.Header.Get("Origin")
	if origin != "http://localhost:4200" {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}

	// 4) Upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrade failed", http.StatusBadRequest)
		return
	}

	// 5) Create client with authenticated ID
	client := NewClient(userID, conn)
	m.Register() <- client

	go client.read(m)
	go client.write(m)

}

// auth.go (internal/chat or internal/auth)
func VerifyJWT(token string) (userID string, err error) {
	// Use a JWT lib like github.com/golang-jwt/jwt/v5
	// Parse, validate signature/expiry, extract sub/uid claim.
	return "user-123", nil
}
