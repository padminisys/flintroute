package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/padminisys/flintroute/internal/database"
	"github.com/padminisys/flintroute/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *database.DB {
	t.Helper()

	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	logger := zap.NewNop()

	db, err := database.Initialize(dbPath, logger)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	return db
}

// SetupTestDBWithData creates a test database with sample data
func SetupTestDBWithData(t *testing.T) (*database.DB, *models.User, *models.BGPPeer) {
	t.Helper()

	db := SetupTestDB(t)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	user := &models.User{
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
		Role:         "admin",
		Active:       true,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Create test BGP peer
	peer := &models.BGPPeer{
		Name:        "Test Peer",
		IPAddress:   "192.168.1.1",
		ASN:         65001,
		RemoteASN:   65002,
		Description: "Test peer description",
		Enabled:     true,
		Multihop:    1,
	}
	if err := db.Create(peer).Error; err != nil {
		t.Fatalf("Failed to create test peer: %v", err)
	}

	return db, user, peer
}

// CleanupTestDB closes the database connection
func CleanupTestDB(t *testing.T, db *database.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Errorf("Failed to close test database: %v", err)
	}
}

// CreateTestLogger creates a no-op logger for testing
func CreateTestLogger() *zap.Logger {
	return zap.NewNop()
}

// CreateTestConfig creates a temporary config file for testing
func CreateTestConfig(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test config: %v", err)
	}

	return configPath
}

// SetupInMemoryDB creates a pure in-memory database (no file)
func SetupInMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.BGPPeer{},
		&models.BGPSession{},
		&models.ConfigVersion{},
		&models.Alert{},
		&models.RefreshToken{},
	); err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}