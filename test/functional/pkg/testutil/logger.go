package testutil

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestLogger wraps zap.Logger with test-specific functionality
type TestLogger struct {
	logger *zap.Logger
	file   *os.File
}

// NewTestLogger creates a new test logger
func NewTestLogger(logPath string, level string) (*TestLogger, error) {
	// Parse log level
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create file writer
	var file *os.File
	var err error
	if logPath != "" {
		file, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}

	// Create core
	var core zapcore.Core
	if file != nil {
		// Log to both file and console
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			zapLevel,
		)
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
		core = zapcore.NewTee(fileCore, consoleCore)
	} else {
		// Log to console only
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			zapLevel,
		)
	}

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &TestLogger{
		logger: logger,
		file:   file,
	}, nil
}

// Debug logs a debug message
func (tl *TestLogger) Debug(msg string, fields ...zap.Field) {
	tl.logger.Debug(msg, fields...)
}

// Info logs an info message
func (tl *TestLogger) Info(msg string, fields ...zap.Field) {
	tl.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (tl *TestLogger) Warn(msg string, fields ...zap.Field) {
	tl.logger.Warn(msg, fields...)
}

// Error logs an error message
func (tl *TestLogger) Error(msg string, fields ...zap.Field) {
	tl.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (tl *TestLogger) Fatal(msg string, fields ...zap.Field) {
	tl.logger.Fatal(msg, fields...)
}

// LogRequest logs an HTTP request
func (tl *TestLogger) LogRequest(method, url string, body interface{}) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("url", url),
	}
	if body != nil {
		fields = append(fields, zap.Any("body", body))
	}
	tl.logger.Info("HTTP Request", fields...)
}

// LogResponse logs an HTTP response
func (tl *TestLogger) LogResponse(statusCode int, body interface{}, duration time.Duration) {
	fields := []zap.Field{
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
	}
	if body != nil {
		fields = append(fields, zap.Any("body", body))
	}
	tl.logger.Info("HTTP Response", fields...)
}

// LogTestStart logs the start of a test
func (tl *TestLogger) LogTestStart(testName string) {
	tl.logger.Info("Test started",
		zap.String("test", testName),
		zap.Time("start_time", time.Now()),
	)
}

// LogTestEnd logs the end of a test
func (tl *TestLogger) LogTestEnd(testName string, passed bool, duration time.Duration) {
	status := "PASSED"
	if !passed {
		status = "FAILED"
	}
	tl.logger.Info("Test completed",
		zap.String("test", testName),
		zap.String("status", status),
		zap.Duration("duration", duration),
	)
}

// LogTestSkipped logs a skipped test
func (tl *TestLogger) LogTestSkipped(testName, reason string) {
	tl.logger.Info("Test skipped",
		zap.String("test", testName),
		zap.String("reason", reason),
	)
}

// LogSetup logs setup operations
func (tl *TestLogger) LogSetup(operation string) {
	tl.logger.Info("Setup operation", zap.String("operation", operation))
}

// LogTeardown logs teardown operations
func (tl *TestLogger) LogTeardown(operation string) {
	tl.logger.Info("Teardown operation", zap.String("operation", operation))
}

// LogAssertion logs an assertion
func (tl *TestLogger) LogAssertion(description string, passed bool) {
	if passed {
		tl.logger.Debug("Assertion passed", zap.String("assertion", description))
	} else {
		tl.logger.Error("Assertion failed", zap.String("assertion", description))
	}
}

// LogDatabaseOperation logs a database operation
func (tl *TestLogger) LogDatabaseOperation(operation, table string, count int) {
	tl.logger.Debug("Database operation",
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Int("count", count),
	)
}

// LogFixtureLoad logs fixture loading
func (tl *TestLogger) LogFixtureLoad(fixtureType, name string) {
	tl.logger.Debug("Fixture loaded",
		zap.String("type", fixtureType),
		zap.String("name", name),
	)
}

// Sync flushes any buffered log entries
func (tl *TestLogger) Sync() error {
	return tl.logger.Sync()
}

// Close closes the logger and any associated files
func (tl *TestLogger) Close() error {
	// Sync before closing
	if err := tl.logger.Sync(); err != nil {
		// Ignore sync errors on stdout/stderr
		if tl.file != nil {
			return err
		}
	}

	// Close file if open
	if tl.file != nil {
		if err := tl.file.Close(); err != nil {
			return fmt.Errorf("failed to close log file: %w", err)
		}
	}

	return nil
}

// GetZapLogger returns the underlying zap.Logger
func (tl *TestLogger) GetZapLogger() *zap.Logger {
	return tl.logger
}

// With creates a child logger with additional fields
func (tl *TestLogger) With(fields ...zap.Field) *TestLogger {
	return &TestLogger{
		logger: tl.logger.With(fields...),
		file:   tl.file,
	}
}

// Named creates a named logger
func (tl *TestLogger) Named(name string) *TestLogger {
	return &TestLogger{
		logger: tl.logger.Named(name),
		file:   tl.file,
	}
}