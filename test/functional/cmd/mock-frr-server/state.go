package main

import (
	"fmt"
	"sync"
	"time"
)

// BGPState manages the in-memory state of BGP peers and sessions
type BGPState struct {
	mu       sync.RWMutex
	peers    map[string]*PeerState
	sessions map[string]*SessionState
}

// PeerState represents the configuration state of a BGP peer
type PeerState struct {
	IPAddress       string
	ASN             uint32
	RemoteASN       uint32
	Password        string
	Multihop        int32
	UpdateSource    string
	RouteMapIn      string
	RouteMapOut     string
	PrefixListIn    string
	PrefixListOut   string
	MaxPrefixes     int32
	LocalPreference int32
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SessionState represents the runtime state of a BGP session
type SessionState struct {
	IPAddress        string
	State            string
	Uptime           int64
	PrefixesReceived int32
	PrefixesSent     int32
	MessagesReceived int64
	MessagesSent     int64
	LastError        string
	StateChangedAt   time.Time
}

// BGP session states
const (
	StateIdle        = "Idle"
	StateConnect     = "Connect"
	StateActive      = "Active"
	StateOpenSent    = "OpenSent"
	StateOpenConfirm = "OpenConfirm"
	StateEstablished = "Established"
)

// NewBGPState creates a new BGP state manager
func NewBGPState() *BGPState {
	return &BGPState{
		peers:    make(map[string]*PeerState),
		sessions: make(map[string]*SessionState),
	}
}

// AddPeer adds a new BGP peer to the state
func (s *BGPState) AddPeer(peer *PeerState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.peers[peer.IPAddress]; exists {
		return fmt.Errorf("peer %s already exists", peer.IPAddress)
	}

	now := time.Now()
	peer.CreatedAt = now
	peer.UpdatedAt = now

	s.peers[peer.IPAddress] = peer

	// Initialize session state
	session := &SessionState{
		IPAddress:        peer.IPAddress,
		State:            StateIdle,
		Uptime:           0,
		PrefixesReceived: 0,
		PrefixesSent:     0,
		MessagesReceived: 0,
		MessagesSent:     0,
		LastError:        "",
		StateChangedAt:   now,
	}
	s.sessions[peer.IPAddress] = session

	return nil
}

// RemovePeer removes a BGP peer from the state
func (s *BGPState) RemovePeer(ipAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.peers[ipAddress]; !exists {
		return fmt.Errorf("peer %s not found", ipAddress)
	}

	delete(s.peers, ipAddress)
	delete(s.sessions, ipAddress)

	return nil
}

// UpdatePeer updates an existing BGP peer configuration
func (s *BGPState) UpdatePeer(peer *PeerState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.peers[peer.IPAddress]
	if !exists {
		return fmt.Errorf("peer %s not found", peer.IPAddress)
	}

	// Preserve creation time
	peer.CreatedAt = existing.CreatedAt
	peer.UpdatedAt = time.Now()

	s.peers[peer.IPAddress] = peer

	return nil
}

// GetPeer retrieves a BGP peer by IP address
func (s *BGPState) GetPeer(ipAddress string) (*PeerState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, exists := s.peers[ipAddress]
	if !exists {
		return nil, fmt.Errorf("peer %s not found", ipAddress)
	}

	// Return a copy to prevent external modifications
	peerCopy := *peer
	return &peerCopy, nil
}

// GetAllPeers retrieves all BGP peers
func (s *BGPState) GetAllPeers() []*PeerState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers := make([]*PeerState, 0, len(s.peers))
	for _, peer := range s.peers {
		peerCopy := *peer
		peers = append(peers, &peerCopy)
	}

	return peers
}

// GetSessionState retrieves the session state for a peer
func (s *BGPState) GetSessionState(ipAddress string) (*SessionState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[ipAddress]
	if !exists {
		return nil, fmt.Errorf("session for peer %s not found", ipAddress)
	}

	// Calculate uptime if session is established
	sessionCopy := *session
	if session.State == StateEstablished {
		sessionCopy.Uptime = int64(time.Since(session.StateChangedAt).Seconds())
	}

	return &sessionCopy, nil
}

// GetAllSessions retrieves all BGP session states
func (s *BGPState) GetAllSessions() []*SessionState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sessions := make([]*SessionState, 0, len(s.sessions))
	now := time.Now()

	for _, session := range s.sessions {
		sessionCopy := *session
		// Calculate uptime if session is established
		if session.State == StateEstablished {
			sessionCopy.Uptime = int64(now.Sub(session.StateChangedAt).Seconds())
		}
		sessions = append(sessions, &sessionCopy)
	}

	return sessions
}

// UpdateSessionState updates the session state for a peer
func (s *BGPState) UpdateSessionState(ipAddress, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[ipAddress]
	if !exists {
		return fmt.Errorf("session for peer %s not found", ipAddress)
	}

	// Only update if state actually changed
	if session.State != state {
		session.State = state
		session.StateChangedAt = time.Now()

		// Reset uptime when transitioning to non-established states
		if state != StateEstablished {
			session.Uptime = 0
		}
	}

	return nil
}

// SimulateSessionEstablishment simulates the BGP session establishment process
func (s *BGPState) SimulateSessionEstablishment(ipAddress string, delay time.Duration) {
	states := []string{StateConnect, StateActive, StateOpenSent, StateOpenConfirm, StateEstablished}

	go func() {
		for _, state := range states {
			time.Sleep(delay)
			s.mu.Lock()
			if session, exists := s.sessions[ipAddress]; exists {
				session.State = state
				session.StateChangedAt = time.Now()

				// Simulate some traffic when established
				if state == StateEstablished {
					session.PrefixesReceived = 100
					session.PrefixesSent = 50
					session.MessagesReceived = 1000
					session.MessagesSent = 900
				}
			}
			s.mu.Unlock()
		}
	}()
}

// IncrementSessionCounters increments message counters for a session
func (s *BGPState) IncrementSessionCounters(ipAddress string, received, sent int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[ipAddress]
	if !exists {
		return fmt.Errorf("session for peer %s not found", ipAddress)
	}

	session.MessagesReceived += received
	session.MessagesSent += sent

	return nil
}

// SetSessionError sets an error message for a session
func (s *BGPState) SetSessionError(ipAddress, errorMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[ipAddress]
	if !exists {
		return fmt.Errorf("session for peer %s not found", ipAddress)
	}

	session.LastError = errorMsg
	session.State = StateIdle
	session.StateChangedAt = time.Now()

	return nil
}

// GetPeerCount returns the number of configured peers
func (s *BGPState) GetPeerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peers)
}

// GetEstablishedSessionCount returns the number of established sessions
func (s *BGPState) GetEstablishedSessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, session := range s.sessions {
		if session.State == StateEstablished {
			count++
		}
	}
	return count
}