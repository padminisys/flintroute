package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// MockFRRServer implements a mock FRR gRPC service
type MockFRRServer struct {
	state      *BGPState
	config     *ServerConfig
	logger     *zap.Logger
	grpcServer *grpc.Server
	httpServer *http.Server
}

// NewMockFRRServer creates a new mock FRR server instance
func NewMockFRRServer(config *ServerConfig, logger *zap.Logger) *MockFRRServer {
	return &MockFRRServer{
		state:  NewBGPState(),
		config: config,
		logger: logger,
	}
}

// Start starts the mock FRR server
func (s *MockFRRServer) Start() error {
	// Create gRPC server
	s.grpcServer = grpc.NewServer()

	// Start gRPC listener
	lis, err := net.Listen("tcp", s.config.GetAddress())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("Mock FRR server starting",
		zap.String("address", s.config.GetAddress()),
	)

	// Start HTTP server for testing/debugging
	go s.startHTTPServer()

	// Start gRPC server
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

// Stop stops the mock FRR server
func (s *MockFRRServer) Stop() {
	s.logger.Info("Stopping mock FRR server")

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	if s.httpServer != nil {
		ctx := context.Background()
		s.httpServer.Shutdown(ctx)
	}
}

// startHTTPServer starts an HTTP server for testing and debugging
func (s *MockFRRServer) startHTTPServer() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", s.handleHealth)

	// Stats endpoint
	mux.HandleFunc("/stats", s.handleStats)

	// Peer management endpoints
	mux.HandleFunc("/peers", s.handlePeers)
	mux.HandleFunc("/peers/add", s.handleAddPeer)
	mux.HandleFunc("/peers/remove", s.handleRemovePeer)
	mux.HandleFunc("/peers/update", s.handleUpdatePeer)

	// Session endpoints
	mux.HandleFunc("/sessions", s.handleGetAllSessions)
	mux.HandleFunc("/sessions/state", s.handleGetSessionState)

	// Config endpoint
	mux.HandleFunc("/config", s.handleGetConfig)

	httpPort := s.config.Server.Port + 1000 // HTTP on port+1000
	httpAddr := fmt.Sprintf("%s:%d", s.config.Server.Host, httpPort)

	s.httpServer = &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	s.logger.Info("HTTP debug server starting", zap.String("address", httpAddr))

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("HTTP server error", zap.Error(err))
	}
}

// HTTP Handlers

func (s *MockFRRServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"peers":  s.state.GetPeerCount(),
	})
}

func (s *MockFRRServer) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := map[string]interface{}{
		"total_peers":          s.state.GetPeerCount(),
		"established_sessions": s.state.GetEstablishedSessionCount(),
	}
	json.NewEncoder(w).Encode(stats)
}

func (s *MockFRRServer) handlePeers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	peers := s.state.GetAllPeers()
	json.NewEncoder(w).Encode(peers)
}

func (s *MockFRRServer) handleAddPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var peer PeerState
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for error injection
	if s.config.Simulation.ErrorInjection {
		http.Error(w, "simulated error: failed to add peer", http.StatusInternalServerError)
		return
	}

	if err := s.state.AddPeer(&peer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Simulate session establishment
	s.state.SimulateSessionEstablishment(peer.IPAddress, s.config.Simulation.SessionStateDelay)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "peer added successfully",
	})
}

func (s *MockFRRServer) handleRemovePeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		IPAddress string `json:"ip_address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for error injection
	if s.config.Simulation.ErrorInjection {
		http.Error(w, "simulated error: failed to remove peer", http.StatusInternalServerError)
		return
	}

	if err := s.state.RemovePeer(req.IPAddress); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "peer removed successfully",
	})
}

func (s *MockFRRServer) handleUpdatePeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var peer PeerState
	if err := json.NewDecoder(r.Body).Decode(&peer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check for error injection
	if s.config.Simulation.ErrorInjection {
		http.Error(w, "simulated error: failed to update peer", http.StatusInternalServerError)
		return
	}

	if err := s.state.UpdatePeer(&peer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "peer updated successfully",
	})
}

func (s *MockFRRServer) handleGetAllSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sessions := s.state.GetAllSessions()
	json.NewEncoder(w).Encode(sessions)
}

func (s *MockFRRServer) handleGetSessionState(w http.ResponseWriter, r *http.Request) {
	ipAddress := r.URL.Query().Get("ip")
	if ipAddress == "" {
		http.Error(w, "ip parameter is required", http.StatusBadRequest)
		return
	}

	session, err := s.state.GetSessionState(ipAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

func (s *MockFRRServer) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config := s.generateMockConfig()
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(config))
}

// generateMockConfig generates a mock FRR configuration string
func (s *MockFRRServer) generateMockConfig() string {
	peers := s.state.GetAllPeers()

	config := "!\n"
	config += "! FRR Mock Configuration\n"
	config += "!\n"
	config += "frr version 8.0\n"
	config += "frr defaults traditional\n"
	config += "!\n"

	if len(peers) > 0 {
		config += "router bgp 65000\n"
		for _, peer := range peers {
			config += fmt.Sprintf(" neighbor %s remote-as %d\n", peer.IPAddress, peer.RemoteASN)

			if peer.Password != "" {
				config += fmt.Sprintf(" neighbor %s password %s\n", peer.IPAddress, peer.Password)
			}

			if peer.Multihop > 0 {
				config += fmt.Sprintf(" neighbor %s ebgp-multihop %d\n", peer.IPAddress, peer.Multihop)
			}

			if peer.UpdateSource != "" {
				config += fmt.Sprintf(" neighbor %s update-source %s\n", peer.IPAddress, peer.UpdateSource)
			}

			if peer.RouteMapIn != "" {
				config += fmt.Sprintf(" neighbor %s route-map %s in\n", peer.IPAddress, peer.RouteMapIn)
			}

			if peer.RouteMapOut != "" {
				config += fmt.Sprintf(" neighbor %s route-map %s out\n", peer.IPAddress, peer.RouteMapOut)
			}

			if peer.PrefixListIn != "" {
				config += fmt.Sprintf(" neighbor %s prefix-list %s in\n", peer.IPAddress, peer.PrefixListIn)
			}

			if peer.PrefixListOut != "" {
				config += fmt.Sprintf(" neighbor %s prefix-list %s out\n", peer.IPAddress, peer.PrefixListOut)
			}

			if peer.MaxPrefixes > 0 {
				config += fmt.Sprintf(" neighbor %s maximum-prefix %d\n", peer.IPAddress, peer.MaxPrefixes)
			}
		}
		config += "!\n"
	}

	config += "line vty\n"
	config += "!\n"
	config += "end\n"

	return config
}