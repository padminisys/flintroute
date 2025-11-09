package websocket

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewHub(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Create new hub", func(t *testing.T) {
		hub := NewHub(logger)
		assert.NotNil(t, hub)
		assert.NotNil(t, hub.clients)
		assert.NotNil(t, hub.broadcast)
		assert.NotNil(t, hub.register)
		assert.NotNil(t, hub.unregister)
		assert.NotNil(t, hub.logger)
	})
}

func TestClientCount(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Initial client count is zero", func(t *testing.T) {
		hub := NewHub(logger)
		assert.Equal(t, 0, hub.ClientCount())
	})

	// Note: Testing hub.Run() requires careful channel management
	// These tests are simplified to avoid race conditions
}

func TestBroadcast(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Broadcast message", func(t *testing.T) {
		hub := NewHub(logger)

		err := hub.Broadcast("test_type", map[string]string{"key": "value"})
		assert.NoError(t, err)
	})

	// Note: Broadcast with running hub requires complex synchronization
	// Tested via integration tests

	t.Run("Broadcast with invalid payload", func(t *testing.T) {
		hub := NewHub(logger)

		// Create a payload that cannot be marshaled to JSON
		invalidPayload := make(chan int)

		err := hub.Broadcast("test_type", invalidPayload)
		assert.Error(t, err)
	})
}

func TestBroadcastSessionUpdate(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	t.Run("Broadcast session update", func(t *testing.T) {
		session := map[string]interface{}{
			"peer_id": 1,
			"state":   "Established",
		}

		err := hub.BroadcastSessionUpdate(session)
		assert.NoError(t, err)
	})
}

func TestBroadcastAlert(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	t.Run("Broadcast alert", func(t *testing.T) {
		alert := map[string]interface{}{
			"type":     "peer_down",
			"severity": "critical",
			"message":  "Peer is down",
		}

		err := hub.BroadcastAlert(alert)
		assert.NoError(t, err)
	})
}

func TestBroadcastPeerUpdate(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	t.Run("Broadcast peer update", func(t *testing.T) {
		peer := map[string]interface{}{
			"id":         1,
			"ip_address": "192.168.1.1",
			"enabled":    true,
		}

		err := hub.BroadcastPeerUpdate(peer)
		assert.NoError(t, err)
	})
}

func TestMessage(t *testing.T) {
	t.Run("Create message", func(t *testing.T) {
		msg := Message{
			Type:    "test_type",
			Payload: map[string]string{"key": "value"},
		}

		assert.Equal(t, "test_type", msg.Type)
		assert.NotNil(t, msg.Payload)
	})

	t.Run("Marshal message to JSON", func(t *testing.T) {
		msg := Message{
			Type: "test_event",
			Payload: map[string]interface{}{
				"message": "hello",
				"count":   123,
			},
		}

		data, err := json.Marshal(msg)
		assert.NoError(t, err)
		assert.Contains(t, string(data), "test_event")
		assert.Contains(t, string(data), "hello")
	})

	t.Run("Unmarshal JSON to message", func(t *testing.T) {
		jsonData := `{"type":"test_type","payload":{"key":"value"}}`

		var msg Message
		err := json.Unmarshal([]byte(jsonData), &msg)
		assert.NoError(t, err)
		assert.Equal(t, "test_type", msg.Type)
	})
}

func TestClient(t *testing.T) {
	logger := zap.NewNop()
	hub := NewHub(logger)

	t.Run("Create client", func(t *testing.T) {
		client := &Client{
			hub:  hub,
			send: make(chan []byte, 256),
			id:   "test-client",
		}

		assert.NotNil(t, client.hub)
		assert.NotNil(t, client.send)
		assert.Equal(t, "test-client", client.id)
	})
}

func TestHubRun(t *testing.T) {
	// Note: Hub.Run() tests require careful channel management
	// These are better suited for integration tests
	t.Skip("Hub.Run() tests are complex and better suited for integration tests")
}

func TestConcurrentOperations(t *testing.T) {
	// Note: Concurrent operations with hub.Run() are complex
	// Better tested in integration tests
	t.Skip("Concurrent operations are better suited for integration tests")
}