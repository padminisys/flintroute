package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	authpkg "github.com/padminisys/flintroute/internal/auth"
	"github.com/padminisys/flintroute/internal/bgp"
	"github.com/padminisys/flintroute/internal/config"
	"github.com/padminisys/flintroute/internal/database"
	"github.com/padminisys/flintroute/internal/frr"
	"github.com/padminisys/flintroute/internal/websocket"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	db         *database.DB
	wsHub      *websocket.Hub
	bgpService *bgp.Service
	jwtManager *authpkg.JWTManager
	logger     *zap.Logger
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.Config, db *database.DB, wsHub *websocket.Hub, logger *zap.Logger) *Server {
	// Parse token expiry durations
	tokenExpiry, err := time.ParseDuration(cfg.Auth.TokenExpiry)
	if err != nil {
		tokenExpiry = 15 * time.Minute
	}

	refreshExpiry, err := time.ParseDuration(cfg.Auth.RefreshExpiry)
	if err != nil {
		refreshExpiry = 168 * time.Hour // 7 days
	}

	// Create JWT manager
	jwtManager := authpkg.NewJWTManager(cfg.Auth.JWTSecret, tokenExpiry, refreshExpiry)

	// Create FRR client
	frrClient, err := frr.NewClient(cfg.FRR.GRPCHost, cfg.FRR.GRPCPort, logger)
	if err != nil {
		logger.Error("Failed to create FRR client", zap.Error(err))
	}

	// Create BGP service
	bgpService := bgp.NewService(db, frrClient, wsHub, logger)

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(loggingMiddleware(logger))

	server := &Server{
		router:     router,
		config:     cfg,
		db:         db,
		wsHub:      wsHub,
		bgpService: bgpService,
		jwtManager: jwtManager,
		logger:     logger,
	}

	// Setup routes
	server.setupRoutes()

	// Start BGP monitoring
	go bgpService.StartMonitoring(context.Background(), 30*time.Second)

	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.handleHealth)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", s.handleLogin)
			auth.POST("/refresh", s.handleRefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(authpkg.AuthMiddleware(s.jwtManager))
		{
			// Auth
			protected.POST("/auth/logout", s.handleLogout)

			// BGP Peers
			peers := protected.Group("/bgp/peers")
			{
				peers.GET("", s.handleListPeers)
				peers.POST("", s.handleCreatePeer)
				peers.GET("/:id", s.handleGetPeer)
				peers.PUT("/:id", s.handleUpdatePeer)
				peers.DELETE("/:id", s.handleDeletePeer)
			}

			// BGP Sessions
			sessions := protected.Group("/bgp/sessions")
			{
				sessions.GET("", s.handleListSessions)
				sessions.GET("/:id", s.handleGetSession)
			}

			// Configuration
			configRoutes := protected.Group("/config")
			{
				configRoutes.GET("/versions", s.handleListConfigVersions)
				configRoutes.POST("/backup", s.handleBackupConfig)
				configRoutes.POST("/restore/:id", s.handleRestoreConfig)
			}

			// Alerts
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", s.handleListAlerts)
				alerts.POST("/:id/acknowledge", s.handleAcknowledgeAlert)
			}

			// WebSocket
			protected.GET("/ws", func(c *gin.Context) {
				s.wsHub.HandleWebSocket(c)
			})
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("Starting HTTP server", zap.String("address", addr))
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}