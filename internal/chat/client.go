package chat

import "github.com/gorilla/websocket"

// Client represents a single connected user/peer.
type Client struct {
	ID   string
	Conn *websocket.Conn // set during WebSocket upgrade
	Send chan []byte     // outbound messages to this client
}

// NewClient constructs a client with a buffered outbound channel.
func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
}

// read pumps messages from the WebSocket -> manager.broadcast.
// (We’ll implement this in the next step.)
func (c *Client) read(m *Manager) {
	// TODO: implement
}

// write pumps messages from c.Send -> WebSocket.
// (We’ll implement this in the next step.)
func (c *Client) write(m *Manager) {
	// TODO: implement
}
