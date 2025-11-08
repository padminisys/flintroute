package websocket

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

// Message represents a WebSocket message
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a WebSocket client
type Client struct {
	hub  *Hub
	send chan []byte
	id   string
}

// Hub maintains active WebSocket connections
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *zap.Logger
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Info("WebSocket client connected", zap.String("client_id", client.id))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.logger.Info("WebSocket client disconnected", zap.String("client_id", client.id))
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, close it
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(msgType string, payload interface{}) error {
	msg := Message{
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.broadcast <- data
	return nil
}

// BroadcastSessionUpdate sends a BGP session update to all clients
func (h *Hub) BroadcastSessionUpdate(session interface{}) error {
	return h.Broadcast("session_update", session)
}

// BroadcastAlert sends an alert to all clients
func (h *Hub) BroadcastAlert(alert interface{}) error {
	return h.Broadcast("alert", alert)
}

// BroadcastPeerUpdate sends a peer update to all clients
func (h *Hub) BroadcastPeerUpdate(peer interface{}) error {
	return h.Broadcast("peer_update", peer)
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}