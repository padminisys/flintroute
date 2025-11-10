# FlintRoute Functional Testing Framework

This directory contains the functional testing framework for FlintRoute, providing comprehensive end-to-end testing capabilities for the BGP management system.

## Overview

The functional testing framework validates FlintRoute's behavior through realistic scenarios, including:
- Authentication and authorization
- BGP peer management
- Session state management
- Configuration operations
- Alert handling
- Error scenarios
- Complete workflows

## Quick Start

### Prerequisites

- Go 1.21 or later
- SQLite3
- Make

### Running Tests

```bash
# Run all functional tests
make test-functional

# Run specific test suite
go test -v ./test/functional/tests/01_authentication/...

# Run with custom configuration
TEST_CONFIG=./config/custom-config.yaml go test -v ./test/functional/tests/...
```

## Directory Structure

```
test/functional/
├── cmd/mock-frr-server/    # Mock FRR gRPC server for testing
├── pkg/client/             # API client for testing
├── pkg/testutil/           # Testing utilities and helpers
├── pkg/runner/             # Test execution framework
├── tests/                  # Test suites organized by feature
│   ├── 01_authentication/
│   ├── 02_peer_management/
│   ├── 03_session_management/
│   ├── 04_configuration/
│   ├── 05_alerts/
│   ├── 06_error_handling/
│   └── 07_workflows/
├── fixtures/               # Test data and fixtures
│   ├── peers/
│   ├── users/
│   └── sessions/
├── scripts/                # Helper scripts
├── config/                 # Test configurations
├── logs/                   # Test execution logs
├── results/                # Test results and reports
└── tmp/                    # Temporary files (databases, etc.)
```

## Configuration

Test configuration is managed through YAML files in the [`config/`](config/) directory:

- [`test-config.yaml`](config/test-config.yaml) - Main test configuration
- [`mock-frr-config.yaml`](config/mock-frr-config.yaml) - Mock FRR server settings
- [`logging-config.yaml`](config/logging-config.yaml) - Logging configuration

## Test Organization

Tests are organized into numbered suites for sequential execution:

1. **Authentication** - User login, token management, permissions
2. **Peer Management** - CRUD operations for BGP peers
3. **Session Management** - Session state tracking and updates
4. **Configuration** - Configuration validation and application
5. **Alerts** - Alert generation and handling
6. **Error Handling** - Error scenarios and recovery
7. **Workflows** - End-to-end user workflows

## Writing Tests

### Test Structure

```go
func TestFeature(t *testing.T) {
    // Setup
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Execute
    result, err := ctx.Client.DoSomething()
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Using Fixtures

```go
// Load peer fixture
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/peer1.json")

// Use in test
err := ctx.Client.CreatePeer(peer)
require.NoError(t, err)
```

## Mock FRR Server

The mock FRR server simulates FRRouting's gRPC interface for testing without requiring a real FRR installation.

### Starting the Mock Server

```bash
cd cmd/mock-frr-server
go run main.go
```

### Features

- Simulates BGP session state changes
- Configurable delays and error injection
- Supports all FRR gRPC operations used by FlintRoute

## Test Results

Test results are stored in the [`results/`](results/) directory in multiple formats:

- JSON - Machine-readable results
- XML - JUnit format for CI integration
- HTML - Human-readable reports

## Logging

Test execution logs are written to [`logs/`](logs/):

- `test-execution.log` - Main test execution log
- `mock-frr-server.log` - Mock server activity
- Individual test logs as needed

## CI Integration

The functional tests are integrated into the CI pipeline:

```yaml
# .github/workflows/test.yml
- name: Run Functional Tests
  run: make test-functional
```

## Troubleshooting

### Tests Hanging

- Check if mock FRR server is running
- Verify port availability (default: 50051)
- Review timeout settings in configuration

### Database Errors

- Ensure [`tmp/`](tmp/) directory is writable
- Check for leftover database files from previous runs
- Run cleanup: `make clean-test-artifacts`

### Authentication Failures

- Verify JWT secret matches between test config and server
- Check token expiry settings
- Review user fixtures for correct credentials

## Contributing

When adding new tests:

1. Place tests in the appropriate numbered suite
2. Add fixtures to [`fixtures/`](fixtures/) if needed
3. Update this README with new test categories
4. Ensure tests clean up resources properly
5. Add appropriate logging for debugging

## Resources

- [Testing Guide](../../docs/testing/TESTING_GUIDE.md)
- [Architecture Overview](../../docs/architecture/overview.md)
- [Development Setup](../../docs/development/setup.md)