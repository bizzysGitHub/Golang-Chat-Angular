package chat

import (
	"context"
	"log"
)

// Manager coordinates all clients and message fan-out.
type Manager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// NewManager constructs a single manager instance.
func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256), // small buffer for backpressure
	}
}

// Run processes register/unregister/broadcast until ctx is canceled.
func (m *Manager) Run(ctx context.Context) {
	for {
		select {
		case c := <-m.register:
			m.clients[c] = true
			log.Printf("client connected: %s", c.ID)

		case c := <-m.unregister:
			if _, ok := m.clients[c]; ok {
				delete(m.clients, c)
				close(c.Send) // signal writer to exit
				log.Printf("client disconnected: %s", c.ID)
			}

		case msg := <-m.broadcast:
			for c := range m.clients {
				select {
				case c.Send <- msg:
				default:
					// slow consumer: drop and disconnect to protect hub
					close(c.Send)
					delete(m.clients, c)
				}
			}

		case <-ctx.Done():
			// graceful shutdown: close all clients
			for c := range m.clients {
				close(c.Send)
				delete(m.clients, c)
			}
			return
		}
	}
}

// Expose channels via methods (so fields stay private & flexible).
func (m *Manager) Register() chan<- *Client   { return m.register }
func (m *Manager) Unregister() chan<- *Client { return m.unregister }
func (m *Manager) Broadcast() chan<- []byte   { return m.broadcast }
