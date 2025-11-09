package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("Load with default values", func(t *testing.T) {
		// Create temp directory without config file
		tmpDir := t.TempDir()
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// Check default values
		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "./data/flintroute.db", cfg.Database.Path)
		assert.Equal(t, "localhost", cfg.FRR.GRPCHost)
		assert.Equal(t, 50051, cfg.FRR.GRPCPort)
		assert.Equal(t, "changeme-in-production", cfg.Auth.JWTSecret)
		assert.Equal(t, "15m", cfg.Auth.TokenExpiry)
		assert.Equal(t, "168h", cfg.Auth.RefreshExpiry)
	})

	t.Run("Load from config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
server:
  host: 127.0.0.1
  port: 9090
database:
  path: /tmp/test.db
frr:
  grpc_host: frr-server
  grpc_port: 50052
auth:
  jwt_secret: my-secret-key
  token_expiry: 30m
  refresh_expiry: 336h
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "127.0.0.1", cfg.Server.Host)
		assert.Equal(t, 9090, cfg.Server.Port)
		assert.Equal(t, "/tmp/test.db", cfg.Database.Path)
		assert.Equal(t, "frr-server", cfg.FRR.GRPCHost)
		assert.Equal(t, 50052, cfg.FRR.GRPCPort)
		assert.Equal(t, "my-secret-key", cfg.Auth.JWTSecret)
		assert.Equal(t, "30m", cfg.Auth.TokenExpiry)
		assert.Equal(t, "336h", cfg.Auth.RefreshExpiry)
	})

	t.Run("Load with environment variables", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		// Set environment variables
		os.Setenv("FLINTROUTE_SERVER_PORT", "7070")
		os.Setenv("FLINTROUTE_AUTH_JWT_SECRET", "env-secret")
		defer func() {
			os.Unsetenv("FLINTROUTE_SERVER_PORT")
			os.Unsetenv("FLINTROUTE_AUTH_JWT_SECRET")
		}()

		cfg, err := Load()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, 7070, cfg.Server.Port)
		assert.Equal(t, "env-secret", cfg.Auth.JWTSecret)
	})

	t.Run("Invalid YAML file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		invalidContent := `
server:
  host: 127.0.0.1
  port: invalid_port
  nested:
    - item1
    - item2
  invalid_yaml: [
`
		err := os.WriteFile(configPath, []byte(invalidContent), 0644)
		assert.NoError(t, err)

		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		cfg, err := Load()
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestValidate(t *testing.T) {
	t.Run("Valid configuration", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
			FRR: FRRConfig{
				GRPCHost: "localhost",
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret: "secure-secret",
			},
		}

		err := validate(cfg)
		assert.NoError(t, err)
	})

	t.Run("Invalid server port - too low", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 0,
			},
			FRR: FRRConfig{
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret: "secret",
			},
		}

		err := validate(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server port")
	})

	t.Run("Invalid server port - too high", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 70000,
			},
			FRR: FRRConfig{
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret: "secret",
			},
		}

		err := validate(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid server port")
	})

	t.Run("Invalid FRR gRPC port - too low", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 8080,
			},
			FRR: FRRConfig{
				GRPCPort: -1,
			},
			Auth: AuthConfig{
				JWTSecret: "secret",
			},
		}

		err := validate(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid FRR gRPC port")
	})

	t.Run("Invalid FRR gRPC port - too high", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 8080,
			},
			FRR: FRRConfig{
				GRPCPort: 100000,
			},
			Auth: AuthConfig{
				JWTSecret: "secret",
			},
		}

		err := validate(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid FRR gRPC port")
	})

	t.Run("Warning for default JWT secret", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 8080,
			},
			FRR: FRRConfig{
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret: "changeme-in-production",
			},
		}

		// Capture stderr to check for warning
		// Note: This test validates that the function doesn't error,
		// but the warning is printed to stderr
		err := validate(cfg)
		assert.NoError(t, err)
	})

	t.Run("Empty JWT secret", func(t *testing.T) {
		cfg := &Config{
			Server: ServerConfig{
				Port: 8080,
			},
			FRR: FRRConfig{
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret: "",
			},
		}

		err := validate(cfg)
		assert.NoError(t, err) // Empty secret triggers warning but doesn't error
	})
}

func TestConfigStructures(t *testing.T) {
	t.Run("ServerConfig structure", func(t *testing.T) {
		server := ServerConfig{
			Host: "localhost",
			Port: 8080,
		}

		assert.Equal(t, "localhost", server.Host)
		assert.Equal(t, 8080, server.Port)
	})

	t.Run("DatabaseConfig structure", func(t *testing.T) {
		db := DatabaseConfig{
			Path: "/path/to/db",
		}

		assert.Equal(t, "/path/to/db", db.Path)
	})

	t.Run("FRRConfig structure", func(t *testing.T) {
		frr := FRRConfig{
			GRPCHost: "frr-host",
			GRPCPort: 50051,
		}

		assert.Equal(t, "frr-host", frr.GRPCHost)
		assert.Equal(t, 50051, frr.GRPCPort)
	})

	t.Run("AuthConfig structure", func(t *testing.T) {
		auth := AuthConfig{
			JWTSecret:     "secret",
			TokenExpiry:   "15m",
			RefreshExpiry: "168h",
		}

		assert.Equal(t, "secret", auth.JWTSecret)
		assert.Equal(t, "15m", auth.TokenExpiry)
		assert.Equal(t, "168h", auth.RefreshExpiry)
	})

	t.Run("Complete Config structure", func(t *testing.T) {
		cfg := Config{
			Server: ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
			},
			Database: DatabaseConfig{
				Path: "./data/db",
			},
			FRR: FRRConfig{
				GRPCHost: "localhost",
				GRPCPort: 50051,
			},
			Auth: AuthConfig{
				JWTSecret:     "secret",
				TokenExpiry:   "15m",
				RefreshExpiry: "168h",
			},
		}

		assert.Equal(t, "0.0.0.0", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "./data/db", cfg.Database.Path)
		assert.Equal(t, "localhost", cfg.FRR.GRPCHost)
		assert.Equal(t, 50051, cfg.FRR.GRPCPort)
		assert.Equal(t, "secret", cfg.Auth.JWTSecret)
	})
}

func TestConfigPrecedence(t *testing.T) {
	t.Run("Environment variables override config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
server:
  port: 8080
auth:
  jwt_secret: file-secret
`
		err := os.WriteFile(configPath, []byte(configContent), 0644)
		assert.NoError(t, err)

		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		// Set environment variable
		os.Setenv("FLINTROUTE_SERVER_PORT", "9090")
		os.Setenv("FLINTROUTE_AUTH_JWT_SECRET", "env-secret")
		defer func() {
			os.Unsetenv("FLINTROUTE_SERVER_PORT")
			os.Unsetenv("FLINTROUTE_AUTH_JWT_SECRET")
		}()

		cfg, err := Load()
		assert.NoError(t, err)

		// Environment variable should override file
		assert.Equal(t, 9090, cfg.Server.Port)
		assert.Equal(t, "env-secret", cfg.Auth.JWTSecret)
	})
}