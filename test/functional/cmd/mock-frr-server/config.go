package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ServerConfig represents the mock FRR server configuration
type ServerConfig struct {
	Server     ServerSettings     `yaml:"server"`
	Simulation SimulationSettings `yaml:"simulation"`
	Logging    LoggingSettings    `yaml:"logging"`
}

// ServerSettings contains server connection settings
type ServerSettings struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// SimulationSettings contains behavior simulation settings
type SimulationSettings struct {
	SessionStateDelay time.Duration `yaml:"session_state_delay"`
	ErrorInjection    bool          `yaml:"error_injection"`
}

// LoggingSettings contains logging configuration
type LoggingSettings struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ServerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate validates the configuration
func (c *ServerConfig) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server host is required")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if c.Simulation.SessionStateDelay < 0 {
		return fmt.Errorf("session state delay must be non-negative")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.Logging.Level)
	}

	return nil
}

// GetAddress returns the server address in host:port format
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}