package client

import (
	"fmt"
	"sync"
	"time"
)

// TokenManager manages authentication tokens with automatic refresh
type TokenManager struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
	mu           sync.RWMutex
	client       *APIClient
}

// NewTokenManager creates a new token manager
func NewTokenManager(client *APIClient) *TokenManager {
	return &TokenManager{
		client: client,
	}
}

// SetTokens sets the access and refresh tokens
func (tm *TokenManager) SetTokens(accessToken, refreshToken string, expiresIn int64) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	tm.accessToken = accessToken
	tm.refreshToken = refreshToken
	tm.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// GetAccessToken returns the current access token, refreshing if necessary
func (tm *TokenManager) GetAccessToken() (string, error) {
	tm.mu.RLock()
	
	// Check if token is still valid (with 30 second buffer)
	if time.Now().Add(30 * time.Second).Before(tm.expiresAt) {
		token := tm.accessToken
		tm.mu.RUnlock()
		return token, nil
	}
	
	refreshToken := tm.refreshToken
	tm.mu.RUnlock()
	
	// Token is expired or about to expire, refresh it
	if refreshToken == "" {
		return "", fmt.Errorf("no refresh token available")
	}
	
	// Refresh the token
	response, err := tm.client.RefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	
	// Update tokens
	tm.SetTokens(response.AccessToken, response.RefreshToken, response.ExpiresIn)
	
	return response.AccessToken, nil
}

// GetRefreshToken returns the current refresh token
func (tm *TokenManager) GetRefreshToken() string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.refreshToken
}

// Clear clears all tokens
func (tm *TokenManager) Clear() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	tm.accessToken = ""
	tm.refreshToken = ""
	tm.expiresAt = time.Time{}
}

// IsAuthenticated returns true if we have valid tokens
func (tm *TokenManager) IsAuthenticated() bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	return tm.accessToken != "" && time.Now().Before(tm.expiresAt)
}

// GetAuthorizationHeader returns the Authorization header value
func (tm *TokenManager) GetAuthorizationHeader() (string, error) {
	token, err := tm.GetAccessToken()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Bearer %s", token), nil
}