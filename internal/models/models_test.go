package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&User{},
		&BGPPeer{},
		&BGPSession{},
		&ConfigVersion{},
		&Alert{},
		&RefreshToken{},
	)
	assert.NoError(t, err)

	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create user", func(t *testing.T) {
		user := User{
			Username:     "testuser",
			PasswordHash: "hashedpassword",
			Email:        "test@example.com",
			Role:         "admin",
			Active:       true,
		}

		err := db.Create(&user).Error
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotZero(t, user.CreatedAt)
	})

	t.Run("Unique username constraint", func(t *testing.T) {
		user1 := User{
			Username:     "uniqueuser",
			PasswordHash: "hash1",
			Email:        "user1@example.com",
			Role:         "user",
			Active:       true,
		}
		err := db.Create(&user1).Error
		assert.NoError(t, err)

		user2 := User{
			Username:     "uniqueuser",
			PasswordHash: "hash2",
			Email:        "user2@example.com",
			Role:         "user",
			Active:       true,
		}
		err = db.Create(&user2).Error
		assert.Error(t, err)
	})

	t.Run("Table name", func(t *testing.T) {
		user := User{}
		assert.Equal(t, "users", user.TableName())
	})
}

func TestBGPPeerModel(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Create BGP peer", func(t *testing.T) {
		peer := BGPPeer{
			Name:            "Test Peer",
			IPAddress:       "192.168.1.1",
			ASN:             65001,
			RemoteASN:       65002,
			Description:     "Test description",
			Enabled:         true,
			Multihop:        1,
			MaxPrefixes:     1000,
			LocalPreference: 100,
		}

		err := db.Create(&peer).Error
		assert.NoError(t, err)
		assert.NotZero(t, peer.ID)
		assert.Equal(t, "Test Peer", peer.Name)
		assert.Equal(t, "192.168.1.1", peer.IPAddress)
	})

	t.Run("Unique IP address constraint", func(t *testing.T) {
		peer1 := BGPPeer{
			Name:      "Peer 1",
			IPAddress: "10.0.0.1",
			ASN:       65001,
			RemoteASN: 65002,
			Enabled:   true,
		}
		err := db.Create(&peer1).Error
		assert.NoError(t, err)

		peer2 := BGPPeer{
			Name:      "Peer 2",
			IPAddress: "10.0.0.1",
			ASN:       65001,
			RemoteASN: 65003,
			Enabled:   true,
		}
		err = db.Create(&peer2).Error
		assert.Error(t, err)
	})

	t.Run("Table name", func(t *testing.T) {
		peer := BGPPeer{}
		assert.Equal(t, "bgp_peers", peer.TableName())
	})
}

func TestBGPSessionModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a peer first
	peer := BGPPeer{
		Name:      "Session Peer",
		IPAddress: "192.168.2.1",
		ASN:       65001,
		RemoteASN: 65002,
		Enabled:   true,
	}
	err := db.Create(&peer).Error
	assert.NoError(t, err)

	t.Run("Create BGP session", func(t *testing.T) {
		session := BGPSession{
			PeerID:           peer.ID,
			State:            "Established",
			Uptime:           3600,
			PrefixesReceived: 100,
			PrefixesSent:     50,
			MessagesReceived: 1000,
			MessagesSent:     900,
		}

		err := db.Create(&session).Error
		assert.NoError(t, err)
		assert.NotZero(t, session.ID)
		assert.Equal(t, "Established", session.State)
	})

	t.Run("Load session with peer", func(t *testing.T) {
		session := BGPSession{
			PeerID: peer.ID,
			State:  "Active",
		}
		db.Create(&session)

		var loadedSession BGPSession
		err := db.Preload("Peer").First(&loadedSession, session.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, peer.ID, loadedSession.Peer.ID)
		assert.Equal(t, "Session Peer", loadedSession.Peer.Name)
	})

	t.Run("Table name", func(t *testing.T) {
		session := BGPSession{}
		assert.Equal(t, "bgp_sessions", session.TableName())
	})
}

func TestConfigVersionModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a user first
	user := User{
		Username:     "configuser",
		PasswordHash: "hash",
		Email:        "config@example.com",
		Role:         "admin",
		Active:       true,
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)

	t.Run("Create config version", func(t *testing.T) {
		version := ConfigVersion{
			Description: "Test backup",
			Config:      "router bgp 65001",
			Hash:        "abc123",
			CreatedBy:   user.ID,
		}

		err := db.Create(&version).Error
		assert.NoError(t, err)
		assert.NotZero(t, version.ID)
		assert.Equal(t, "Test backup", version.Description)
	})

	t.Run("Unique hash constraint", func(t *testing.T) {
		version1 := ConfigVersion{
			Description: "Version 1",
			Config:      "config1",
			Hash:        "hash123",
			CreatedBy:   user.ID,
		}
		err := db.Create(&version1).Error
		assert.NoError(t, err)

		version2 := ConfigVersion{
			Description: "Version 2",
			Config:      "config2",
			Hash:        "hash123",
			CreatedBy:   user.ID,
		}
		err = db.Create(&version2).Error
		assert.Error(t, err)
	})

	t.Run("Table name", func(t *testing.T) {
		version := ConfigVersion{}
		assert.Equal(t, "config_versions", version.TableName())
	})
}

func TestAlertModel(t *testing.T) {
	db := setupTestDB(t)

	// Create peer and user
	peer := BGPPeer{
		Name:      "Alert Peer",
		IPAddress: "192.168.3.1",
		ASN:       65001,
		RemoteASN: 65002,
		Enabled:   true,
	}
	db.Create(&peer)

	user := User{
		Username:     "alertuser",
		PasswordHash: "hash",
		Email:        "alert@example.com",
		Role:         "admin",
		Active:       true,
	}
	db.Create(&user)

	t.Run("Create alert", func(t *testing.T) {
		alert := Alert{
			Type:     "peer_down",
			Severity: "critical",
			Message:  "Peer is down",
			Details:  "Connection lost",
			PeerID:   &peer.ID,
		}

		err := db.Create(&alert).Error
		assert.NoError(t, err)
		assert.NotZero(t, alert.ID)
		assert.Equal(t, "peer_down", alert.Type)
		assert.False(t, alert.Acknowledged)
	})

	t.Run("Acknowledge alert", func(t *testing.T) {
		alert := Alert{
			Type:     "peer_up",
			Severity: "info",
			Message:  "Peer is up",
		}
		db.Create(&alert)

		now := time.Now()
		alert.Acknowledged = true
		alert.AcknowledgedAt = &now
		alert.AcknowledgedBy = &user.ID

		err := db.Save(&alert).Error
		assert.NoError(t, err)
		assert.True(t, alert.Acknowledged)
		assert.NotNil(t, alert.AcknowledgedAt)
	})

	t.Run("Table name", func(t *testing.T) {
		alert := Alert{}
		assert.Equal(t, "alerts", alert.TableName())
	})
}

func TestRefreshTokenModel(t *testing.T) {
	db := setupTestDB(t)

	// Create user
	user := User{
		Username:     "tokenuser",
		PasswordHash: "hash",
		Email:        "token@example.com",
		Role:         "user",
		Active:       true,
	}
	db.Create(&user)

	t.Run("Create refresh token", func(t *testing.T) {
		token := RefreshToken{
			UserID:    user.ID,
			Token:     "refresh_token_123",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}

		err := db.Create(&token).Error
		assert.NoError(t, err)
		assert.NotZero(t, token.ID)
		assert.False(t, token.Revoked)
	})

	t.Run("Revoke token", func(t *testing.T) {
		token := RefreshToken{
			UserID:    user.ID,
			Token:     "token_to_revoke",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			Revoked:   false,
		}
		db.Create(&token)

		token.Revoked = true
		err := db.Save(&token).Error
		assert.NoError(t, err)
		assert.True(t, token.Revoked)
	})

	t.Run("Table name", func(t *testing.T) {
		token := RefreshToken{}
		assert.Equal(t, "refresh_tokens", token.TableName())
	})
}