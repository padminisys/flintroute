# FlintRoute Testing Guide

## Overview

This guide provides comprehensive information about the testing framework implemented for the FlintRoute Go backend. The testing infrastructure ensures code quality, functional correctness, and maintainability.

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Test Coverage](#test-coverage)
5. [Writing Tests](#writing-tests)
6. [Mocking](#mocking)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)

## Testing Philosophy

The FlintRoute testing framework follows these principles:

- **Unit Tests Only**: Focus on testing individual components in isolation
- **Fast Execution**: All tests should complete in under 30 seconds
- **No External Dependencies**: All external services (FRR, databases) are mocked
- **High Coverage**: Target minimum 80% code coverage
- **Clear and Maintainable**: Tests should be easy to read and understand

## Test Structure

### Directory Organization

```
internal/
├── models/
│   ├── models.go
│   └── models_test.go
├── auth/
│   ├── jwt.go
│   ├── jwt_test.go
│   ├── middleware.go
│   └── middleware_test.go
├── config/
│   ├── config.go
│   └── config_test.go
├── database/
│   ├── database.go
│   └── database_test.go
├── frr/
│   ├── client.go
│   ├── client_test.go
│   └── mock_client.go
├── bgp/
│   ├── service.go
│   └── service_test.go
├── websocket/
│   ├── hub.go
│   ├── hub_test.go
│   ├── handler.go
│   └── handler_test.go
├── api/
│   ├── auth_handlers.go
│   ├── auth_handlers_test.go
│   ├── bgp_handlers.go
│   ├── bgp_handlers_test.go
│   ├── config_handlers.go
│   └── config_handlers_test.go
└── testutil/
    └── helpers.go
```

### Test File Naming

- Test files are named `*_test.go`
- Mock files are named `mock_*.go`
- Test files are placed alongside the source files they test

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./internal/...

# Run tests with verbose output
go test ./internal/... -v

# Run tests for a specific package
go test ./internal/auth -v

# Run a specific test
go test ./internal/auth -v -run TestGenerateToken

# Run tests with race detector
go test ./internal/... -race
```

### Using Makefile

```bash
# Run all tests
make -f Makefile.test test

# Run tests with coverage
make -f Makefile.test test-coverage

# Generate HTML coverage report
make -f Makefile.test test-coverage-html

# Run tests with race detector
make -f Makefile.test test-race

# Check coverage threshold (80%)
make -f Makefile.test test-coverage-check
```

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage profile
go test ./internal/... -coverprofile=coverage.out

# View coverage summary
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Targets

- **Overall Target**: 80% minimum
- **Critical Packages**: 90%+ (auth, database, api)
- **Utility Packages**: 70%+ (config, models)

### Current Coverage Status

Run `make -f Makefile.test test-coverage` to see current coverage statistics.

## Writing Tests

### Test Structure

```go
func TestFunctionName(t *testing.T) {
    t.Run("descriptive test case name", func(t *testing.T) {
        // Arrange
        // ... setup test data
        
        // Act
        // ... call function under test
        
        // Assert
        // ... verify results
    })
}
```

### Using Testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
    result := SomeFunction()
    
    assert.NotNil(t, result)
    assert.Equal(t, expected, result)
    assert.NoError(t, err)
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid input", "test", true, false},
        {"invalid input", "", false, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

## Mocking

### Using Test Helpers

```go
import "github.com/padminisys/flintroute/internal/testutil"

func TestWithDatabase(t *testing.T) {
    db := testutil.SetupTestDB(t)
    defer testutil.CleanupTestDB(t, db)
    
    // Use db for testing
}
```

### FRR Client Mock

```go
import "github.com/padminisys/flintroute/internal/frr"

func TestWithMockFRR(t *testing.T) {
    mockClient := frr.NewMockClient()
    mockClient.On("AddBGPPeer", ctx, config).Return(nil)
    
    // Use mockClient in tests
    
    mockClient.AssertExpectations(t)
}
```

### HTTP Handler Testing

```go
func TestHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/endpoint", handler)
    
    req := httptest.NewRequest(http.MethodPost, "/endpoint", body)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Best Practices

### 1. Test Independence

- Each test should be independent and not rely on other tests
- Use `t.Run()` for subtests
- Clean up resources after tests

### 2. Clear Test Names

```go
// Good
t.Run("returns error when user not found", func(t *testing.T) {})

// Bad
t.Run("test1", func(t *testing.T) {})
```

### 3. Use Test Helpers

```go
// Create reusable setup functions
func setupTestServer(t *testing.T) *Server {
    t.Helper()
    // ... setup code
    return server
}
```

### 4. Test Error Cases

```go
func TestFunction(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // Test happy path
    })
    
    t.Run("error case - invalid input", func(t *testing.T) {
        // Test error handling
    })
}
```

### 5. Avoid Test Flakiness

- Don't rely on timing (use channels/synchronization)
- Don't use random data without seeding
- Clean up resources properly

### 6. Keep Tests Fast

- Use in-memory databases
- Mock external services
- Avoid unnecessary sleeps

## Troubleshooting

### Common Issues

#### 1. Database Lock Errors

**Problem**: SQLite database locked errors

**Solution**: Use separate database files for each test or use in-memory databases

```go
db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
```

#### 2. Port Already in Use

**Problem**: Test server fails to start

**Solution**: Use random ports or ensure cleanup

```go
listener, _ := net.Listen("tcp", ":0")
port := listener.Addr().(*net.TCPAddr).Port
```

#### 3. Race Conditions

**Problem**: Tests fail intermittently

**Solution**: Use proper synchronization

```go
done := make(chan bool)
go func() {
    // async operation
    done <- true
}()
<-done // Wait for completion
```

#### 4. Mock Expectations Not Met

**Problem**: `AssertExpectations` fails

**Solution**: Ensure all mocked methods are called

```go
mockClient.On("Method", args).Return(result)
// ... call code that uses Method
mockClient.AssertExpectations(t)
```

### Running Specific Tests

```bash
# Run only failing tests
go test ./internal/... -v -run TestName

# Run with more verbose output
go test ./internal/... -v -count=1

# Skip cache
go test ./internal/... -count=1
```

### Debugging Tests

```go
// Add debug output
t.Logf("Debug: value = %v", value)

// Use testify's require for immediate failure
require.NoError(t, err) // Stops test immediately on failure
```

## Test Utilities

### Available Helpers

Located in `internal/testutil/helpers.go`:

- `SetupTestDB(t)` - Creates test database
- `SetupTestDBWithData(t)` - Creates database with sample data
- `CleanupTestDB(t, db)` - Closes database connection
- `CreateTestLogger()` - Creates no-op logger
- `CreateTestConfig(t, content)` - Creates temporary config file
- `SetupInMemoryDB(t)` - Creates pure in-memory database

### Example Usage

```go
func TestWithHelpers(t *testing.T) {
    db, user, peer := testutil.SetupTestDBWithData(t)
    defer testutil.CleanupTestDB(t, db)
    
    logger := testutil.CreateTestLogger()
    
    // Use db, user, peer, logger in tests
}
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.24'
      - run: go test ./internal/... -v -coverprofile=coverage.out
      - run: go tool cover -func=coverage.out
```

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Test Comments](https://golang.org/cmd/go/#hdr-Testing_flags)

## Contributing

When adding new features:

1. Write tests first (TDD approach recommended)
2. Ensure all tests pass
3. Maintain or improve coverage
4. Update this documentation if needed

## Support

For questions or issues with tests:

1. Check this documentation
2. Review existing test examples
3. Check the troubleshooting section
4. Open an issue on GitHub