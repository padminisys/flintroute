package testutil

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// FixtureLoader loads test fixtures from YAML files
type FixtureLoader struct {
	basePath string
	logger   *zap.Logger
}

// NewFixtureLoader creates a new fixture loader
func NewFixtureLoader(basePath string, logger *zap.Logger) *FixtureLoader {
	return &FixtureLoader{
		basePath: basePath,
		logger:   logger,
	}
}

// PeerFixture represents a peer fixture
type PeerFixture struct {
	Name            string `yaml:"name"`
	IPAddress       string `yaml:"ip_address"`
	ASN             uint32 `yaml:"asn"`
	RemoteASN       uint32 `yaml:"remote_asn"`
	Description     string `yaml:"description"`
	Enabled         bool   `yaml:"enabled"`
	Password        string `yaml:"password"`
	Multihop        int    `yaml:"multihop"`
	UpdateSource    string `yaml:"update_source"`
	RouteMapIn      string `yaml:"route_map_in"`
	RouteMapOut     string `yaml:"route_map_out"`
	PrefixListIn    string `yaml:"prefix_list_in"`
	PrefixListOut   string `yaml:"prefix_list_out"`
	MaxPrefixes     int    `yaml:"max_prefixes"`
	LocalPreference int    `yaml:"local_preference"`
}

// UserFixture represents a user fixture
type UserFixture struct {
	Username string `yaml:"username"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Role     string `yaml:"role"`
	Active   bool   `yaml:"active"`
}

// SessionFixture represents a session fixture
type SessionFixture struct {
	PeerID           uint   `yaml:"peer_id"`
	State            string `yaml:"state"`
	Uptime           int64  `yaml:"uptime"`
	PrefixesReceived int    `yaml:"prefixes_received"`
	PrefixesSent     int    `yaml:"prefixes_sent"`
	MessagesReceived int64  `yaml:"messages_received"`
	MessagesSent     int64  `yaml:"messages_sent"`
	LastError        string `yaml:"last_error"`
}

// LoadPeer loads a peer fixture by name
func (fl *FixtureLoader) LoadPeer(name string) (*PeerFixture, error) {
	path := filepath.Join(fl.basePath, "peers", name+".yaml")
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read peer fixture %s: %w", name, err)
	}

	var peer PeerFixture
	if err := yaml.Unmarshal(data, &peer); err != nil {
		return nil, fmt.Errorf("failed to parse peer fixture %s: %w", name, err)
	}

	fl.logger.Debug("Peer fixture loaded", zap.String("name", name))
	return &peer, nil
}

// LoadUser loads a user fixture by name
func (fl *FixtureLoader) LoadUser(name string) (*UserFixture, error) {
	path := filepath.Join(fl.basePath, "users", name+".yaml")
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read user fixture %s: %w", name, err)
	}

	var user UserFixture
	if err := yaml.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user fixture %s: %w", name, err)
	}

	fl.logger.Debug("User fixture loaded", zap.String("name", name))
	return &user, nil
}

// LoadSession loads a session fixture by name
func (fl *FixtureLoader) LoadSession(name string) (*SessionFixture, error) {
	path := filepath.Join(fl.basePath, "sessions", name+".yaml")
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read session fixture %s: %w", name, err)
	}

	var session SessionFixture
	if err := yaml.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse session fixture %s: %w", name, err)
	}

	fl.logger.Debug("Session fixture loaded", zap.String("name", name))
	return &session, nil
}

// LoadAllPeers loads all peer fixtures matching a pattern
func (fl *FixtureLoader) LoadAllPeers(pattern string) ([]*PeerFixture, error) {
	peersDir := filepath.Join(fl.basePath, "peers")
	
	// If pattern is empty, use wildcard
	if pattern == "" {
		pattern = "*.yaml"
	}
	
	matches, err := filepath.Glob(filepath.Join(peersDir, pattern))
	if err != nil {
		return nil, fmt.Errorf("failed to glob peer fixtures: %w", err)
	}

	var peers []*PeerFixture
	for _, match := range matches {
		// Get filename without extension
		name := filepath.Base(match)
		name = name[:len(name)-len(filepath.Ext(name))]
		
		peer, err := fl.LoadPeer(name)
		if err != nil {
			fl.logger.Warn("Failed to load peer fixture", zap.String("name", name), zap.Error(err))
			continue
		}
		
		peers = append(peers, peer)
	}

	fl.logger.Debug("Peer fixtures loaded", zap.Int("count", len(peers)))
	return peers, nil
}

// LoadAllUsers loads all user fixtures matching a pattern
func (fl *FixtureLoader) LoadAllUsers(pattern string) ([]*UserFixture, error) {
	usersDir := filepath.Join(fl.basePath, "users")
	
	if pattern == "" {
		pattern = "*.yaml"
	}
	
	matches, err := filepath.Glob(filepath.Join(usersDir, pattern))
	if err != nil {
		return nil, fmt.Errorf("failed to glob user fixtures: %w", err)
	}

	var users []*UserFixture
	for _, match := range matches {
		name := filepath.Base(match)
		name = name[:len(name)-len(filepath.Ext(name))]
		
		user, err := fl.LoadUser(name)
		if err != nil {
			fl.logger.Warn("Failed to load user fixture", zap.String("name", name), zap.Error(err))
			continue
		}
		
		users = append(users, user)
	}

	fl.logger.Debug("User fixtures loaded", zap.Int("count", len(users)))
	return users, nil
}

// LoadAllSessions loads all session fixtures matching a pattern
func (fl *FixtureLoader) LoadAllSessions(pattern string) ([]*SessionFixture, error) {
	sessionsDir := filepath.Join(fl.basePath, "sessions")
	
	if pattern == "" {
		pattern = "*.yaml"
	}
	
	matches, err := filepath.Glob(filepath.Join(sessionsDir, pattern))
	if err != nil {
		return nil, fmt.Errorf("failed to glob session fixtures: %w", err)
	}

	var sessions []*SessionFixture
	for _, match := range matches {
		name := filepath.Base(match)
		name = name[:len(name)-len(filepath.Ext(name))]
		
		session, err := fl.LoadSession(name)
		if err != nil {
			fl.logger.Warn("Failed to load session fixture", zap.String("name", name), zap.Error(err))
			continue
		}
		
		sessions = append(sessions, session)
	}

	fl.logger.Debug("Session fixtures loaded", zap.Int("count", len(sessions)))
	return sessions, nil
}