package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/padminisys/flintroute/internal/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestInitialize(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Initialize database successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer db.Close()

		// Verify database file was created
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})

	t.Run("Create directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "nested", "dir", "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		defer db.Close()

		// Verify nested directory was created
		_, err = os.Stat(filepath.Dir(dbPath))
		assert.NoError(t, err)
	})

	t.Run("Auto-migrate all models", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		// Verify tables exist by attempting to query them
		var count int64

		err = db.Model(&models.User{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&models.BGPPeer{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&models.BGPSession{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&models.ConfigVersion{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&models.Alert{}).Count(&count).Error
		assert.NoError(t, err)

		err = db.Model(&models.RefreshToken{}).Count(&count).Error
		assert.NoError(t, err)
	})

	t.Run("Create default admin user", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		// Verify default admin user exists
		var user models.User
		err = db.Where("username = ?", "admin").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "admin", user.Username)
		assert.Equal(t, "admin@flintroute.local", user.Email)
		assert.Equal(t, "admin", user.Role)
		assert.True(t, user.Active)

		// Verify password is hashed correctly
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("admin"))
		assert.NoError(t, err)
	})

	t.Run("Do not create duplicate admin user", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		// Count admin users
		var count int64
		err = db.Model(&models.User{}).Where("username = ?", "admin").Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Close and reinitialize
		db.Close()

		db, err = Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		// Verify still only one admin user
		err = db.Model(&models.User{}).Where("username = ?", "admin").Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Invalid database path", func(t *testing.T) {
		// Try to create database in a read-only location (if possible)
		dbPath := "/root/readonly/test.db"

		db, err := Initialize(dbPath, logger)
		if err == nil {
			// If no error (running as root or path is writable), clean up
			db.Close()
			os.Remove(dbPath)
		} else {
			// Expected error for non-writable path
			assert.Error(t, err)
			assert.Nil(t, db)
		}
	})
}

func TestCreateDefaultUser(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Create default user when none exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		var count int64
		err = db.Model(&models.User{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Skip creation when users exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)

		// Create another user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := models.User{
			Username:     "testuser",
			PasswordHash: string(hashedPassword),
			Email:        "test@example.com",
			Role:         "user",
			Active:       true,
		}
		err = db.Create(&user).Error
		assert.NoError(t, err)

		db.Close()

		// Reinitialize
		db, err = Initialize(dbPath, logger)
		assert.NoError(t, err)
		defer db.Close()

		// Should have 2 users (admin + testuser)
		var count int64
		err = db.Model(&models.User{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestGetDB(t *testing.T) {
	logger := zap.NewNop()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Initialize(dbPath, logger)
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Get underlying GORM DB", func(t *testing.T) {
		gormDB := db.GetDB()
		assert.NotNil(t, gormDB)

		// Verify it's the same instance
		assert.Equal(t, db.DB, gormDB)
	})

	t.Run("Use returned DB for queries", func(t *testing.T) {
		gormDB := db.GetDB()

		var count int64
		err := gormDB.Model(&models.User{}).Count(&count).Error
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(1))
	})
}

func TestClose(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Close database successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)

		err = db.Close()
		assert.NoError(t, err)
	})

	t.Run("Operations fail after close", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		db, err := Initialize(dbPath, logger)
		assert.NoError(t, err)

		err = db.Close()
		assert.NoError(t, err)

		// Try to perform operation after close
		var count int64
		err = db.Model(&models.User{}).Count(&count).Error
		assert.Error(t, err)
	})
}

func TestDatabaseOperations(t *testing.T) {
	logger := zap.NewNop()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Initialize(dbPath, logger)
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Create and retrieve user", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		user := models.User{
			Username:     "createtest",
			PasswordHash: string(hashedPassword),
			Email:        "create@test.com",
			Role:         "user",
			Active:       true,
		}

		err := db.Create(&user).Error
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)

		var retrieved models.User
		err = db.First(&retrieved, user.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, user.Username, retrieved.Username)
	})

	t.Run("Create and retrieve BGP peer", func(t *testing.T) {
		peer := models.BGPPeer{
			Name:      "Test Peer",
			IPAddress: "192.168.1.1",
			ASN:       65001,
			RemoteASN: 65002,
			Enabled:   true,
		}

		err := db.Create(&peer).Error
		assert.NoError(t, err)
		assert.NotZero(t, peer.ID)

		var retrieved models.BGPPeer
		err = db.First(&retrieved, peer.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, peer.Name, retrieved.Name)
	})

	t.Run("Create and retrieve alert", func(t *testing.T) {
		alert := models.Alert{
			Type:     "test_alert",
			Severity: "info",
			Message:  "Test message",
		}

		err := db.Create(&alert).Error
		assert.NoError(t, err)
		assert.NotZero(t, alert.ID)

		var retrieved models.Alert
		err = db.First(&retrieved, alert.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, alert.Message, retrieved.Message)
	})
}

func TestDatabaseConcurrency(t *testing.T) {
	logger := zap.NewNop()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := Initialize(dbPath, logger)
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Concurrent writes", func(t *testing.T) {
		done := make(chan bool, 10)

		for i := 0; i < 10; i++ {
			go func(index int) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("pass"), bcrypt.DefaultCost)
				user := models.User{
					Username:     "concurrent" + string(rune(index)),
					PasswordHash: string(hashedPassword),
					Email:        "concurrent" + string(rune(index)) + "@test.com",
					Role:         "user",
					Active:       true,
				}
				db.Create(&user)
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify all users were created
		var count int64
		err := db.Model(&models.User{}).Count(&count).Error
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, count, int64(10))
	})
}