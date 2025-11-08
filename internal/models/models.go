package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Email        string         `gorm:"uniqueIndex" json:"email"`
	Role         string         `gorm:"not null;default:'user'" json:"role"` // admin, user
	Active       bool           `gorm:"not null;default:true" json:"active"`
}

// BGPPeer represents a BGP peer configuration
type BGPPeer struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string         `gorm:"not null" json:"name"`
	IPAddress       string         `gorm:"uniqueIndex;not null" json:"ip_address"`
	ASN             uint32         `gorm:"not null" json:"asn"`
	RemoteASN       uint32         `gorm:"not null" json:"remote_asn"`
	Description     string         `json:"description"`
	Enabled         bool           `gorm:"not null;default:true" json:"enabled"`
	Password        string         `json:"password,omitempty"`
	Multihop        int            `gorm:"default:1" json:"multihop"`
	UpdateSource    string         `json:"update_source"`
	RouteMapIn      string         `json:"route_map_in"`
	RouteMapOut     string         `json:"route_map_out"`
	PrefixListIn    string         `json:"prefix_list_in"`
	PrefixListOut   string         `json:"prefix_list_out"`
	MaxPrefixes     int            `json:"max_prefixes"`
	LocalPreference int            `json:"local_preference"`
}

// BGPSession represents the runtime state of a BGP session
type BGPSession struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PeerID           uint      `gorm:"not null;index" json:"peer_id"`
	Peer             BGPPeer   `gorm:"foreignKey:PeerID" json:"peer,omitempty"`
	State            string    `gorm:"not null" json:"state"` // Idle, Connect, Active, OpenSent, OpenConfirm, Established
	Uptime           int64     `json:"uptime"`                // seconds
	PrefixesReceived int       `json:"prefixes_received"`
	PrefixesSent     int       `json:"prefixes_sent"`
	MessagesReceived int64     `json:"messages_received"`
	MessagesSent     int64     `json:"messages_sent"`
	LastError        string    `json:"last_error"`
	LastReset        time.Time `json:"last_reset"`
}

// ConfigVersion represents a configuration backup
type ConfigVersion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Config      string    `gorm:"type:text;not null" json:"config"`
	Hash        string    `gorm:"uniqueIndex;not null" json:"hash"`
	CreatedBy   uint      `json:"created_by"`
	User        User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

// Alert represents a system alert
type Alert struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Type          string         `gorm:"not null;index" json:"type"` // peer_down, peer_up, config_change, etc.
	Severity      string         `gorm:"not null" json:"severity"`   // info, warning, error, critical
	Message       string         `gorm:"not null" json:"message"`
	Details       string         `gorm:"type:text" json:"details"`
	PeerID        *uint          `gorm:"index" json:"peer_id,omitempty"`
	Peer          *BGPPeer       `gorm:"foreignKey:PeerID" json:"peer,omitempty"`
	Acknowledged  bool           `gorm:"not null;default:false" json:"acknowledged"`
	AcknowledgedAt *time.Time    `json:"acknowledged_at,omitempty"`
	AcknowledgedBy *uint         `json:"acknowledged_by,omitempty"`
	User          *User          `gorm:"foreignKey:AcknowledgedBy" json:"user,omitempty"`
}

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	Revoked   bool      `gorm:"not null;default:false" json:"revoked"`
}

// TableName overrides for GORM
func (User) TableName() string          { return "users" }
func (BGPPeer) TableName() string       { return "bgp_peers" }
func (BGPSession) TableName() string    { return "bgp_sessions" }
func (ConfigVersion) TableName() string { return "config_versions" }
func (Alert) TableName() string         { return "alerts" }
func (RefreshToken) TableName() string  { return "refresh_tokens" }