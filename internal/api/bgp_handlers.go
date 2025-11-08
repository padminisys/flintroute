package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/padminisys/flintroute/internal/models"
	"go.uber.org/zap"
)

// CreatePeerRequest represents a request to create a BGP peer
type CreatePeerRequest struct {
	Name            string `json:"name" binding:"required"`
	IPAddress       string `json:"ip_address" binding:"required"`
	ASN             uint32 `json:"asn" binding:"required"`
	RemoteASN       uint32 `json:"remote_asn" binding:"required"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	Password        string `json:"password"`
	Multihop        int    `json:"multihop"`
	UpdateSource    string `json:"update_source"`
	RouteMapIn      string `json:"route_map_in"`
	RouteMapOut     string `json:"route_map_out"`
	PrefixListIn    string `json:"prefix_list_in"`
	PrefixListOut   string `json:"prefix_list_out"`
	MaxPrefixes     int    `json:"max_prefixes"`
	LocalPreference int    `json:"local_preference"`
}

// UpdatePeerRequest represents a request to update a BGP peer
type UpdatePeerRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Enabled         bool   `json:"enabled"`
	Password        string `json:"password"`
	Multihop        int    `json:"multihop"`
	UpdateSource    string `json:"update_source"`
	RouteMapIn      string `json:"route_map_in"`
	RouteMapOut     string `json:"route_map_out"`
	PrefixListIn    string `json:"prefix_list_in"`
	PrefixListOut   string `json:"prefix_list_out"`
	MaxPrefixes     int    `json:"max_prefixes"`
	LocalPreference int    `json:"local_preference"`
}

// handleListPeers handles listing all BGP peers
func (s *Server) handleListPeers(c *gin.Context) {
	peers, err := s.bgpService.ListPeers(c.Request.Context())
	if err != nil {
		s.logger.Error("Failed to list peers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list peers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"peers": peers})
}

// handleGetPeer handles getting a specific BGP peer
func (s *Server) handleGetPeer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid peer ID"})
		return
	}

	peer, err := s.bgpService.GetPeer(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peer not found"})
		return
	}

	c.JSON(http.StatusOK, peer)
}

// handleCreatePeer handles creating a new BGP peer
func (s *Server) handleCreatePeer(c *gin.Context) {
	var req CreatePeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	peer := &models.BGPPeer{
		Name:            req.Name,
		IPAddress:       req.IPAddress,
		ASN:             req.ASN,
		RemoteASN:       req.RemoteASN,
		Description:     req.Description,
		Enabled:         req.Enabled,
		Password:        req.Password,
		Multihop:        req.Multihop,
		UpdateSource:    req.UpdateSource,
		RouteMapIn:      req.RouteMapIn,
		RouteMapOut:     req.RouteMapOut,
		PrefixListIn:    req.PrefixListIn,
		PrefixListOut:   req.PrefixListOut,
		MaxPrefixes:     req.MaxPrefixes,
		LocalPreference: req.LocalPreference,
	}

	if err := s.bgpService.CreatePeer(c.Request.Context(), peer); err != nil {
		s.logger.Error("Failed to create peer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create peer"})
		return
	}

	c.JSON(http.StatusCreated, peer)
}

// handleUpdatePeer handles updating a BGP peer
func (s *Server) handleUpdatePeer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid peer ID"})
		return
	}

	var req UpdatePeerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	updates := &models.BGPPeer{
		Name:            req.Name,
		Description:     req.Description,
		Enabled:         req.Enabled,
		Password:        req.Password,
		Multihop:        req.Multihop,
		UpdateSource:    req.UpdateSource,
		RouteMapIn:      req.RouteMapIn,
		RouteMapOut:     req.RouteMapOut,
		PrefixListIn:    req.PrefixListIn,
		PrefixListOut:   req.PrefixListOut,
		MaxPrefixes:     req.MaxPrefixes,
		LocalPreference: req.LocalPreference,
	}

	if err := s.bgpService.UpdatePeer(c.Request.Context(), uint(id), updates); err != nil {
		s.logger.Error("Failed to update peer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update peer"})
		return
	}

	peer, _ := s.bgpService.GetPeer(c.Request.Context(), uint(id))
	c.JSON(http.StatusOK, peer)
}

// handleDeletePeer handles deleting a BGP peer
func (s *Server) handleDeletePeer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid peer ID"})
		return
	}

	if err := s.bgpService.DeletePeer(c.Request.Context(), uint(id)); err != nil {
		s.logger.Error("Failed to delete peer", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete peer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peer deleted successfully"})
}

// handleListSessions handles listing all BGP sessions
func (s *Server) handleListSessions(c *gin.Context) {
	sessions, err := s.bgpService.ListSessions(c.Request.Context())
	if err != nil {
		s.logger.Error("Failed to list sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// handleGetSession handles getting a specific BGP session
func (s *Server) handleGetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := s.bgpService.GetSession(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}