package frr

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the FRR client for testing
type MockClient struct {
	mock.Mock
}

// NewMockClient creates a new mock FRR client
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Connect mocks the Connect method
func (m *MockClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Close mocks the Close method
func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

// IsConnected mocks the IsConnected method
func (m *MockClient) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

// AddBGPPeer mocks the AddBGPPeer method
func (m *MockClient) AddBGPPeer(ctx context.Context, config *BGPPeerConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// RemoveBGPPeer mocks the RemoveBGPPeer method
func (m *MockClient) RemoveBGPPeer(ctx context.Context, ipAddress string) error {
	args := m.Called(ctx, ipAddress)
	return args.Error(0)
}

// UpdateBGPPeer mocks the UpdateBGPPeer method
func (m *MockClient) UpdateBGPPeer(ctx context.Context, config *BGPPeerConfig) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

// GetBGPSessionState mocks the GetBGPSessionState method
func (m *MockClient) GetBGPSessionState(ctx context.Context, ipAddress string) (*BGPSessionState, error) {
	args := m.Called(ctx, ipAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BGPSessionState), args.Error(1)
}

// GetAllBGPSessions mocks the GetAllBGPSessions method
func (m *MockClient) GetAllBGPSessions(ctx context.Context) ([]*BGPSessionState, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*BGPSessionState), args.Error(1)
}

// GetRunningConfig mocks the GetRunningConfig method
func (m *MockClient) GetRunningConfig(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}