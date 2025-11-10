package authentication_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/flintroute/test/functional/pkg/client"
	"github.com/yourusername/flintroute/test/functional/pkg/testutil"
)

// TestLogin tests the login endpoint functionality
func TestLogin(t *testing.T) {
	// Setup logger
	logger, err := testutil.NewTestLogger("../../logs/test.log", "debug")
	require.NoError(t, err, "Failed to create test logger")
	defer logger.Close()

	logger.LogTestStart("TestLogin")
	startTime := time.Now()

	// Create API client
	apiClient := client.NewAPIClient("http://localhost:8080", logger.GetZapLogger())

	// Load fixture
	fixtureLoader := testutil.NewFixtureLoader("../../fixtures", logger.GetZapLogger())
	adminUser, err := fixtureLoader.LoadUser("admin_user")
	require.NoError(t, err, "Failed to load admin user fixture")

	// Test: Successful login
	t.Run("successful_login", func(t *testing.T) {
		logger.Info("Testing successful login")

		resp, err := apiClient.Login(adminUser.Username, adminUser.Password)
		require.NoError(t, err, "Login should succeed with valid credentials")

		// Verify response structure
		assert.NotEmpty(t, resp.AccessToken, "Access token should not be empty")
		assert.NotEmpty(t, resp.RefreshToken, "Refresh token should not be empty")
		assert.Greater(t, resp.ExpiresIn, int64(0), "ExpiresIn should be positive")

		// Verify user information
		assert.Equal(t, adminUser.Username, resp.User.Username, "Username should match")
		assert.Equal(t, adminUser.Email, resp.User.Email, "Email should match")
		assert.Equal(t, adminUser.Role, resp.User.Role, "Role should match")

		logger.Info("Successful login test passed")
	})

	// Test: Invalid credentials
	t.Run("invalid_credentials", func(t *testing.T) {
		logger.Info("Testing invalid credentials")

		_, err := apiClient.Login("invalid_user", "wrong_password")
		assert.Error(t, err, "Login should fail with invalid credentials")

		logger.Info("Invalid credentials test passed")
	})

	// Test: Empty username
	t.Run("empty_username", func(t *testing.T) {
		logger.Info("Testing empty username")

		_, err := apiClient.Login("", adminUser.Password)
		assert.Error(t, err, "Login should fail with empty username")

		logger.Info("Empty username test passed")
	})

	// Test: Empty password
	t.Run("empty_password", func(t *testing.T) {
		logger.Info("Testing empty password")

		_, err := apiClient.Login(adminUser.Username, "")
		assert.Error(t, err, "Login should fail with empty password")

		logger.Info("Empty password test passed")
	})

	// Test: Token authentication
	t.Run("authenticated_request", func(t *testing.T) {
		logger.Info("Testing authenticated request")

		// First login
		_, err := apiClient.Login(adminUser.Username, adminUser.Password)
		require.NoError(t, err, "Login should succeed")

		// Verify client is authenticated
		assert.True(t, apiClient.IsAuthenticated(), "Client should be authenticated after login")

		// Make an authenticated request (health check doesn't require auth, but we can test the mechanism)
		err = apiClient.HealthCheck()
		assert.NoError(t, err, "Health check should succeed")

		logger.Info("Authenticated request test passed")
	})

	// Test: Logout
	t.Run("logout", func(t *testing.T) {
		logger.Info("Testing logout")

		// First login
		_, err := apiClient.Login(adminUser.Username, adminUser.Password)
		require.NoError(t, err, "Login should succeed")

		// Logout
		err = apiClient.Logout()
		assert.NoError(t, err, "Logout should succeed")

		// Verify client is no longer authenticated
		assert.False(t, apiClient.IsAuthenticated(), "Client should not be authenticated after logout")

		logger.Info("Logout test passed")
	})

	duration := time.Since(startTime)
	logger.LogTestEnd("TestLogin", !t.Failed(), duration)
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	// Setup logger
	logger, err := testutil.NewTestLogger("../../logs/test.log", "debug")
	require.NoError(t, err, "Failed to create test logger")
	defer logger.Close()

	logger.LogTestStart("TestHealthCheck")
	startTime := time.Now()

	// Create API client
	apiClient := client.NewAPIClient("http://localhost:8080", logger.GetZapLogger())

	// Test: Health check without authentication
	t.Run("health_check_no_auth", func(t *testing.T) {
		logger.Info("Testing health check without authentication")

		err := apiClient.HealthCheck()
		assert.NoError(t, err, "Health check should succeed without authentication")

		logger.Info("Health check test passed")
	})

	duration := time.Since(startTime)
	logger.LogTestEnd("TestHealthCheck", !t.Failed(), duration)
}

// TestTokenRefresh tests the token refresh functionality
func TestTokenRefresh(t *testing.T) {
	// Setup logger
	logger, err := testutil.NewTestLogger("../../logs/test.log", "debug")
	require.NoError(t, err, "Failed to create test logger")
	defer logger.Close()

	logger.LogTestStart("TestTokenRefresh")
	startTime := time.Now()

	// Create API client
	apiClient := client.NewAPIClient("http://localhost:8080", logger.GetZapLogger())

	// Load fixture
	fixtureLoader := testutil.NewFixtureLoader("../../fixtures", logger.GetZapLogger())
	adminUser, err := fixtureLoader.LoadUser("admin_user")
	require.NoError(t, err, "Failed to load admin user fixture")

	// Test: Token refresh
	t.Run("token_refresh", func(t *testing.T) {
		logger.Info("Testing token refresh")

		// First login
		loginResp, err := apiClient.Login(adminUser.Username, adminUser.Password)
		require.NoError(t, err, "Login should succeed")

		// Get the refresh token
		refreshToken := loginResp.RefreshToken
		assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")

		// Refresh the token
		tokenResp, err := apiClient.RefreshToken(refreshToken)
		require.NoError(t, err, "Token refresh should succeed")

		// Verify new tokens
		assert.NotEmpty(t, tokenResp.AccessToken, "New access token should not be empty")
		assert.NotEmpty(t, tokenResp.RefreshToken, "New refresh token should not be empty")
		assert.Greater(t, tokenResp.ExpiresIn, int64(0), "ExpiresIn should be positive")

		// Verify tokens are different (optional, depends on implementation)
		// Note: Some implementations may return the same refresh token
		logger.Info("Token refresh test passed")
	})

	// Test: Invalid refresh token
	t.Run("invalid_refresh_token", func(t *testing.T) {
		logger.Info("Testing invalid refresh token")

		_, err := apiClient.RefreshToken("invalid_token")
		assert.Error(t, err, "Token refresh should fail with invalid token")

		logger.Info("Invalid refresh token test passed")
	})

	duration := time.Since(startTime)
	logger.LogTestEnd("TestTokenRefresh", !t.Failed(), duration)
}