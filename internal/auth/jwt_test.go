package auth

import (
	"testing"
	"time"

	"github.com/padminisys/flintroute/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTManager(t *testing.T) {
	secretKey := "test-secret-key"
	tokenExpiry := 15 * time.Minute
	refreshExpiry := 7 * 24 * time.Hour

	manager := NewJWTManager(secretKey, tokenExpiry, refreshExpiry)

	assert.NotNil(t, manager)
	assert.Equal(t, secretKey, manager.secretKey)
	assert.Equal(t, tokenExpiry, manager.tokenExpiry)
	assert.Equal(t, refreshExpiry, manager.refreshExpiry)
}

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
	}

	t.Run("Generate valid token", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Token contains correct claims", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, user.Role, claims.Role)
	})
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     "user",
	}

	t.Run("Generate valid refresh token", func(t *testing.T) {
		token, expiresAt, err := manager.GenerateRefreshToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.True(t, expiresAt.After(time.Now()))
	})

	t.Run("Refresh token has correct expiry", func(t *testing.T) {
		_, expiresAt, err := manager.GenerateRefreshToken(user)
		assert.NoError(t, err)

		expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
		// Allow 1 second tolerance
		assert.WithinDuration(t, expectedExpiry, expiresAt, time.Second)
	})
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Role:     "admin",
	}

	t.Run("Validate valid token", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Username, claims.Username)
		assert.Equal(t, user.Role, claims.Role)
	})

	t.Run("Reject invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.here"
		claims, err := manager.ValidateToken(invalidToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, ErrInvalidToken, err)
	})

	t.Run("Reject token with wrong secret", func(t *testing.T) {
		wrongManager := NewJWTManager("wrong-secret", 15*time.Minute, 7*24*time.Hour)
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := wrongManager.ValidateToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Reject expired token", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", 1*time.Millisecond, 7*24*time.Hour)
		token, err := shortManager.GenerateToken(user)
		assert.NoError(t, err)

		// Wait for token to expire
		time.Sleep(10 * time.Millisecond)

		claims, err := shortManager.ValidateToken(token)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, ErrExpiredToken, err)
	})

	t.Run("Validate token before NotBefore time", func(t *testing.T) {
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		// Token should be valid immediately
		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
	})
}

func TestTokenClaims(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	t.Run("Claims contain all required fields", func(t *testing.T) {
		user := &models.User{
			ID:       42,
			Username: "claimsuser",
			Role:     "user",
		}

		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, uint(42), claims.UserID)
		assert.Equal(t, "claimsuser", claims.Username)
		assert.Equal(t, "user", claims.Role)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.NotBefore)
	})

	t.Run("Different users generate different tokens", func(t *testing.T) {
		user1 := &models.User{ID: 1, Username: "user1", Role: "admin"}
		user2 := &models.User{ID: 2, Username: "user2", Role: "user"}

		token1, err := manager.GenerateToken(user1)
		assert.NoError(t, err)

		token2, err := manager.GenerateToken(user2)
		assert.NoError(t, err)

		assert.NotEqual(t, token1, token2)

		claims1, _ := manager.ValidateToken(token1)
		claims2, _ := manager.ValidateToken(token2)

		assert.NotEqual(t, claims1.UserID, claims2.UserID)
		assert.NotEqual(t, claims1.Username, claims2.Username)
	})
}

func TestTokenExpiry(t *testing.T) {
	t.Run("Token expiry is set correctly", func(t *testing.T) {
		expiry := 30 * time.Minute
		manager := NewJWTManager("test-secret", expiry, 7*24*time.Hour)

		user := &models.User{ID: 1, Username: "testuser", Role: "user"}
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		claims, err := manager.ValidateToken(token)
		assert.NoError(t, err)

		expectedExpiry := time.Now().Add(expiry)
		actualExpiry := claims.ExpiresAt.Time

		// Allow 1 second tolerance
		assert.WithinDuration(t, expectedExpiry, actualExpiry, time.Second)
	})
}

func TestRefreshTokenExpiry(t *testing.T) {
	t.Run("Refresh token expiry is set correctly", func(t *testing.T) {
		refreshExpiry := 14 * 24 * time.Hour
		manager := NewJWTManager("test-secret", 15*time.Minute, refreshExpiry)

		user := &models.User{ID: 1, Username: "testuser", Role: "user"}
		token, expiresAt, err := manager.GenerateRefreshToken(user)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		expectedExpiry := time.Now().Add(refreshExpiry)
		// Allow 1 second tolerance
		assert.WithinDuration(t, expectedExpiry, expiresAt, time.Second)
	})
}