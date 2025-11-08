package frr

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents an FRR gRPC client
type Client struct {
	conn   *grpc.ClientConn
	logger *zap.Logger
	host   string
	port   int
}

// NewClient creates a new FRR gRPC client
func NewClient(host string, port int, logger *zap.Logger) (*Client, error) {
	return &Client{
		host:   host,
		port:   port,
		logger: logger,
	}, nil
}

// Connect establishes connection to FRR gRPC server
func (c *Client) Connect(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to FRR gRPC server: %w", err)
	}

	c.conn = conn
	c.logger.Info("Connected to FRR gRPC server", zap.String("address", addr))
	return nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected checks if the client is connected
func (c *Client) IsConnected() bool {
	return c.conn != nil
}

// BGPPeerConfig represents BGP peer configuration for FRR
type BGPPeerConfig struct {
	IPAddress       string
	ASN             uint32
	RemoteASN       uint32
	Password        string
	Multihop        int
	UpdateSource    string
	RouteMapIn      string
	RouteMapOut     string
	PrefixListIn    string
	PrefixListOut   string
	MaxPrefixes     int
	LocalPreference int
}

// BGPSessionState represents BGP session state from FRR
type BGPSessionState struct {
	IPAddress        string
	State            string
	Uptime           int64
	PrefixesReceived int
	PrefixesSent     int
	MessagesReceived int64
	MessagesSent     int64
	LastError        string
}

// AddBGPPeer adds a BGP peer to FRR configuration
func (c *Client) AddBGPPeer(ctx context.Context, config *BGPPeerConfig) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	// For now, this is a stub that logs the operation
	c.logger.Info("Adding BGP peer",
		zap.String("ip", config.IPAddress),
		zap.Uint32("remote_asn", config.RemoteASN),
	)

	return nil
}

// RemoveBGPPeer removes a BGP peer from FRR configuration
func (c *Client) RemoveBGPPeer(ctx context.Context, ipAddress string) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	c.logger.Info("Removing BGP peer", zap.String("ip", ipAddress))

	return nil
}

// UpdateBGPPeer updates a BGP peer configuration
func (c *Client) UpdateBGPPeer(ctx context.Context, config *BGPPeerConfig) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	c.logger.Info("Updating BGP peer",
		zap.String("ip", config.IPAddress),
		zap.Uint32("remote_asn", config.RemoteASN),
	)

	return nil
}

// GetBGPSessionState retrieves BGP session state for a peer
func (c *Client) GetBGPSessionState(ctx context.Context, ipAddress string) (*BGPSessionState, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	// For now, return mock data
	c.logger.Debug("Getting BGP session state", zap.String("ip", ipAddress))

	return &BGPSessionState{
		IPAddress:        ipAddress,
		State:            "Established",
		Uptime:           3600,
		PrefixesReceived: 100,
		PrefixesSent:     50,
		MessagesReceived: 1000,
		MessagesSent:     900,
		LastError:        "",
	}, nil
}

// GetAllBGPSessions retrieves all BGP session states
func (c *Client) GetAllBGPSessions(ctx context.Context) ([]*BGPSessionState, error) {
	if !c.IsConnected() {
		return nil, fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	c.logger.Debug("Getting all BGP session states")

	return []*BGPSessionState{}, nil
}

// GetRunningConfig retrieves the current FRR running configuration
func (c *Client) GetRunningConfig(ctx context.Context) (string, error) {
	if !c.IsConnected() {
		return "", fmt.Errorf("not connected to FRR gRPC server")
	}

	// TODO: Implement actual gRPC call to FRR
	c.logger.Debug("Getting running configuration")

	return "! FRR Configuration\n", nil
}