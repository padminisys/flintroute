package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	FRR      FRRConfig      `mapstructure:"frr"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// FRRConfig represents FRR gRPC configuration
type FRRConfig struct {
	GRPCHost string `mapstructure:"grpc_host"`
	GRPCPort int    `mapstructure:"grpc_port"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	TokenExpiry   string `mapstructure:"token_expiry"`
	RefreshExpiry string `mapstructure:"refresh_expiry"`
}

// Load loads configuration from file or environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.path", "./data/flintroute.db")
	v.SetDefault("frr.grpc_host", "localhost")
	v.SetDefault("frr.grpc_port", 50051)
	v.SetDefault("auth.jwt_secret", "changeme-in-production")
	v.SetDefault("auth.token_expiry", "15m")
	v.SetDefault("auth.refresh_expiry", "168h") // 7 days

	// Set config file name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Enable environment variable override
	v.SetEnvPrefix("FLINTROUTE")
	v.AutomaticEnv()

	// Read config file if it exists
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; using defaults
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validate validates the configuration
func validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.FRR.GRPCPort < 1 || cfg.FRR.GRPCPort > 65535 {
		return fmt.Errorf("invalid FRR gRPC port: %d", cfg.FRR.GRPCPort)
	}

	if cfg.Auth.JWTSecret == "" || cfg.Auth.JWTSecret == "changeme-in-production" {
		fmt.Fprintf(os.Stderr, "WARNING: Using default JWT secret. Please set a secure secret in production!\n")
	}

	return nil
}