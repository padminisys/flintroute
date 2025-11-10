package testutil

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseManager manages test database operations
type DatabaseManager struct {
	dbPath string
	db     *gorm.DB
	logger *zap.Logger
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(dbPath string, logger *zap.Logger) (*DatabaseManager, error) {
	dm := &DatabaseManager{
		dbPath: dbPath,
		logger: logger,
	}
	return dm, nil
}

// Initialize initializes the database connection and schema
func (dm *DatabaseManager) Initialize() error {
	// Configure GORM logger to be silent in tests
	gormLogger := logger.Default.LogMode(logger.Silent)

	// Open database connection
	db, err := gorm.Open(sqlite.Open(dm.dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	dm.db = db
	dm.logger.Info("Database connection established", zap.String("path", dm.dbPath))

	// Auto-migrate schema
	if err := dm.migrateSchema(); err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	return nil
}

// migrateSchema runs database migrations
func (dm *DatabaseManager) migrateSchema() error {
	// Define all models that need to be migrated
	models := []interface{}{
		&User{},
		&BGPPeer{},
		&BGPSession{},
		&ConfigVersion{},
		&Alert{},
		&RefreshToken{},
	}

	for _, model := range models {
		if err := dm.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model: %w", err)
		}
	}

	dm.logger.Info("Database schema migrated successfully")
	return nil
}

// Clean removes all data from tables but keeps the schema
func (dm *DatabaseManager) Clean() error {
	return dm.CleanTables()
}

// CleanTables removes all data from all tables
func (dm *DatabaseManager) CleanTables() error {
	tables := []string{
		"refresh_tokens",
		"alerts",
		"bgp_sessions",
		"bgp_peers",
		"config_versions",
		"users",
	}

	for _, table := range tables {
		if err := dm.db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	dm.logger.Info("All tables cleaned")
	return nil
}

// DropAllTables drops all tables from the database
func (dm *DatabaseManager) DropAllTables() error {
	tables := []string{
		"refresh_tokens",
		"alerts",
		"bgp_sessions",
		"bgp_peers",
		"config_versions",
		"users",
	}

	for _, table := range tables {
		if err := dm.db.Migrator().DropTable(table); err != nil {
			dm.logger.Warn("Failed to drop table", zap.String("table", table), zap.Error(err))
		}
	}

	dm.logger.Info("All tables dropped")
	return nil
}

// GetDB returns the underlying GORM database instance
func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// Close closes the database connection
func (dm *DatabaseManager) Close() error {
	if dm.db == nil {
		return nil
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	dm.logger.Info("Database connection closed")
	return nil
}

// VerifyPeerCount verifies the number of peers in the database
func (dm *DatabaseManager) VerifyPeerCount(expected int) error {
	var count int64
	if err := dm.db.Model(&BGPPeer{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count peers: %w", err)
	}

	if int(count) != expected {
		return fmt.Errorf("expected %d peers, got %d", expected, count)
	}

	dm.logger.Debug("Peer count verified", zap.Int("count", expected))
	return nil
}

// VerifySessionCount verifies the number of sessions in the database
func (dm *DatabaseManager) VerifySessionCount(expected int) error {
	var count int64
	if err := dm.db.Model(&BGPSession{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count sessions: %w", err)
	}

	if int(count) != expected {
		return fmt.Errorf("expected %d sessions, got %d", expected, count)
	}

	dm.logger.Debug("Session count verified", zap.Int("count", expected))
	return nil
}

// VerifyAlertCount verifies the number of alerts in the database
func (dm *DatabaseManager) VerifyAlertCount(expected int) error {
	var count int64
	if err := dm.db.Model(&Alert{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count alerts: %w", err)
	}

	if int(count) != expected {
		return fmt.Errorf("expected %d alerts, got %d", expected, count)
	}

	dm.logger.Debug("Alert count verified", zap.Int("count", expected))
	return nil
}

// CreateTestUser creates a test user in the database
func (dm *DatabaseManager) CreateTestUser(username, email, password, role string) (*User, error) {
	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: password, // In tests, this might be pre-hashed
		Role:         role,
		Active:       true,
	}

	if err := dm.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create test user: %w", err)
	}

	dm.logger.Debug("Test user created", zap.String("username", username))
	return user, nil
}

// GetPeerByIP retrieves a peer by IP address
func (dm *DatabaseManager) GetPeerByIP(ipAddress string) (*BGPPeer, error) {
	var peer BGPPeer
	if err := dm.db.Where("ip_address = ?", ipAddress).First(&peer).Error; err != nil {
		return nil, fmt.Errorf("failed to get peer by IP: %w", err)
	}
	return &peer, nil
}

// GetSessionByPeerID retrieves a session by peer ID
func (dm *DatabaseManager) GetSessionByPeerID(peerID uint) (*BGPSession, error) {
	var session BGPSession
	if err := dm.db.Where("peer_id = ?", peerID).First(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to get session by peer ID: %w", err)
	}
	return &session, nil
}

// CountUnacknowledgedAlerts counts unacknowledged alerts
func (dm *DatabaseManager) CountUnacknowledgedAlerts() (int64, error) {
	var count int64
	if err := dm.db.Model(&Alert{}).Where("acknowledged = ?", false).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count unacknowledged alerts: %w", err)
	}
	return count, nil
}