package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/padminisys/flintroute/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB wraps the GORM database connection
type DB struct {
	*gorm.DB
	logger *zap.Logger
}

// Initialize creates and initializes the database
func Initialize(dbPath string, log *zap.Logger) (*DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Silent)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
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
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	database := &DB{
		DB:     db,
		logger: log,
	}

	// Create default admin user if no users exist
	if err := database.createDefaultUser(); err != nil {
		return nil, fmt.Errorf("failed to create default user: %w", err)
	}

	log.Info("Database initialized successfully", zap.String("path", dbPath))

	return database, nil
}

// createDefaultUser creates a default admin user if no users exist
func (db *DB) createDefaultUser() error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // Users already exist
	}

	// Hash default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create default admin user
	user := models.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		Email:        "admin@flintroute.local",
		Role:         "admin",
		Active:       true,
	}

	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create default user: %w", err)
	}

	db.logger.Info("Created default admin user",
		zap.String("username", "admin"),
		zap.String("password", "admin"),
	)
	db.logger.Warn("Please change the default admin password immediately!")

	return nil
}

// GetDB returns the underlying GORM DB instance
func (db *DB) GetDB() *gorm.DB {
	return db.DB
}

// Close closes the database connection
func (db *DB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}