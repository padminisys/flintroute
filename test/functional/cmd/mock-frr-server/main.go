package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultConfigPath = "./test/functional/config/mock-frr-config.yaml"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", defaultConfigPath, "Path to configuration file")
	flag.Parse()

	// Load configuration
	config, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := initLogger(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Mock FRR server starting",
		zap.String("config", *configPath),
		zap.String("address", config.GetAddress()),
	)

	// Create mock FRR server
	server := NewMockFRRServer(config, logger)

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(); err != nil {
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		logger.Error("Server error", zap.Error(err))
		os.Exit(1)
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		server.Stop()
		logger.Info("Mock FRR server stopped")
	}
}

// initLogger initializes the zap logger based on configuration
func initLogger(config *ServerConfig) (*zap.Logger, error) {
	// Parse log level
	var level zapcore.Level
	switch config.Logging.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create console encoder for stdout
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// Create file encoder for log file
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// Ensure log directory exists
	if config.Logging.File != "" {
		logDir := filepath.Dir(config.Logging.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Open log file
		logFile, err := os.OpenFile(config.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		// Create multi-writer core (console + file)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level),
		)

		return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
	}

	// Console only
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}