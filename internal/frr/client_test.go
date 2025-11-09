package frr

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Create new client", func(t *testing.T) {
		client, err := NewClient("localhost", 50051, logger)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "localhost", client.host)
		assert.Equal(t, 50051, client.port)
		assert.NotNil(t, client.logger)
	})

	t.Run("Create client with different host and port", func(t *testing.T) {
		client, err := NewClient("frr-server", 9090, logger)
		assert.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "frr-server", client.host)
		assert.Equal(t, 9090, client.port)
	})
}

func TestIsConnected(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Not connected initially", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)
		assert.False(t, client.IsConnected())
	})

	// Note: Testing actual connection requires a running FRR gRPC server
	// This is tested via integration tests, not unit tests
}

func TestClose(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Close without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)
		err := client.Close()
		assert.NoError(t, err)
	})
}

func TestBGPPeerConfig(t *testing.T) {
	t.Run("Create BGP peer config", func(t *testing.T) {
		config := &BGPPeerConfig{
			IPAddress:       "192.168.1.1",
			ASN:             65001,
			RemoteASN:       65002,
			Password:        "secret",
			Multihop:        2,
			UpdateSource:    "eth0",
			RouteMapIn:      "RM-IN",
			RouteMapOut:     "RM-OUT",
			PrefixListIn:    "PL-IN",
			PrefixListOut:   "PL-OUT",
			MaxPrefixes:     1000,
			LocalPreference: 100,
		}

		assert.Equal(t, "192.168.1.1", config.IPAddress)
		assert.Equal(t, uint32(65001), config.ASN)
		assert.Equal(t, uint32(65002), config.RemoteASN)
		assert.Equal(t, "secret", config.Password)
		assert.Equal(t, 2, config.Multihop)
		assert.Equal(t, "eth0", config.UpdateSource)
		assert.Equal(t, "RM-IN", config.RouteMapIn)
		assert.Equal(t, "RM-OUT", config.RouteMapOut)
		assert.Equal(t, "PL-IN", config.PrefixListIn)
		assert.Equal(t, "PL-OUT", config.PrefixListOut)
		assert.Equal(t, 1000, config.MaxPrefixes)
		assert.Equal(t, 100, config.LocalPreference)
	})
}

func TestBGPSessionState(t *testing.T) {
	t.Run("Create BGP session state", func(t *testing.T) {
		state := &BGPSessionState{
			IPAddress:        "192.168.1.1",
			State:            "Established",
			Uptime:           3600,
			PrefixesReceived: 100,
			PrefixesSent:     50,
			MessagesReceived: 1000,
			MessagesSent:     900,
			LastError:        "",
		}

		assert.Equal(t, "192.168.1.1", state.IPAddress)
		assert.Equal(t, "Established", state.State)
		assert.Equal(t, int64(3600), state.Uptime)
		assert.Equal(t, 100, state.PrefixesReceived)
		assert.Equal(t, 50, state.PrefixesSent)
		assert.Equal(t, int64(1000), state.MessagesReceived)
		assert.Equal(t, int64(900), state.MessagesSent)
		assert.Equal(t, "", state.LastError)
	})

	t.Run("Session state with error", func(t *testing.T) {
		state := &BGPSessionState{
			IPAddress: "192.168.1.2",
			State:     "Idle",
			LastError: "Connection refused",
		}

		assert.Equal(t, "Idle", state.State)
		assert.Equal(t, "Connection refused", state.LastError)
	})
}

func TestAddBGPPeer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Add peer without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)
		config := &BGPPeerConfig{
			IPAddress: "192.168.1.1",
			ASN:       65001,
			RemoteASN: 65002,
		}

		err := client.AddBGPPeer(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
	// Use mock client for testing business logic
}

func TestRemoveBGPPeer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Remove peer without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)

		err := client.RemoveBGPPeer(ctx, "192.168.1.1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
}

func TestUpdateBGPPeer(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Update peer without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)
		config := &BGPPeerConfig{
			IPAddress: "192.168.1.1",
			ASN:       65001,
			RemoteASN: 65002,
		}

		err := client.UpdateBGPPeer(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
}

func TestGetBGPSessionState(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Get session state without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)

		state, err := client.GetBGPSessionState(ctx, "192.168.1.1")
		assert.Error(t, err)
		assert.Nil(t, state)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
}

func TestGetAllBGPSessions(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Get all sessions without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)

		sessions, err := client.GetAllBGPSessions(ctx)
		assert.Error(t, err)
		assert.Nil(t, sessions)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
}

func TestGetRunningConfig(t *testing.T) {
	logger := zap.NewNop()
	ctx := context.Background()

	t.Run("Get config without connection", func(t *testing.T) {
		client, _ := NewClient("localhost", 50051, logger)

		config, err := client.GetRunningConfig(ctx)
		assert.Error(t, err)
		assert.Empty(t, config)
		assert.Contains(t, err.Error(), "not connected")
	})

	// Note: Testing with actual connection requires a running FRR gRPC server
}

func TestMockClient(t *testing.T) {
	ctx := context.Background()

	t.Run("Mock Connect", func(t *testing.T) {
		mockClient := NewMockClient()
		mockClient.On("Connect", ctx).Return(nil)

		err := mockClient.Connect(ctx)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Mock IsConnected", func(t *testing.T) {
		mockClient := NewMockClient()
		mockClient.On("IsConnected").Return(true)

		connected := mockClient.IsConnected()
		assert.True(t, connected)
		mockClient.AssertExpectations(t)
	})

	t.Run("Mock AddBGPPeer", func(t *testing.T) {
		mockClient := NewMockClient()
		config := &BGPPeerConfig{
			IPAddress: "192.168.1.1",
			ASN:       65001,
			RemoteASN: 65002,
		}
		mockClient.On("AddBGPPeer", ctx, config).Return(nil)

		err := mockClient.AddBGPPeer(ctx, config)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Mock GetBGPSessionState", func(t *testing.T) {
		mockClient := NewMockClient()
		expectedState := &BGPSessionState{
			IPAddress: "192.168.1.1",
			State:     "Established",
			Uptime:    3600,
		}
		mockClient.On("GetBGPSessionState", ctx, "192.168.1.1").Return(expectedState, nil)

		state, err := mockClient.GetBGPSessionState(ctx, "192.168.1.1")
		assert.NoError(t, err)
		assert.Equal(t, expectedState, state)
		mockClient.AssertExpectations(t)
	})

	t.Run("Mock GetRunningConfig", func(t *testing.T) {
		mockClient := NewMockClient()
		expectedConfig := "router bgp 65001"
		mockClient.On("GetRunningConfig", ctx).Return(expectedConfig, nil)

		config, err := mockClient.GetRunningConfig(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, config)
		mockClient.AssertExpectations(t)
	})
}

func TestConnectTimeout(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Connect with timeout", func(t *testing.T) {
		client, _ := NewClient("invalid-host", 50051, logger)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := client.Connect(ctx)
		assert.Error(t, err)
	})
}