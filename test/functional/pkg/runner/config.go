package runner

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// TestConfig represents the test configuration
type TestConfig struct {
	ServerURL        string        `yaml:"server_url"`
	DatabasePath     string        `yaml:"database_path"`
	MockFRRURL       string        `yaml:"mock_frr_url"`
	Timeout          time.Duration `yaml:"timeout"`
	CleanupOnSuccess bool          `yaml:"cleanup_on_success"`
	LogLevel         string        `yaml:"log_level"`
	Parallel         bool          `yaml:"parallel"`
	FixturesPath     string        `yaml:"fixtures_path"`
	ResultsPath      string        `yaml:"results_path"`
	LogsPath         string        `yaml:"logs_path"`
	MaxRetries       int           `yaml:"max_retries"`
	RetryDelay       time.Duration `yaml:"retry_delay"`
}

// DefaultConfig returns a default test configuration
func DefaultConfig() *TestConfig {
	return &TestConfig{
		ServerURL:        "http://localhost:8080",
		DatabasePath:     "./tmp/test.db",
		MockFRRURL:       "localhost:50051",
		Timeout:          30 * time.Second,
		CleanupOnSuccess: true,
		LogLevel:         "info",
		Parallel:         false,
		FixturesPath:     "./fixtures",
		ResultsPath:      "./results",
		LogsPath:         "./logs",
		MaxRetries:       3,
		RetryDelay:       1 * time.Second,
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*TestConfig, error) {
	// Start with default config
	config := DefaultConfig()

	// If no config path provided, return defaults
	if configPath == "" {
		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *TestConfig) Validate() error {
	if c.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}

	if c.DatabasePath == "" {
		return fmt.Errorf("database_path is required")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	if c.FixturesPath == "" {
		c.FixturesPath = "./fixtures"
	}

	if c.ResultsPath == "" {
		c.ResultsPath = "./results"
	}

	if c.LogsPath == "" {
		c.LogsPath = "./logs"
	}

	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}

	if c.RetryDelay <= 0 {
		c.RetryDelay = 1 * time.Second
	}

	return nil
}

// SaveConfig saves the configuration to a YAML file
func (c *TestConfig) SaveConfig(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}