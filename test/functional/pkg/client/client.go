package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// APIClient is a client for the FlintRoute REST API
type APIClient struct {
	baseURL      string
	httpClient   *http.Client
	tokenManager *TokenManager
	logger       *zap.Logger
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string, logger *zap.Logger) *APIClient {
	client := &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
	client.tokenManager = NewTokenManager(client)
	return client
}

// SetTimeout sets the HTTP client timeout
func (c *APIClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// doRequest performs an HTTP request with automatic authentication
func (c *APIClient) doRequest(method, path string, body interface{}, authenticated bool) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
		c.logger.Debug("Request body", zap.String("body", string(jsonData)))
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication if required
	if authenticated {
		authHeader, err := c.tokenManager.GetAuthorizationHeader()
		if err != nil {
			return nil, fmt.Errorf("failed to get authorization header: %w", err)
		}
		req.Header.Set("Authorization", authHeader)
	}

	c.logger.Debug("Making request",
		zap.String("method", method),
		zap.String("url", fullURL),
		zap.Bool("authenticated", authenticated),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	c.logger.Debug("Response received",
		zap.Int("status", resp.StatusCode),
		zap.String("status_text", resp.Status),
	)

	return resp, nil
}

// parseResponse parses the response body into the target struct
func (c *APIClient) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	c.logger.Debug("Response body", zap.String("body", string(body)))

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Error)
	}

	if target != nil {
		if err := json.Unmarshal(body, target); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// Login authenticates with the API
func (c *APIClient) Login(username, password string) (*LoginResponse, error) {
	req := LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.doRequest("POST", "/api/v1/auth/login", req, false)
	if err != nil {
		return nil, err
	}

	var loginResp LoginResponse
	if err := c.parseResponse(resp, &loginResp); err != nil {
		return nil, err
	}

	// Store tokens
	c.tokenManager.SetTokens(loginResp.AccessToken, loginResp.RefreshToken, loginResp.ExpiresIn)

	c.logger.Info("Logged in successfully", zap.String("username", username))

	return &loginResp, nil
}

// RefreshToken refreshes the access token
func (c *APIClient) RefreshToken(refreshToken string) (*TokenResponse, error) {
	req := RefreshRequest{
		RefreshToken: refreshToken,
	}

	resp, err := c.doRequest("POST", "/api/v1/auth/refresh", req, false)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := c.parseResponse(resp, &tokenResp); err != nil {
		return nil, err
	}

	c.logger.Info("Token refreshed successfully")

	return &tokenResp, nil
}

// Logout logs out from the API
func (c *APIClient) Logout() error {
	resp, err := c.doRequest("POST", "/api/v1/auth/logout", nil, true)
	if err != nil {
		return err
	}

	var msgResp MessageResponse
	if err := c.parseResponse(resp, &msgResp); err != nil {
		return err
	}

	// Clear tokens
	c.tokenManager.Clear()

	c.logger.Info("Logged out successfully")

	return nil
}

// CreatePeer creates a new BGP peer
func (c *APIClient) CreatePeer(peer *PeerRequest) (*Peer, error) {
	resp, err := c.doRequest("POST", "/api/v1/bgp/peers", peer, true)
	if err != nil {
		return nil, err
	}

	var createdPeer Peer
	if err := c.parseResponse(resp, &createdPeer); err != nil {
		return nil, err
	}

	c.logger.Info("Peer created", zap.Uint("id", createdPeer.ID), zap.String("name", createdPeer.Name))

	return &createdPeer, nil
}

// ListPeers lists all BGP peers
func (c *APIClient) ListPeers() ([]*Peer, error) {
	resp, err := c.doRequest("GET", "/api/v1/bgp/peers", nil, true)
	if err != nil {
		return nil, err
	}

	var peersResp PeersResponse
	if err := c.parseResponse(resp, &peersResp); err != nil {
		return nil, err
	}

	c.logger.Debug("Peers listed", zap.Int("count", len(peersResp.Peers)))

	return peersResp.Peers, nil
}

// GetPeer gets a specific BGP peer
func (c *APIClient) GetPeer(id uint) (*Peer, error) {
	path := fmt.Sprintf("/api/v1/bgp/peers/%d", id)
	resp, err := c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var peer Peer
	if err := c.parseResponse(resp, &peer); err != nil {
		return nil, err
	}

	c.logger.Debug("Peer retrieved", zap.Uint("id", id))

	return &peer, nil
}

// UpdatePeer updates a BGP peer
func (c *APIClient) UpdatePeer(id uint, updates *PeerRequest) (*Peer, error) {
	path := fmt.Sprintf("/api/v1/bgp/peers/%d", id)
	resp, err := c.doRequest("PUT", path, updates, true)
	if err != nil {
		return nil, err
	}

	var peer Peer
	if err := c.parseResponse(resp, &peer); err != nil {
		return nil, err
	}

	c.logger.Info("Peer updated", zap.Uint("id", id))

	return &peer, nil
}

// DeletePeer deletes a BGP peer
func (c *APIClient) DeletePeer(id uint) error {
	path := fmt.Sprintf("/api/v1/bgp/peers/%d", id)
	resp, err := c.doRequest("DELETE", path, nil, true)
	if err != nil {
		return err
	}

	var msgResp MessageResponse
	if err := c.parseResponse(resp, &msgResp); err != nil {
		return err
	}

	c.logger.Info("Peer deleted", zap.Uint("id", id))

	return nil
}

// ListSessions lists all BGP sessions
func (c *APIClient) ListSessions() ([]*Session, error) {
	resp, err := c.doRequest("GET", "/api/v1/bgp/sessions", nil, true)
	if err != nil {
		return nil, err
	}

	var sessionsResp SessionsResponse
	if err := c.parseResponse(resp, &sessionsResp); err != nil {
		return nil, err
	}

	c.logger.Debug("Sessions listed", zap.Int("count", len(sessionsResp.Sessions)))

	return sessionsResp.Sessions, nil
}

// GetSession gets a specific BGP session
func (c *APIClient) GetSession(id uint) (*Session, error) {
	path := fmt.Sprintf("/api/v1/bgp/sessions/%d", id)
	resp, err := c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := c.parseResponse(resp, &session); err != nil {
		return nil, err
	}

	c.logger.Debug("Session retrieved", zap.Uint("id", id))

	return &session, nil
}

// ListConfigVersions lists all configuration versions
func (c *APIClient) ListConfigVersions() ([]*ConfigVersion, error) {
	resp, err := c.doRequest("GET", "/api/v1/config/versions", nil, true)
	if err != nil {
		return nil, err
	}

	var versionsResp ConfigVersionsResponse
	if err := c.parseResponse(resp, &versionsResp); err != nil {
		return nil, err
	}

	c.logger.Debug("Config versions listed", zap.Int("count", len(versionsResp.Versions)))

	return versionsResp.Versions, nil
}

// BackupConfig creates a configuration backup
func (c *APIClient) BackupConfig(description string) (*ConfigVersion, error) {
	req := BackupConfigRequest{
		Description: description,
	}

	resp, err := c.doRequest("POST", "/api/v1/config/backup", req, true)
	if err != nil {
		return nil, err
	}

	var version ConfigVersion
	if err := c.parseResponse(resp, &version); err != nil {
		return nil, err
	}

	c.logger.Info("Config backed up", zap.Uint("version_id", version.ID))

	return &version, nil
}

// RestoreConfig restores a configuration version
func (c *APIClient) RestoreConfig(id uint) error {
	path := fmt.Sprintf("/api/v1/config/restore/%d", id)
	resp, err := c.doRequest("POST", path, nil, true)
	if err != nil {
		return err
	}

	var msgResp MessageResponse
	if err := c.parseResponse(resp, &msgResp); err != nil {
		return err
	}

	c.logger.Info("Config restore initiated", zap.Uint("version_id", id))

	return nil
}

// ListAlerts lists alerts with optional filters
func (c *APIClient) ListAlerts(params *AlertQueryParams) ([]*Alert, error) {
	path := "/api/v1/alerts"
	
	if params != nil {
		query := url.Values{}
		if params.Acknowledged != nil {
			if *params.Acknowledged {
				query.Set("acknowledged", "true")
			} else {
				query.Set("acknowledged", "false")
			}
		}
		if params.Severity != "" {
			query.Set("severity", params.Severity)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	resp, err := c.doRequest("GET", path, nil, true)
	if err != nil {
		return nil, err
	}

	var alertsResp AlertsResponse
	if err := c.parseResponse(resp, &alertsResp); err != nil {
		return nil, err
	}

	c.logger.Debug("Alerts listed", zap.Int("count", len(alertsResp.Alerts)))

	return alertsResp.Alerts, nil
}

// AcknowledgeAlert acknowledges an alert
func (c *APIClient) AcknowledgeAlert(id uint) error {
	path := fmt.Sprintf("/api/v1/alerts/%d/acknowledge", id)
	resp, err := c.doRequest("POST", path, nil, true)
	if err != nil {
		return err
	}

	var alert Alert
	if err := c.parseResponse(resp, &alert); err != nil {
		return err
	}

	c.logger.Info("Alert acknowledged", zap.Uint("id", id))

	return nil
}

// HealthCheck performs a health check
func (c *APIClient) HealthCheck() error {
	resp, err := c.doRequest("GET", "/health", nil, false)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	resp.Body.Close()
	c.logger.Debug("Health check passed")

	return nil
}

// IsAuthenticated returns true if the client is authenticated
func (c *APIClient) IsAuthenticated() bool {
	return c.tokenManager.IsAuthenticated()
}