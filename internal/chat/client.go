package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

const (
	readWait   = 60 * time.Second // how long we wait to hear from peer
	writeWait  = 10 * time.Second // how long we allow each write
	pingPeriod = 54 * time.Second // must be < readWait so pongs arrive in time
)

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan []byte, 256), // small buffer
	}
}

// read: socket -> manager.Broadcast
func (c *Client) read(m *Manager) {
	defer func() {
		// ensure cleanup even on error/exit
		m.Unregister() <- c
		_ = c.Conn.Close()
	}()

	// 1) Set initial read deadline
	_ = c.Conn.SetReadDeadline(time.Now().Add(readWait))

	// 2) On each Pong from the peer, extend the deadline
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(readWait))
	})

	for {
		// block until client sends a WS frame
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			// connection closed or error
			log.Printf("read err (%s): %v", c.ID, err)
			return
		}

		// wrap incoming text as a Message { type:"chat", from, payload }
		msg := Message{
			Type: "chat",
			From: c.ID,
			// Payload is raw JSON; we’ll encode the string as JSON string
		}
		payload, _ := json.Marshal(ChatPayload{Text: string(data)})
		msg.Payload = payload

		out, err := json.Marshal(msg)
		if err != nil {
			log.Printf("marshal err (%s): %v", c.ID, err)
			continue
		}

		// hand off to manager for fan-out
		m.Broadcast() <- out
	}
}

// write: manager -> socket
func (c *Client) write(m *Manager) {
	ticker := time.NewTicker(pingPeriod) // NEW: heartbeat
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			// Set a write deadline for all writes
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// Channel closed by manager: tell peer we’re done
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Normal app message
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			// Periodic ping to keep the connection alive and detect dead peers
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
