# Frequently Asked Questions (FAQ)

**Common Questions About FlintRoute Functional Testing**

---

## Table of Contents

1. [General Questions](#general-questions)
2. [Setup and Installation](#setup-and-installation)
3. [Running Tests](#running-tests)
4. [Writing Tests](#writing-tests)
5. [Mock FRR Server](#mock-frr-server)
6. [Troubleshooting](#troubleshooting)
7. [Performance](#performance)
8. [CI/CD Integration](#cicd-integration)

---

## General Questions

### What is the FlintRoute Functional Testing Framework?

The FlintRoute Functional Testing Framework is a comprehensive end-to-end testing solution for the FlintRoute BGP management system. It provides realistic testing capabilities through a mock FRR server, extensive test fixtures, and complete test automation.

### Why do we need functional tests?

Functional tests validate that the entire system works correctly together, including:
- API endpoints
- Database operations
- FRR integration
- Authentication and authorization
- Real-world workflows

Unlike unit tests that test individual components, functional tests ensure the complete system behaves correctly.

### What's the difference between unit tests and functional tests?

| Aspect | Unit Tests | Functional Tests |
|--------|-----------|------------------|
| **Scope** | Single function/method | Complete system |
| **Dependencies** | Mocked | Real (or realistic mock) |
| **Speed** | Very fast (< 1ms) | Slower (seconds) |
| **Purpose** | Verify logic | Verify behavior |
| **Location** | `internal/*/` | `test/functional/` |

### Do I need to run functional tests for every change?

**For development**: Run relevant test suites for your changes.  
**Before committing**: Run all tests with `./run-clean.sh`.  
**In CI/CD**: All tests run automatically.

---

## Setup and Installation

### What are the prerequisites?

**Required**:
- Go 1.21 or later
- SQLite3
- Make
- Git

**Optional but recommended**:
- curl (for manual testing)
- jq (for JSON processing)
- lsof (for port checking)

Check with: `./scripts/check-prerequisites.sh`

### How do I install Go 1.21?

**Linux**:
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**macOS**:
```bash
brew install go@1.21
```

**Verify**:
```bash
go version  # Should show 1.21 or later
```

### Do I need a real FRR installation?

**No!** The mock FRR server simulates FRR's behavior, so you don't need a real FRR installation for testing.

### Can I run tests on Windows?

Yes, but you need **WSL2** (Windows Subsystem for Linux). The tests are designed for Unix-like environments.

### How much disk space do I need?

- **Minimum**: 100MB
- **Recommended**: 500MB (for logs and artifacts)

---

## Running Tests

### How do I run all tests?

```bash
cd test/functional
./run-clean.sh
```

This performs a clean run with full cleanup before testing.

### How do I run a specific test suite?

```bash
./run-tests.sh --pattern ./tests/01_authentication/...
```

### How do I run a single test?

```bash
cd tests/01_authentication
go test -v -run TestLogin
```

### How long do tests take to run?

- **Full suite**: 2-5 minutes
- **Single suite**: 10-30 seconds
- **Single test**: < 1 second (most tests)

### Can I run tests in parallel?

Currently, tests run sequentially by default. Parallel execution is planned for future releases.

### How do I see more detailed output?

```bash
# Verbose output
./run-tests.sh --verbose

# Debug logging
./run-tests.sh --log-level debug

# Both
./run-tests.sh --verbose --log-level debug
```

### Where are test results stored?

- **JSON results**: `results/test-results-TIMESTAMP.json`
- **HTML reports**: `results/test-results-TIMESTAMP.html`
- **Text summary**: `results/test-summary.txt`
- **Logs**: `logs/test-execution.log`

### How do I keep artifacts after test failure?

```bash
./run-tests.sh --no-cleanup
```

This preserves:
- Test database (`tmp/test.db`)
- Log files (`logs/`)
- Test results (`results/`)

---

## Writing Tests

### How do I create a new test?

1. **Choose the appropriate suite**:
   ```bash
   cd tests/02_peer_management
   ```

2. **Create test file**:
   ```bash
   touch my_feature_test.go
   ```

3. **Write test**:
   ```go
   package peermanagement_test
   
   import (
       "testing"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
   )
   
   func TestMyFeature(t *testing.T) {
       ctx := testutil.NewTestContext(t)
       defer ctx.Cleanup()
       
       // Your test code
   }
   ```

### How do I use fixtures in tests?

```go
// Load peer fixture
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")

// Use in test
err := ctx.Client.CreatePeer(peer)
require.NoError(t, err)
```

### How do I test error scenarios?

```go
func TestInvalidPeer(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Load invalid fixture
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/invalid/invalid_ip_format.yaml")
    
    // Expect error
    err := ctx.Client.CreatePeer(peer)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid IP address")
}
```

### How do I wait for async operations?

```go
// Wait for session to establish
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
require.NoError(t, err)

// Or use generic wait
err := testutil.WaitForCondition(t, 10*time.Second, func() bool {
    session, _ := ctx.Client.GetSession("peer1")
    return session.State == "Established"
})
```

### Should I use `assert` or `require`?

- **Use `require`** for critical checks that should stop the test:
  ```go
  require.NoError(t, err)  // Stop if error
  require.NotNil(t, peer)  // Stop if nil
  ```

- **Use `assert`** for non-critical checks:
  ```go
  assert.Equal(t, "Established", session.State)  // Continue if fails
  assert.Greater(t, count, 0)  // Continue if fails
  ```

---

## Mock FRR Server

### What is the Mock FRR Server?

A complete simulation of FRRouting's gRPC interface that allows testing without a real FRR installation. It simulates BGP session states, peer management, and configuration operations.

### How do I start the Mock FRR Server?

```bash
cd cmd/mock-frr-server
./mock-frr-server
```

Or use the test runner which starts it automatically.

### How do I test the Mock FRR Server manually?

```bash
# Health check
curl http://localhost:51051/health

# View statistics
curl http://localhost:51051/stats

# List peers
curl http://localhost:51051/peers

# Add a peer
curl -X POST http://localhost:51051/peers/add \
  -H "Content-Type: application/json" \
  -d '{"IPAddress":"192.168.1.1","ASN":65000,"RemoteASN":65001}'
```

### Can I configure the Mock FRR Server?

Yes! Edit `config/mock-frr-config.yaml`:

```yaml
server:
  host: localhost
  port: 50051
  
simulation:
  session_state_delay: 100ms  # Adjust delay
  error_injection: false      # Enable for error testing
  
logging:
  level: info                 # debug, info, warn, error
  file: ./logs/mock-frr-server.log
```

### How does BGP session simulation work?

The mock server simulates realistic BGP state transitions:

1. **Idle** (initial state)
2. **Connect** (after delay)
3. **Active** (after delay)
4. **OpenSent** (after delay)
5. **OpenConfirm** (after delay)
6. **Established** (after delay)

Default delay: 100ms per state (configurable)

### Can I inject errors for testing?

Yes! Enable error injection in config:

```yaml
simulation:
  error_injection: true
```

All operations will return errors, allowing you to test error handling.

---

## Troubleshooting

### Tests are hanging, what should I do?

**Check if mock server is running**:
```bash
ps aux | grep mock-frr-server
```

**Check ports**:
```bash
lsof -i :50051
lsof -i :51051
```

**Restart mock server**:
```bash
cd cmd/mock-frr-server
./mock-frr-server &
```

### I'm getting "database locked" errors

**Clean the database**:
```bash
./scripts/cleanup-db.sh
```

**Check for multiple test instances**:
```bash
ps aux | grep "go test"
```

**Ensure proper cleanup**:
```go
func TestExample(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()  // Always defer cleanup!
}
```

### I'm getting "port already in use" errors

**Find the process**:
```bash
lsof -i :8080
lsof -i :50051
```

**Kill the process**:
```bash
kill -9 <PID>
```

**Or cleanup all**:
```bash
./scripts/cleanup-all.sh
```

### Authentication tests are failing

**Check JWT secret**:
```bash
grep jwt_secret config/test-config.yaml
```

**Verify user fixtures**:
```bash
cat fixtures/users/admin_user.yaml
```

**Check token expiry**:
```yaml
auth:
  token_expiry: 5m  # Increase if needed
```

### Tests pass locally but fail in CI

**Common causes**:
1. **Timing issues**: Increase timeouts in CI
2. **Port conflicts**: Ensure ports are available
3. **Resource limits**: Check CI resource allocation
4. **Environment differences**: Verify Go version, dependencies

**Debug in CI**:
```bash
# Add verbose logging
./run-tests.sh --verbose --log-level debug

# Preserve artifacts
./run-tests.sh --no-cleanup
```

### How do I debug a specific test?

**Add debug logging**:
```go
func TestDebug(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    ctx.Logger.SetLevel(logrus.DebugLevel)
    defer ctx.Cleanup()
    
    t.Logf("Debug info: %+v", someVariable)
}
```

**Use delve debugger**:
```bash
cd tests/01_authentication
dlv test -- -test.run TestLogin
```

**Inspect database**:
```bash
sqlite3 tmp/test.db
SELECT * FROM bgp_peers;
```

---

## Performance

### Why are tests slow?

**Common reasons**:
1. **BGP simulation delays**: Configurable in mock-frr-config.yaml
2. **Database operations**: SQLite is slower than in-memory
3. **Network operations**: HTTP requests take time
4. **Sequential execution**: Tests run one at a time

### How can I speed up tests?

**Reduce simulation delays**:
```yaml
simulation:
  session_state_delay: 10ms  # Reduce from 100ms
```

**Run specific suites**:
```bash
./run-tests.sh --pattern ./tests/01_authentication/...
```

**Skip cleanup during development**:
```bash
./run-tests.sh --no-cleanup
```

### Can tests run in parallel?

Not currently, but it's planned for future releases. Tests share resources (database, mock server) that would need isolation for parallel execution.

---

## CI/CD Integration

### How do I integrate with GitHub Actions?

Create `.github/workflows/functional-tests.yml`:

```yaml
name: Functional Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          cd test/functional
          ./run-clean.sh
```

### How do I integrate with GitLab CI?

Create `.gitlab-ci.yml`:

```yaml
functional-tests:
  stage: test
  image: golang:1.21
  script:
    - cd test/functional
    - ./run-clean.sh
  artifacts:
    when: always
    paths:
      - test/functional/results/
```

### How do I generate JUnit XML reports?

```bash
# Install go-junit-report
go install github.com/jstemmer/go-junit-report@latest

# Generate report
go test -v ./... 2>&1 | go-junit-report > results/junit.xml
```

### Should tests run on every commit?

**Recommended approach**:
- **On every push**: Run all tests
- **On pull requests**: Run all tests
- **On main branch**: Run all tests + generate reports

### How do I handle flaky tests in CI?

**Short term**:
- Increase timeouts
- Add retries for network operations
- Improve test isolation

**Long term**:
- Identify root cause
- Fix underlying issues
- Improve test reliability

---

## Additional Questions

### Where can I find more help?

- **Documentation**: [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **API Reference**: [API_REFERENCE.md](API_REFERENCE.md)
- **Quick Reference**: [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)

### How do I report a bug?

1. Check if it's already reported in GitHub issues
2. Create a new issue with:
   - Clear description
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details
   - Relevant logs

### How do I request a feature?

1. Check if it's already requested
2. Create a feature request with:
   - Clear description
   - Use case
   - Proposed solution
   - Alternatives considered

### How do I contribute?

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:
- Setting up development environment
- Writing tests
- Submitting pull requests
- Code style and standards

### Is there a community chat?

Check the main FlintRoute repository for:
- GitHub Discussions
- Community channels
- Support resources

---

## Quick Tips

### Best Practices

✅ **DO**:
- Always use `defer ctx.Cleanup()`
- Load fixtures for test data
- Use `require` for critical checks
- Wait for async operations
- Add descriptive assertion messages

❌ **DON'T**:
- Share state between tests
- Hardcode test data
- Use `time.Sleep()` (use `WaitForCondition`)
- Ignore errors
- Make tests too complex

### Common Patterns

**Basic test**:
```go
func TestFeature(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Test code
}
```

**Parameterized test**:
```go
testCases := []struct {
    name    string
    fixture string
    wantErr bool
}{
    {"Valid", "fixtures/valid.yaml", false},
    {"Invalid", "fixtures/invalid.yaml", true},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        // Test code
    })
}
```

**Async operation**:
```go
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
require.NoError(t, err)
```

---

## Still Have Questions?

If your question isn't answered here:

1. **Check the documentation**:
   - [TESTING_GUIDE.md](TESTING_GUIDE.md) - Comprehensive guide
   - [API_REFERENCE.md](API_REFERENCE.md) - Complete API docs
   - [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Quick reference

2. **Search existing issues**: Someone may have asked before

3. **Ask in discussions**: GitHub Discussions for questions

4. **Create an issue**: For bugs or feature requests

---

**Document Version**: 1.0.0  
**Last Updated**: November 10, 2025  
**Maintained By**: FlintRoute Development Team