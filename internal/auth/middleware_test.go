package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/padminisys/flintroute/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthMiddleware(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	router := setupTestRouter()

	// Protected endpoint
	router.GET("/protected", AuthMiddleware(manager), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	t.Run("Allow valid token", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "testuser", Role: "admin"}
		token, err := manager.GenerateToken(user)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Reject missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("Reject invalid authorization format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("Reject token without Bearer prefix", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "testuser", Role: "admin"}
		token, _ := manager.GenerateToken(user)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Reject invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})

	t.Run("Reject expired token", func(t *testing.T) {
		shortManager := NewJWTManager("test-secret", 1*time.Millisecond, 7*24*time.Hour)
		user := &models.User{ID: 1, Username: "testuser", Role: "admin"}
		token, _ := shortManager.GenerateToken(user)

		time.Sleep(10 * time.Millisecond)

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Set user context on valid token", func(t *testing.T) {
		user := &models.User{ID: 42, Username: "contextuser", Role: "user"}
		token, _ := manager.GenerateToken(user)

		router := setupTestRouter()
		router.GET("/check-context", AuthMiddleware(manager), func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, uint(42), userID)

			username, exists := c.Get("username")
			assert.True(t, exists)
			assert.Equal(t, "contextuser", username)

			role, exists := c.Get("role")
			assert.True(t, exists)
			assert.Equal(t, "user", role)

			c.JSON(http.StatusOK, gin.H{"message": "context set"})
		})

		req := httptest.NewRequest(http.MethodGet, "/check-context", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAdminMiddleware(t *testing.T) {
	router := setupTestRouter()

	// Admin-only endpoint
	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", c.Query("role"))
		c.Next()
	}, AdminMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin access"})
	})

	t.Run("Allow admin role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin?role=admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Reject non-admin role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin?role=user", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Admin access required")
	})

	t.Run("Reject missing role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGetUserID(t *testing.T) {
	t.Run("Get existing user ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", uint(123))

		userID, exists := GetUserID(c)
		assert.True(t, exists)
		assert.Equal(t, uint(123), userID)
	})

	t.Run("Get non-existing user ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		userID, exists := GetUserID(c)
		assert.False(t, exists)
		assert.Equal(t, uint(0), userID)
	})

	t.Run("Get invalid type user ID", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("user_id", "not-a-uint")

		userID, exists := GetUserID(c)
		assert.False(t, exists)
		assert.Equal(t, uint(0), userID)
	})
}

func TestGetUsername(t *testing.T) {
	t.Run("Get existing username", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("username", "testuser")

		username, exists := GetUsername(c)
		assert.True(t, exists)
		assert.Equal(t, "testuser", username)
	})

	t.Run("Get non-existing username", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		username, exists := GetUsername(c)
		assert.False(t, exists)
		assert.Equal(t, "", username)
	})

	t.Run("Get invalid type username", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("username", 123)

		username, exists := GetUsername(c)
		assert.False(t, exists)
		assert.Equal(t, "", username)
	})
}

func TestGetRole(t *testing.T) {
	t.Run("Get existing role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("role", "admin")

		role, exists := GetRole(c)
		assert.True(t, exists)
		assert.Equal(t, "admin", role)
	})

	t.Run("Get non-existing role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())

		role, exists := GetRole(c)
		assert.False(t, exists)
		assert.Equal(t, "", role)
	})

	t.Run("Get invalid type role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("role", 456)

		role, exists := GetRole(c)
		assert.False(t, exists)
		assert.Equal(t, "", role)
	})
}

func TestMiddlewareChaining(t *testing.T) {
	manager := NewJWTManager("test-secret", 15*time.Minute, 7*24*time.Hour)
	router := setupTestRouter()

	// Chain auth and admin middleware
	router.GET("/admin-protected", AuthMiddleware(manager), AdminMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "admin success"})
	})

	t.Run("Allow admin with valid token", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "admin", Role: "admin"}
		token, _ := manager.GenerateToken(user)

		req := httptest.NewRequest(http.MethodGet, "/admin-protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Reject non-admin with valid token", func(t *testing.T) {
		user := &models.User{ID: 2, Username: "user", Role: "user"}
		token, _ := manager.GenerateToken(user)

		req := httptest.NewRequest(http.MethodGet, "/admin-protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Reject admin without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin-protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}