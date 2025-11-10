package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/flintroute/test/functional/pkg/client"
	"github.com/yourusername/flintroute/test/functional/pkg/testutil"
)

// TestExecutor manages test execution
type TestExecutor struct {
	config      *TestConfig
	apiClient   *client.APIClient
	dbManager   *testutil.DatabaseManager
	logger      *testutil.TestLogger
	results     *TestResults
	fixtureLoader *testutil.FixtureLoader
}

// NewTestExecutor creates a new test executor
func NewTestExecutor(config *TestConfig) (*TestExecutor, error) {
	if config == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	executor := &TestExecutor{
		config:  config,
		results: NewTestResults(),
	}

	return executor, nil
}

// Setup initializes the test environment
func (e *TestExecutor) Setup() error {
	// Create necessary directories
	dirs := []string{
		e.config.LogsPath,
		e.config.ResultsPath,
		filepath.Dir(e.config.DatabasePath),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Initialize logger
	logPath := filepath.Join(e.config.LogsPath, fmt.Sprintf("test-%s.log", time.Now().Format("20060102-150405")))
	logger, err := testutil.NewTestLogger(logPath, e.config.LogLevel)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	e.logger = logger
	e.logger.Info("Test executor initialized")

	// Initialize API client
	e.apiClient = client.NewAPIClient(e.config.ServerURL, logger.GetZapLogger())
	e.apiClient.SetTimeout(e.config.Timeout)
	e.logger.Info("API client initialized")

	// Initialize database manager
	dbManager, err := testutil.NewDatabaseManager(e.config.DatabasePath, logger.GetZapLogger())
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}
	e.dbManager = dbManager

	if err := e.dbManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	e.logger.Info("Database initialized")

	// Initialize fixture loader
	e.fixtureLoader = testutil.NewFixtureLoader(e.config.FixturesPath, logger.GetZapLogger())
	e.logger.Info("Fixture loader initialized")

	// Verify server is reachable
	if err := e.apiClient.HealthCheck(); err != nil {
		return fmt.Errorf("server health check failed: %w", err)
	}
	e.logger.Info("Server health check passed")

	return nil
}

// Teardown cleans up the test environment
func (e *TestExecutor) Teardown() error {
	e.logger.Info("Starting teardown")

	// Close database
	if e.dbManager != nil {
		if err := e.dbManager.Close(); err != nil {
			e.logger.Error("Failed to close database")
		}
	}

	// Cleanup database file if configured
	if e.config.CleanupOnSuccess && !e.results.HasFailures() {
		if err := os.Remove(e.config.DatabasePath); err != nil && !os.IsNotExist(err) {
			e.logger.Warn("Failed to remove database file")
		} else {
			e.logger.Info("Database file removed")
		}
	}

	// Close logger
	if e.logger != nil {
		if err := e.logger.Close(); err != nil {
			return fmt.Errorf("failed to close logger: %w", err)
		}
	}

	return nil
}

// RunTests discovers and runs tests matching the pattern
func (e *TestExecutor) RunTests(pattern string) error {
	e.logger.Info("Starting test run")

	// Discover tests
	tests, err := e.DiscoverTests(pattern)
	if err != nil {
		return fmt.Errorf("failed to discover tests: %w", err)
	}

	if len(tests) == 0 {
		e.logger.Warn("No tests found matching pattern")
		return nil
	}

	e.logger.Info("Tests discovered")

	// Run tests
	for _, testPath := range tests {
		result, err := e.ExecuteTest(testPath)
		if err != nil {
			e.logger.Error("Failed to execute test")
			result = &TestResult{
				Name:     testPath,
				Status:   "failed",
				Error:    err.Error(),
				Duration: 0,
			}
		}
		e.results.AddResult(result)
	}

	// Finalize results
	e.results.Finalize()

	return nil
}

// DiscoverTests finds all test files matching the pattern
func (e *TestExecutor) DiscoverTests(pattern string) ([]string, error) {
	testsDir := "./tests"
	var tests []string

	// If pattern is empty, use wildcard
	if pattern == "" {
		pattern = "*"
	}

	// Walk through tests directory
	err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file matches pattern
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if matched && strings.HasSuffix(path, "_test.go") {
			tests = append(tests, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk tests directory: %w", err)
	}

	return tests, nil
}

// ExecuteTest executes a single test file
func (e *TestExecutor) ExecuteTest(testPath string) (*TestResult, error) {
	startTime := time.Now()
	testName := filepath.Base(testPath)

	e.logger.LogTestStart(testName)

	result := &TestResult{
		Name:   testName,
		Status: "passed",
	}

	// Note: In a real implementation, this would execute the Go test file
	// For now, this is a placeholder that would need to be integrated with
	// the actual test execution mechanism (e.g., using go test command or
	// importing and running test functions directly)

	// Placeholder implementation
	e.logger.Info("Test execution placeholder")

	result.Duration = time.Since(startTime)
	e.logger.LogTestEnd(testName, result.Status == "passed", result.Duration)

	return result, nil
}

// GetResults returns the test results
func (e *TestExecutor) GetResults() *TestResults {
	return e.results
}

// GetAPIClient returns the API client
func (e *TestExecutor) GetAPIClient() *client.APIClient {
	return e.apiClient
}

// GetDatabaseManager returns the database manager
func (e *TestExecutor) GetDatabaseManager() *testutil.DatabaseManager {
	return e.dbManager
}

// GetLogger returns the logger
func (e *TestExecutor) GetLogger() *testutil.TestLogger {
	return e.logger
}

// GetFixtureLoader returns the fixture loader
func (e *TestExecutor) GetFixtureLoader() *testutil.FixtureLoader {
	return e.fixtureLoader
}

// CleanDatabase cleans all data from the database
func (e *TestExecutor) CleanDatabase() error {
	if e.dbManager == nil {
		return fmt.Errorf("database manager not initialized")
	}
	return e.dbManager.Clean()
}

// GenerateReports generates test reports in multiple formats
func (e *TestExecutor) GenerateReports() error {
	timestamp := time.Now().Format("20060102-150405")

	// Generate JSON report
	jsonPath := filepath.Join(e.config.ResultsPath, fmt.Sprintf("results-%s.json", timestamp))
	if err := e.results.GenerateJSONReport(jsonPath); err != nil {
		return fmt.Errorf("failed to generate JSON report: %w", err)
	}
	e.logger.Info("JSON report generated")

	// Generate XML report
	xmlPath := filepath.Join(e.config.ResultsPath, fmt.Sprintf("results-%s.xml", timestamp))
	if err := e.results.GenerateXMLReport(xmlPath); err != nil {
		return fmt.Errorf("failed to generate XML report: %w", err)
	}
	e.logger.Info("XML report generated")

	// Print summary
	e.results.PrintSummary()

	return nil
}