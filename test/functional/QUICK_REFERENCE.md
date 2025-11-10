# FlintRoute Testing Quick Reference

**One-Page Cheat Sheet for Functional Testing**

---

## Quick Start

```bash
# Check prerequisites
./scripts/check-prerequisites.sh

# Run all tests (clean environment)
./run-clean.sh

# Run specific test suite
./run-tests.sh --pattern ./tests/01_authentication/...

# Run with debug logging
./run-tests.sh --log-level debug --verbose
```

---

## Common Commands

### Test Execution

```bash
# Full test run
./run-tests.sh

# Clean run (recommended)
./run-clean.sh

# Specific pattern
./run-tests.sh --pattern ./tests/02_peer_management/...

# Keep artifacts on failure
./run-tests.sh --no-cleanup

# Custom config
./run-tests.sh --config config/custom-config.yaml
```

### Environment Management

```bash
# Check prerequisites
./scripts/check-prerequisites.sh

# Setup environment
./scripts/setup-env.sh

# Teardown environment
./scripts/teardown-env.sh

# Full cleanup
./scripts/cleanup-all.sh

# Clean specific artifacts
./scripts/cleanup-db.sh
./scripts/cleanup-logs.sh
./scripts/cleanup-results.sh
```

### Mock FRR Server

```bash
# Start server
cd cmd/mock-frr-server
./mock-frr-server

# Build server
make build

# Test server
./test-server.sh

# Check health
curl http://localhost:51051/health

# View stats
curl http://localhost:51051/stats
```

---

## Directory Structure

```
test/functional/
├── cmd/mock-frr-server/    # Mock FRR server
├── config/                 # Configuration files
├── fixtures/               # Test data (YAML)
│   ├── peers/             # BGP peer fixtures
│   ├── users/             # User fixtures
│   ├── sessions/          # Session fixtures
│   └── alerts/            # Alert fixtures
├── logs/                   # Test logs
├── results/                # Test results
├── scripts/                # Helper scripts
├── tests/                  # Test suites
│   ├── 01_authentication/
│   ├── 02_peer_management/
│   ├── 03_session_management/
│   ├── 04_configuration/
│   ├── 05_alerts/
│   ├── 06_error_handling/
│   └── 07_workflows/
├── tmp/                    # Temporary files
├── run-tests.sh           # Main test runner
└── run-clean.sh           # Clean run script
```

---

## Test Template

```go
package mytest_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFeature(t *testing.T) {
    // Setup
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Arrange
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    
    // Act
    err := ctx.Client.CreatePeer(peer)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, peer)
}
```

---

## API Client Methods

### Authentication
```go
token, err := ctx.Client.Login("admin", "password")
err = ctx.Client.Logout()
newToken, err := ctx.Client.RefreshToken()
isAuth := ctx.Client.IsAuthenticated()
```

### Peer Management
```go
err = ctx.Client.CreatePeer(peer)
peer, err := ctx.Client.GetPeer("peer1")
peers, err := ctx.Client.ListPeers()
err = ctx.Client.UpdatePeer(peer)
err = ctx.Client.DeletePeer("peer1")
err = ctx.Client.EnablePeer("peer1")
err = ctx.Client.DisablePeer("peer1")
```

### Session Management
```go
session, err := ctx.Client.GetSession("peer1")
sessions, err := ctx.Client.ListSessions()
err = ctx.Client.ResetSession("peer1", false)  // soft
err = ctx.Client.ResetSession("peer1", true)   // hard
```

### Configuration
```go
config, err := ctx.Client.GetConfiguration()
err = ctx.Client.UpdateConfiguration(config)
err = ctx.Client.ValidateConfiguration(config)
backup, err := ctx.Client.BackupConfiguration("description")
err = ctx.Client.RestoreConfiguration(backupID)
```

### Alerts
```go
alerts, err := ctx.Client.GetAlerts()
alert, err := ctx.Client.GetAlert(id)
err = ctx.Client.AcknowledgeAlert(id)
err = ctx.Client.ClearAlert(id)
```

---

## Fixture Loaders

```go
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
user := testutil.LoadUserFixture(t, "fixtures/users/admin_user.yaml")
session := testutil.LoadSessionFixture(t, "fixtures/sessions/established_session.yaml")
alert := testutil.LoadAlertFixture(t, "fixtures/alerts/peer_down_alert.yaml")
config := testutil.LoadConfigFixture(t, "fixtures/config/backup_description.yaml")
```

---

## Assertions

### Standard (testify/assert)
```go
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)
assert.Nil(t, object)
assert.NotNil(t, object)
assert.True(t, condition)
assert.False(t, condition)
assert.Contains(t, str, substring)
assert.Empty(t, str)
assert.NotEmpty(t, str)
assert.Greater(t, actual, expected)
assert.Less(t, actual, expected)
assert.Error(t, err)
assert.NoError(t, err)
assert.Len(t, collection, length)
```

### Required (testify/require)
```go
require.NoError(t, err)        // Stops test if fails
require.NotNil(t, object)
require.Equal(t, expected, actual)
```

### Custom
```go
testutil.AssertPeerExists(t, ctx, "peer1")
testutil.AssertSessionState(t, ctx, "peer1", "Established")
testutil.AssertAlertExists(t, ctx, "PeerDown")
```

---

## Database Queries

```go
// Query single row
var count int
err := ctx.DB.QueryRow("SELECT COUNT(*) FROM bgp_peers").Scan(&count)

// Query multiple rows
rows, err := ctx.DB.Query("SELECT name, remote_ip FROM bgp_peers")
defer rows.Close()
for rows.Next() {
    var name, ip string
    rows.Scan(&name, &ip)
}

// Execute statement
result, err := ctx.DB.Exec("DELETE FROM bgp_peers WHERE name = ?", "peer1")

// Helper functions
count, err := testutil.CountRows(ctx, "bgp_peers")
exists, err := testutil.RowExists(ctx, "bgp_peers", "name", "peer1")
err = testutil.ClearTable(ctx, "bgp_peers")
```

---

## Helper Functions

```go
// Wait for condition
err := testutil.WaitForCondition(t, 10*time.Second, func() bool {
    session, _ := ctx.Client.GetSession("peer1")
    return session.State == "Established"
})

// Wait for session state
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)

// Retry operation
err := testutil.RetryOperation(t, 3, func() error {
    return ctx.Client.CreatePeer(peer)
})

// Generate random data
name := "peer-" + testutil.GenerateRandomString(8)
ip := testutil.GenerateRandomIP()
```

---

## Configuration Files

### test-config.yaml
```yaml
server:
  host: 127.0.0.1
  port: 0
database:
  path: ./tmp/test.db
frr:
  grpc_host: localhost
  grpc_port: 50051
auth:
  jwt_secret: test-secret-key
  token_expiry: 5m
testing:
  timeout: 30s
  log_level: debug
```

### mock-frr-config.yaml
```yaml
server:
  host: localhost
  port: 50051
simulation:
  session_state_delay: 100ms
  error_injection: false
logging:
  level: info
  file: ./logs/mock-frr-server.log
```

---

## Mock FRR HTTP API

```bash
# Health check
curl http://localhost:51051/health

# Statistics
curl http://localhost:51051/stats

# List peers
curl http://localhost:51051/peers

# Add peer
curl -X POST http://localhost:51051/peers/add \
  -H "Content-Type: application/json" \
  -d '{"IPAddress":"192.168.1.1","ASN":65000,"RemoteASN":65001}'

# Remove peer
curl -X POST http://localhost:51051/peers/remove \
  -H "Content-Type: application/json" \
  -d '{"ip_address":"192.168.1.1"}'

# Get session state
curl "http://localhost:51051/sessions/state?ip=192.168.1.1"

# Get all sessions
curl http://localhost:51051/sessions

# Get running config
curl http://localhost:51051/config
```

---

## Debugging

### Enable Verbose Logging
```bash
./run-tests.sh --verbose --log-level debug
```

### Preserve Artifacts
```bash
./run-tests.sh --no-cleanup
```

### View Logs
```bash
cat logs/test-execution.log
cat logs/mock-frr-server.log
cat logs/api-server.log
```

### Inspect Database
```bash
sqlite3 tmp/test.db
SELECT * FROM bgp_peers;
SELECT * FROM bgp_sessions;
```

### Check Ports
```bash
lsof -i :8080
lsof -i :50051
lsof -i :51051
```

---

## Common Issues

### Tests Hang
```bash
# Check mock server
ps aux | grep mock-frr-server

# Check ports
lsof -i :50051

# Restart mock server
cd cmd/mock-frr-server && ./mock-frr-server &
```

### Database Locked
```bash
# Clean database
./scripts/cleanup-db.sh

# Check permissions
ls -la tmp/
```

### Port Already in Use
```bash
# Find process
lsof -i :8080

# Kill process
kill -9 <PID>

# Or cleanup all
./scripts/cleanup-all.sh
```

### Authentication Failures
```bash
# Check JWT secret
grep jwt_secret config/test-config.yaml

# Verify user fixtures
cat fixtures/users/admin_user.yaml
```

---

## Test Patterns

### Parameterized Tests
```go
testCases := []struct {
    name     string
    fixture  string
    wantErr  bool
}{
    {"Valid", "fixtures/peers/valid/basic_peer.yaml", false},
    {"Invalid", "fixtures/peers/invalid/bad_ip.yaml", true},
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        peer := testutil.LoadPeerFixture(t, tc.fixture)
        err := ctx.Client.CreatePeer(peer)
        if tc.wantErr {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
        }
    })
}
```

### Async Operations
```go
// Wait for session establishment
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
require.NoError(t, err)
```

### Cleanup Pattern
```go
func TestWithCleanup(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()  // Always defer
    
    // Test code
}
```

---

## Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| 0 | Success | All tests passed |
| 1 | Test failure | Review failures, fix issues |
| 2 | Setup error | Check prerequisites, config |

---

## Useful Links

- **Testing Guide**: [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **API Reference**: [API_REFERENCE.md](API_REFERENCE.md)
- **Mock FRR Server**: [cmd/mock-frr-server/README.md](cmd/mock-frr-server/README.md)
- **Fixtures Guide**: [fixtures/README.md](fixtures/README.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)
- **FAQ**: [FAQ.md](FAQ.md)

---

## Best Practices

✅ **DO**
- Use `defer ctx.Cleanup()` always
- Load fixtures for test data
- Use `require` for critical checks
- Wait for async operations
- Add descriptive assertion messages
- Test both success and failure cases

❌ **DON'T**
- Share state between tests
- Hardcode test data
- Use `time.Sleep()` (use `WaitForCondition`)
- Ignore errors
- Test implementation details
- Make tests too complex

---

**Quick Reference Version**: 1.0.0  
**Last Updated**: November 10, 2025