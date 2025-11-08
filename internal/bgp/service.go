package bgp

import (
	"context"
	"fmt"
	"time"

	"github.com/padminisys/flintroute/internal/database"
	"github.com/padminisys/flintroute/internal/frr"
	"github.com/padminisys/flintroute/internal/models"
	"github.com/padminisys/flintroute/internal/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service manages BGP operations
type Service struct {
	db        *database.DB
	frrClient *frr.Client
	wsHub     *websocket.Hub
	logger    *zap.Logger
}

// NewService creates a new BGP service
func NewService(db *database.DB, frrClient *frr.Client, wsHub *websocket.Hub, logger *zap.Logger) *Service {
	return &Service{
		db:        db,
		frrClient: frrClient,
		wsHub:     wsHub,
		logger:    logger,
	}
}

// CreatePeer creates a new BGP peer
func (s *Service) CreatePeer(ctx context.Context, peer *models.BGPPeer) error {
	// Save to database
	if err := s.db.Create(peer).Error; err != nil {
		return fmt.Errorf("failed to create peer in database: %w", err)
	}

	// Configure in FRR if enabled
	if peer.Enabled {
		config := &frr.BGPPeerConfig{
			IPAddress:       peer.IPAddress,
			ASN:             peer.ASN,
			RemoteASN:       peer.RemoteASN,
			Password:        peer.Password,
			Multihop:        peer.Multihop,
			UpdateSource:    peer.UpdateSource,
			RouteMapIn:      peer.RouteMapIn,
			RouteMapOut:     peer.RouteMapOut,
			PrefixListIn:    peer.PrefixListIn,
			PrefixListOut:   peer.PrefixListOut,
			MaxPrefixes:     peer.MaxPrefixes,
			LocalPreference: peer.LocalPreference,
		}

		if err := s.frrClient.AddBGPPeer(ctx, config); err != nil {
			s.logger.Error("Failed to add peer to FRR", zap.Error(err))
			// Don't fail the operation, just log the error
		}
	}

	// Broadcast update
	s.wsHub.BroadcastPeerUpdate(peer)

	s.logger.Info("Created BGP peer",
		zap.Uint("id", peer.ID),
		zap.String("ip", peer.IPAddress),
	)

	return nil
}

// GetPeer retrieves a BGP peer by ID
func (s *Service) GetPeer(ctx context.Context, id uint) (*models.BGPPeer, error) {
	var peer models.BGPPeer
	if err := s.db.First(&peer, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("peer not found")
		}
		return nil, err
	}
	return &peer, nil
}

// ListPeers retrieves all BGP peers
func (s *Service) ListPeers(ctx context.Context) ([]*models.BGPPeer, error) {
	var peers []*models.BGPPeer
	if err := s.db.Find(&peers).Error; err != nil {
		return nil, err
	}
	return peers, nil
}

// UpdatePeer updates a BGP peer
func (s *Service) UpdatePeer(ctx context.Context, id uint, updates *models.BGPPeer) error {
	var peer models.BGPPeer
	if err := s.db.First(&peer, id).Error; err != nil {
		return fmt.Errorf("peer not found")
	}

	// Update fields
	peer.Name = updates.Name
	peer.Description = updates.Description
	peer.Enabled = updates.Enabled
	peer.Password = updates.Password
	peer.Multihop = updates.Multihop
	peer.UpdateSource = updates.UpdateSource
	peer.RouteMapIn = updates.RouteMapIn
	peer.RouteMapOut = updates.RouteMapOut
	peer.PrefixListIn = updates.PrefixListIn
	peer.PrefixListOut = updates.PrefixListOut
	peer.MaxPrefixes = updates.MaxPrefixes
	peer.LocalPreference = updates.LocalPreference

	if err := s.db.Save(&peer).Error; err != nil {
		return fmt.Errorf("failed to update peer: %w", err)
	}

	// Update FRR configuration
	config := &frr.BGPPeerConfig{
		IPAddress:       peer.IPAddress,
		ASN:             peer.ASN,
		RemoteASN:       peer.RemoteASN,
		Password:        peer.Password,
		Multihop:        peer.Multihop,
		UpdateSource:    peer.UpdateSource,
		RouteMapIn:      peer.RouteMapIn,
		RouteMapOut:     peer.RouteMapOut,
		PrefixListIn:    peer.PrefixListIn,
		PrefixListOut:   peer.PrefixListOut,
		MaxPrefixes:     peer.MaxPrefixes,
		LocalPreference: peer.LocalPreference,
	}

	if err := s.frrClient.UpdateBGPPeer(ctx, config); err != nil {
		s.logger.Error("Failed to update peer in FRR", zap.Error(err))
	}

	// Broadcast update
	s.wsHub.BroadcastPeerUpdate(&peer)

	s.logger.Info("Updated BGP peer", zap.Uint("id", id))

	return nil
}

// DeletePeer deletes a BGP peer
func (s *Service) DeletePeer(ctx context.Context, id uint) error {
	var peer models.BGPPeer
	if err := s.db.First(&peer, id).Error; err != nil {
		return fmt.Errorf("peer not found")
	}

	// Remove from FRR
	if err := s.frrClient.RemoveBGPPeer(ctx, peer.IPAddress); err != nil {
		s.logger.Error("Failed to remove peer from FRR", zap.Error(err))
	}

	// Delete from database
	if err := s.db.Delete(&peer).Error; err != nil {
		return fmt.Errorf("failed to delete peer: %w", err)
	}

	s.logger.Info("Deleted BGP peer", zap.Uint("id", id))

	return nil
}

// GetSession retrieves a BGP session by peer ID
func (s *Service) GetSession(ctx context.Context, peerID uint) (*models.BGPSession, error) {
	var session models.BGPSession
	if err := s.db.Preload("Peer").Where("peer_id = ?", peerID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	return &session, nil
}

// ListSessions retrieves all BGP sessions
func (s *Service) ListSessions(ctx context.Context) ([]*models.BGPSession, error) {
	var sessions []*models.BGPSession
	if err := s.db.Preload("Peer").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// UpdateSessionStates updates all BGP session states from FRR
func (s *Service) UpdateSessionStates(ctx context.Context) error {
	// Get all peers
	peers, err := s.ListPeers(ctx)
	if err != nil {
		return err
	}

	for _, peer := range peers {
		if !peer.Enabled {
			continue
		}

		// Get session state from FRR
		state, err := s.frrClient.GetBGPSessionState(ctx, peer.IPAddress)
		if err != nil {
			s.logger.Error("Failed to get session state",
				zap.String("ip", peer.IPAddress),
				zap.Error(err),
			)
			continue
		}

		// Update or create session in database
		var session models.BGPSession
		result := s.db.Where("peer_id = ?", peer.ID).First(&session)
		
		if result.Error == gorm.ErrRecordNotFound {
			// Create new session
			session = models.BGPSession{
				PeerID:           peer.ID,
				State:            state.State,
				Uptime:           state.Uptime,
				PrefixesReceived: state.PrefixesReceived,
				PrefixesSent:     state.PrefixesSent,
				MessagesReceived: state.MessagesReceived,
				MessagesSent:     state.MessagesSent,
				LastError:        state.LastError,
			}
			if err := s.db.Create(&session).Error; err != nil {
				s.logger.Error("Failed to create session", zap.Error(err))
				continue
			}
		} else {
			// Update existing session
			oldState := session.State
			session.State = state.State
			session.Uptime = state.Uptime
			session.PrefixesReceived = state.PrefixesReceived
			session.PrefixesSent = state.PrefixesSent
			session.MessagesReceived = state.MessagesReceived
			session.MessagesSent = state.MessagesSent
			session.LastError = state.LastError

			if err := s.db.Save(&session).Error; err != nil {
				s.logger.Error("Failed to update session", zap.Error(err))
				continue
			}

			// Create alert if state changed
			if oldState != state.State {
				s.createStateChangeAlert(peer, oldState, state.State)
			}
		}

		// Broadcast session update
		session.Peer = *peer
		s.wsHub.BroadcastSessionUpdate(&session)
	}

	return nil
}

// createStateChangeAlert creates an alert for BGP state changes
func (s *Service) createStateChangeAlert(peer *models.BGPPeer, oldState, newState string) {
	severity := "info"
	alertType := "peer_up"

	if newState != "Established" {
		severity = "warning"
		alertType = "peer_down"
	}

	alert := models.Alert{
		Type:     alertType,
		Severity: severity,
		Message:  fmt.Sprintf("BGP peer %s (%s) state changed from %s to %s", peer.Name, peer.IPAddress, oldState, newState),
		PeerID:   &peer.ID,
	}

	if err := s.db.Create(&alert).Error; err != nil {
		s.logger.Error("Failed to create alert", zap.Error(err))
		return
	}

	// Broadcast alert
	alert.Peer = peer
	s.wsHub.BroadcastAlert(&alert)

	s.logger.Info("Created state change alert",
		zap.String("peer", peer.Name),
		zap.String("old_state", oldState),
		zap.String("new_state", newState),
	)
}

// GetRunningConfig retrieves the current FRR running configuration
func (s *Service) GetRunningConfig(ctx context.Context) (string, error) {
	return s.frrClient.GetRunningConfig(ctx)
}

// StartMonitoring starts periodic monitoring of BGP sessions
func (s *Service) StartMonitoring(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.logger.Info("Started BGP session monitoring", zap.Duration("interval", interval))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopped BGP session monitoring")
			return
		case <-ticker.C:
			if err := s.UpdateSessionStates(ctx); err != nil {
				s.logger.Error("Failed to update session states", zap.Error(err))
			}
		}
	}
}