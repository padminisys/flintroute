package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	authpkg "github.com/padminisys/flintroute/internal/auth"
	"github.com/padminisys/flintroute/internal/models"
	"go.uber.org/zap"
)

// BackupConfigRequest represents a request to backup configuration
type BackupConfigRequest struct {
	Description string `json:"description"`
}

// handleListConfigVersions handles listing all configuration versions
func (s *Server) handleListConfigVersions(c *gin.Context) {
	var versions []models.ConfigVersion
	if err := s.db.Preload("User").Order("created_at DESC").Find(&versions).Error; err != nil {
		s.logger.Error("Failed to list config versions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list config versions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// handleBackupConfig handles backing up the current configuration
func (s *Server) handleBackupConfig(c *gin.Context) {
	var req BackupConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get current user ID
	userID, exists := authpkg.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current FRR configuration
	config, err := s.bgpService.GetRunningConfig(c.Request.Context())
	if err != nil {
		s.logger.Error("Failed to get running config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get running config"})
		return
	}

	// Calculate hash
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(config)))

	// Check if this config already exists
	var existingVersion models.ConfigVersion
	if err := s.db.Where("hash = ?", hash).First(&existingVersion).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Configuration already backed up",
			"version": existingVersion,
		})
		return
	}

	// Create new version
	version := models.ConfigVersion{
		Description: req.Description,
		Config:      config,
		Hash:        hash,
		CreatedBy:   userID,
	}

	if err := s.db.Create(&version).Error; err != nil {
		s.logger.Error("Failed to create config version", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to backup config"})
		return
	}

	// Load user info
	s.db.Preload("User").First(&version, version.ID)

	s.logger.Info("Configuration backed up",
		zap.Uint("version_id", version.ID),
		zap.Uint("user_id", userID),
	)

	c.JSON(http.StatusCreated, version)
}

// handleRestoreConfig handles restoring a configuration version
func (s *Server) handleRestoreConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	// Get version
	var version models.ConfigVersion
	if err := s.db.First(&version, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	// TODO: Implement actual configuration restore to FRR
	// This would involve applying the configuration to FRR via gRPC
	s.logger.Info("Configuration restore requested",
		zap.Uint("version_id", uint(id)),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration restore initiated",
		"version": version,
	})
}

// handleListAlerts handles listing all alerts
func (s *Server) handleListAlerts(c *gin.Context) {
	// Parse query parameters
	acknowledged := c.Query("acknowledged")
	severity := c.Query("severity")

	query := s.db.Preload("Peer").Preload("User").Order("created_at DESC")

	if acknowledged != "" {
		ack := acknowledged == "true"
		query = query.Where("acknowledged = ?", ack)
	}

	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	var alerts []models.Alert
	if err := query.Find(&alerts).Error; err != nil {
		s.logger.Error("Failed to list alerts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// handleAcknowledgeAlert handles acknowledging an alert
func (s *Server) handleAcknowledgeAlert(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	// Get current user ID
	userID, exists := authpkg.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get alert
	var alert models.Alert
	if err := s.db.First(&alert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	// Check if already acknowledged
	if alert.Acknowledged {
		c.JSON(http.StatusOK, gin.H{"message": "Alert already acknowledged"})
		return
	}

	// Acknowledge alert
	now := time.Now()
	alert.Acknowledged = true
	alert.AcknowledgedAt = &now
	alert.AcknowledgedBy = &userID

	if err := s.db.Save(&alert).Error; err != nil {
		s.logger.Error("Failed to acknowledge alert", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acknowledge alert"})
		return
	}

	s.logger.Info("Alert acknowledged",
		zap.Uint("alert_id", uint(id)),
		zap.Uint("user_id", userID),
	)

	c.JSON(http.StatusOK, alert)
}