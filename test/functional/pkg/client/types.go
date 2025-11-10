package client

import "time"

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo represents user information
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// TokenResponse represents a token refresh response
type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PeerRequest represents a request to create or update a BGP peer
type PeerRequest struct {
	Name            string `json:"name"`
	IPAddress       string `json:"ip_address"`
	ASN             uint32 `json:"asn"`
	RemoteASN       uint32 `json:"remote_asn"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	Password        string `json:"password,omitempty"`
	Multihop        int    `json:"multihop"`
	UpdateSource    string `json:"update_source,omitempty"`
	RouteMapIn      string `json:"route_map_in,omitempty"`
	RouteMapOut     string `json:"route_map_out,omitempty"`
	PrefixListIn    string `json:"prefix_list_in,omitempty"`
	PrefixListOut   string `json:"prefix_list_out,omitempty"`
	MaxPrefixes     int    `json:"max_prefixes"`
	LocalPreference int    `json:"local_preference"`
}

// Peer represents a BGP peer
type Peer struct {
	ID              uint      `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Name            string    `json:"name"`
	IPAddress       string    `json:"ip_address"`
	ASN             uint32    `json:"asn"`
	RemoteASN       uint32    `json:"remote_asn"`
	Description     string    `json:"description"`
	Enabled         bool      `json:"enabled"`
	Password        string    `json:"password,omitempty"`
	Multihop        int       `json:"multihop"`
	UpdateSource    string    `json:"update_source,omitempty"`
	RouteMapIn      string    `json:"route_map_in,omitempty"`
	RouteMapOut     string    `json:"route_map_out,omitempty"`
	PrefixListIn    string    `json:"prefix_list_in,omitempty"`
	PrefixListOut   string    `json:"prefix_list_out,omitempty"`
	MaxPrefixes     int       `json:"max_prefixes"`
	LocalPreference int       `json:"local_preference"`
}

// Session represents a BGP session
type Session struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PeerID           uint      `json:"peer_id"`
	Peer             *Peer     `json:"peer,omitempty"`
	State            string    `json:"state"`
	Uptime           int64     `json:"uptime"`
	PrefixesReceived int       `json:"prefixes_received"`
	PrefixesSent     int       `json:"prefixes_sent"`
	MessagesReceived int64     `json:"messages_received"`
	MessagesSent     int64     `json:"messages_sent"`
	LastError        string    `json:"last_error"`
	LastReset        time.Time `json:"last_reset"`
}

// ConfigVersion represents a configuration backup
type ConfigVersion struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Config      string    `json:"config"`
	Hash        string    `json:"hash"`
	CreatedBy   uint      `json:"created_by"`
	User        *UserInfo `json:"user,omitempty"`
}

// BackupConfigRequest represents a request to backup configuration
type BackupConfigRequest struct {
	Description string `json:"description"`
}

// Alert represents a system alert
type Alert struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Type           string     `json:"type"`
	Severity       string     `json:"severity"`
	Message        string     `json:"message"`
	Details        string     `json:"details"`
	PeerID         *uint      `json:"peer_id,omitempty"`
	Peer           *Peer      `json:"peer,omitempty"`
	Acknowledged   bool       `json:"acknowledged"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty"`
	AcknowledgedBy *uint      `json:"acknowledged_by,omitempty"`
	User           *UserInfo  `json:"user,omitempty"`
}

// AlertQueryParams represents query parameters for listing alerts
type AlertQueryParams struct {
	Acknowledged *bool  `json:"acknowledged,omitempty"`
	Severity     string `json:"severity,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// PeersResponse represents a list of peers response
type PeersResponse struct {
	Peers []*Peer `json:"peers"`
}

// SessionsResponse represents a list of sessions response
type SessionsResponse struct {
	Sessions []*Session `json:"sessions"`
}

// ConfigVersionsResponse represents a list of config versions response
type ConfigVersionsResponse struct {
	Versions []*ConfigVersion `json:"versions"`
}

// AlertsResponse represents a list of alerts response
type AlertsResponse struct {
	Alerts []*Alert `json:"alerts"`
}