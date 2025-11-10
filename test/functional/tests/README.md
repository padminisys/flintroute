# FlintRoute Functional Tests

This directory contains functional tests for the FlintRoute BGP management system. Tests are organized by feature area and numbered for execution order.

## Directory Structure

```
tests/
├── 01_authentication/     # Authentication and authorization tests
├── 02_peer_management/    # BGP peer CRUD operations
├── 03_session_management/ # BGP session monitoring
├── 04_configuration/      # Configuration backup/restore
├── 05_alerts/            # Alert system tests
├── 06_error_handling/    # Error scenarios and edge cases
└── 07_workflows/         # End-to-end workflow tests
```

## Test Organization

Tests are organized using a numbered prefix system:
- `01_` - Core authentication and setup
- `02_` - Basic CRUD operations
- `03_` - Session management
- `04_` - Configuration management
- `05_` - Alert handling
- `06_` - Error scenarios
- `07_` - Complex workflows

## Running Tests

### Run All Tests
```bash
cd test/functional
go test ./tests/... -v
```

### Run Specific Test Suite
```bash
cd test/functional
go test ./tests/01_authentication/... -v
```

### Run Single Test
```bash
cd test/functional
go test ./tests/01_authentication -run TestLogin -v
```

### Run with Coverage
```bash
cd test/functional
go test ./tests/... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

Each test file follows this pattern:

```go
package <feature>_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/yourusername/flintroute/test/functional/pkg/client"
    "github.com/yourusername/flintroute/test/functional/pkg/testutil"
)

func TestFeature(t *testing.T) {
    // Setup
    logger, err := testutil.NewTestLogger("../../logs/test.log", "debug")
    require.NoError(t, err)
    defer logger.Close()
    
    apiClient := client.NewAPIClient("http://localhost:8080", logger.GetZapLogger())
    
    // Load fixtures
    fixtureLoader := testutil.NewFixtureLoader("../../fixtures", logger.GetZapLogger())
    
    // Test cases
    t.Run("test_case_name", func(t *testing.T) {
        // Test implementation
    })
}
```

## Prerequisites

Before running tests, ensure:

1. **FlintRoute Server is Running**
   ```bash
   # From project root
   make run
   ```

2. **Mock FRR Server is Running**
   ```bash
   cd test/functional/cmd/mock-frr-server
   make run
   ```

3. **Database is Clean**
   ```bash
   cd test/functional
   ./scripts/cleanup-db.sh
   ```

4. **Dependencies are Installed**
   ```bash
   cd test/functional
   go mod download
   ```

## Test Utilities

### Logger
The `testutil.TestLogger` provides structured logging:
- Logs to both file and console
- Supports debug, info, warn, error levels
- Includes test lifecycle logging

### API Client
The `client.APIClient` provides:
- Automatic authentication
- Token refresh
- Request/response logging
- Type-safe API methods

### Fixture Loader
The `testutil.FixtureLoader` loads test data from YAML files:
- User fixtures
- Peer fixtures
- Session fixtures
- Alert fixtures

## Writing New Tests

1. **Create Test File**
   ```bash
   touch tests/XX_feature/YY_test_name_test.go
   ```

2. **Use Test Template**
   - Import required packages
   - Setup logger and client
   - Load fixtures
   - Write test cases using subtests

3. **Follow Naming Conventions**
   - Test functions: `TestFeatureName`
   - Subtests: `test_case_description`
   - Files: `XX_descriptive_name_test.go`

4. **Use Assertions**
   - `require.*` for critical checks (stops test on failure)
   - `assert.*` for non-critical checks (continues test)

## Example Test

See [`01_authentication/01_login_test.go`](./01_authentication/01_login_test.go) for a complete example demonstrating:
- Logger setup
- API client usage
- Fixture loading
- Multiple test cases
- Proper assertions
- Error handling

## Debugging Tests

### Enable Verbose Logging
```bash
go test ./tests/... -v -args -log-level=debug
```

### View Test Logs
```bash
tail -f logs/test.log
```

### Run Single Test with Debugging
```bash
go test ./tests/01_authentication -run TestLogin -v
```

## CI/CD Integration

Tests can be run in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Functional Tests
  run: |
    cd test/functional
    ./scripts/setup-env.sh
    go test ./tests/... -v
    ./scripts/teardown-env.sh
```

## Test Data

Test fixtures are located in `test/functional/fixtures/`:
- `users/` - User credentials and profiles
- `peers/` - BGP peer configurations
- `sessions/` - BGP session states
- `alerts/` - Alert definitions
- `config/` - Configuration backups

## Troubleshooting

### Tests Fail to Connect
- Verify FlintRoute server is running on `http://localhost:8080`
- Check Mock FRR server is running on `http://localhost:50051`
- Review logs in `test/functional/logs/`

### Authentication Failures
- Ensure database is clean before tests
- Verify user fixtures are correct
- Check JWT secret configuration

### Fixture Loading Errors
- Verify fixture files exist in `fixtures/` directory
- Check YAML syntax is valid
- Review fixture loader logs

## Contributing

When adding new tests:
1. Follow the existing structure
2. Add comprehensive test cases
3. Include both positive and negative scenarios
4. Update this README if adding new test categories
5. Ensure tests are idempotent and can run in any order

## Resources

- [Testing Guide](../TESTING_GUIDE.md)
- [API Reference](../API_REFERENCE.md)
- [Fixture Documentation](../fixtures/README.md)
- [Test Runner Documentation](../pkg/runner/README.md)