package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/padminisys/flintroute/internal/auth"
	"github.com/padminisys/flintroute/internal/database"
	"github.com/padminisys/flintroute/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestServer(t *testing.T) (*Server, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test.db"
	logger := zap.NewNop()

	dbWrapper, err := database.Initialize(dbPath, logger)
	assert.NoError(t, err)

	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)

	server := &Server{
		db:         dbWrapper,
		logger:     logger,
		jwtManager: jwtManager,
	}

	return server, dbWrapper.GetDB()
}

func TestHandleLogin(t *testing.T) {
	server, db := setupTestServer(t)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := models.User{
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
		Role:         "admin",
		Active:       true,
	}
	db.Create(&user)

	t.Run("Successful login", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", server.handleLogin)

		reqBody := LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, "testuser", response.User.Username)
		assert.Equal(t, "admin", response.User.Role)
	})

	t.Run("Invalid credentials - wrong password", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", server.handleLogin)

		reqBody := LoginRequest{
			Username: "testuser",
			Password: "wrongpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("Invalid credentials - user not found", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", server.handleLogin)

		reqBody := LoginRequest{
			Username: "nonexistent",
			Password: "testpass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Inactive user", func(t *testing.T) {
		// Create inactive user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
		inactiveUser := models.User{
			Username:     "inactive",
			PasswordHash: string(hashedPassword),
			Email:        "inactive@example.com",
			Role:         "user",
		}
		// Create user first, then update Active to false (workaround for GORM default value issue)
		db.Create(&inactiveUser)
		db.Model(&inactiveUser).Update("active", false)

		router := gin.New()
		router.POST("/login", server.handleLogin)

		reqBody := LoginRequest{
			Username: "inactive",
			Password: "pass",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Account is disabled")
	})

	t.Run("Invalid request body", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", server.handleLogin)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		router := gin.New()
		router.POST("/login", server.handleLogin)

		reqBody := map[string]string{
			"username": "testuser",
			// Missing password
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandleRefreshToken(t *testing.T) {
	server, db := setupTestServer(t)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := models.User{
		Username:     "refreshuser",
		PasswordHash: string(hashedPassword),
		Email:        "refresh@example.com",
		Role:         "user",
		Active:       true,
	}
	db.Create(&user)

	t.Run("Successful token refresh", func(t *testing.T) {
		// Generate refresh token
		refreshToken, expiresAt, _ := server.jwtManager.GenerateRefreshToken(&user)
		tokenModel := models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: expiresAt,
		}
		// Create token first, then update Revoked to false (workaround for GORM default value issue)
		db.Create(&tokenModel)
		db.Model(&tokenModel).Update("revoked", false)

		router := gin.New()
		router.POST("/refresh", server.handleRefreshToken)

		reqBody := RefreshRequest{
			RefreshToken: refreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, refreshToken, response.RefreshToken) // Should be a new token
	})

	t.Run("Invalid refresh token", func(t *testing.T) {
		router := gin.New()
		router.POST("/refresh", server.handleRefreshToken)

		reqBody := RefreshRequest{
			RefreshToken: "invalid.token.here",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Revoked refresh token", func(t *testing.T) {
		// Generate and revoke token
		refreshToken, expiresAt, _ := server.jwtManager.GenerateRefreshToken(&user)
		tokenModel := models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: expiresAt,
		}
		db.Create(&tokenModel)
		db.Model(&tokenModel).Update("revoked", true) // Mark as revoked

		router := gin.New()
		router.POST("/refresh", server.handleRefreshToken)

		reqBody := RefreshRequest{
			RefreshToken: refreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Expired refresh token", func(t *testing.T) {
		// Generate expired token
		refreshToken, _, _ := server.jwtManager.GenerateRefreshToken(&user)
		tokenModel := models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		}
		db.Create(&tokenModel)
		db.Model(&tokenModel).Update("revoked", false)

		router := gin.New()
		router.POST("/refresh", server.handleRefreshToken)

		reqBody := RefreshRequest{
			RefreshToken: refreshToken,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandleLogout(t *testing.T) {
	server, db := setupTestServer(t)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := models.User{
		Username:     "logoutuser",
		PasswordHash: string(hashedPassword),
		Email:        "logout@example.com",
		Role:         "user",
		Active:       true,
	}
	db.Create(&user)

	t.Run("Successful logout", func(t *testing.T) {
		// Create refresh tokens
		refreshToken1, expiresAt1, _ := server.jwtManager.GenerateRefreshToken(&user)
		refreshToken2, expiresAt2, _ := server.jwtManager.GenerateRefreshToken(&user)

		token1 := &models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken1,
			ExpiresAt: expiresAt1,
		}
		db.Create(token1)
		db.Model(token1).Update("revoked", false)
		
		token2 := &models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken2,
			ExpiresAt: expiresAt2,
		}
		db.Create(token2)
		db.Model(token2).Update("revoked", false)

		// Generate access token
		accessToken, _ := server.jwtManager.GenerateToken(&user)

		router := gin.New()
		router.POST("/logout", server.handleLogout)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify all refresh tokens are revoked
		var count int64
		db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = ?", user.ID, false).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Logout without authorization header", func(t *testing.T) {
		router := gin.New()
		router.POST("/logout", server.handleLogout)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Logout with invalid token", func(t *testing.T) {
		router := gin.New()
		router.POST("/logout", server.handleLogout)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer invalid.token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLoginResponse(t *testing.T) {
	t.Run("Create login response", func(t *testing.T) {
		response := LoginResponse{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			ExpiresIn:    3600,
			User: UserInfo{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
				Role:     "admin",
			},
		}

		assert.Equal(t, "access_token", response.AccessToken)
		assert.Equal(t, "refresh_token", response.RefreshToken)
		assert.Equal(t, int64(3600), response.ExpiresIn)
		assert.Equal(t, "testuser", response.User.Username)
	})
}

func TestUserInfo(t *testing.T) {
	t.Run("Create user info", func(t *testing.T) {
		userInfo := UserInfo{
			ID:       42,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		}

		assert.Equal(t, uint(42), userInfo.ID)
		assert.Equal(t, "testuser", userInfo.Username)
		assert.Equal(t, "test@example.com", userInfo.Email)
		assert.Equal(t, "user", userInfo.Role)
	})
}