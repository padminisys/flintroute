# FlintRoute Testing Framework API Reference

**Complete API Documentation for Test Development**

Version: 1.0.0  
Last Updated: November 10, 2025

---

## Table of Contents

1. [Overview](#overview)
2. [Test Context API](#test-context-api)
3. [API Client Methods](#api-client-methods)
4. [Database Utilities](#database-utilities)
5. [Fixture Loaders](#fixture-loaders)
6. [Test Assertions](#test-assertions)
7. [Mock FRR Server API](#mock-frr-server-api)
8. [Helper Functions](#helper-functions)

---

## Overview

This document provides complete API reference for the FlintRoute functional testing framework. All APIs are designed to be used within Go test functions.

### Import Statements

```go
import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "github.com/yourusername/flintroute/test/functional/pkg/testutil"
    "github.com/yourusername/flintroute/test/functional/pkg/client"
    "github.com/yourusername/flintroute/internal/models"
)
```

---

## Test Context API

### NewTestContext

Creates a new test context with all necessary components initialized.

**Signature:**
```go
func NewTestContext(t *testing.T) *TestContext
```

**Parameters:**
- `t` - Testing instance

**Returns:**
- `*TestContext` - Initialized test context

**Example:**
```go
func TestExample(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Use ctx for testing
}
```

### TestContext Structure

```go
type TestContext struct {
    T           *testing.T          // Test instance
    Client      *client.APIClient   // API client
    DB          *sql.DB             // Database connection
    Config      *Config             // Test configuration
    Logger      *logrus.Logger      // Logger instance
    MockFRR     *MockFRRServer      // Mock FRR server reference
    BaseURL     string              // API base URL
    Token       string              // Authentication token
}
```

### TestContext Methods

#### Cleanup

Cleans up all resources created during the test.

**Signature:**
```go
func (ctx *TestContext) Cleanup()
```

**Example:**
```go
func TestExample(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()  // Always defer cleanup
}
```

#### SetTimeout

Sets a timeout for the test context.

**Signature:**
```go
func (ctx *TestContext) SetTimeout(duration time.Duration)
```

**Parameters:**
- `duration` - Timeout duration

**Example:**
```go
ctx.SetTimeout(30 * time.Second)
```

#### GetAuthToken

Retrieves the current authentication token.

**Signature:**
```go
func (ctx *TestContext) GetAuthToken() string
```

**Returns:**
- `string` - JWT authentication token

**Example:**
```go
token := ctx.GetAuthToken()
```

#### SetAuthToken

Sets the authentication token for subsequent requests.

**Signature:**
```go
func (ctx *TestContext) SetAuthToken(token string)
```

**Parameters:**
- `token` - JWT authentication token

**Example:**
```go
ctx.SetAuthToken("eyJhbGciOiJIUzI1NiIs...")
```

---

## API Client Methods

### Authentication

#### Login

Authenticates a user and returns a JWT token.

**Signature:**
```go
func (c *APIClient) Login(username, password string) (string, error)
```

**Parameters:**
- `username` - User's username
- `password` - User's password

**Returns:**
- `string` - JWT token
- `error` - Error if authentication fails

**Example:**
```go
token, err := ctx.Client.Login("admin", "admin123")
require.NoError(t, err)
assert.NotEmpty(t, token)
```

#### Logout

Logs out the current user.

**Signature:**
```go
func (c *APIClient) Logout() error
```

**Returns:**
- `error` - Error if logout fails

**Example:**
```go
err := ctx.Client.Logout()
assert.NoError(t, err)
```

#### RefreshToken

Refreshes the authentication token.

**Signature:**
```go
func (c *APIClient) RefreshToken() (string, error)
```

**Returns:**
- `string` - New JWT token
- `error` - Error if refresh fails

**Example:**
```go
newToken, err := ctx.Client.RefreshToken()
require.NoError(t, err)
assert.NotEmpty(t, newToken)
```

#### IsAuthenticated

Checks if the client is currently authenticated.

**Signature:**
```go
func (c *APIClient) IsAuthenticated() bool
```

**Returns:**
- `bool` - True if authenticated

**Example:**
```go
if ctx.Client.IsAuthenticated() {
    // Perform authenticated operations
}
```

### Peer Management

#### CreatePeer

Creates a new BGP peer.

**Signature:**
```go
func (c *APIClient) CreatePeer(peer *models.BGPPeer) error
```

**Parameters:**
- `peer` - BGP peer configuration

**Returns:**
- `error` - Error if creation fails

**Example:**
```go
peer := &models.BGPPeer{
    Name:      "peer1",
    RemoteIP:  "192.168.1.1",
    RemoteASN: 65001,
}
err := ctx.Client.CreatePeer(peer)
require.NoError(t, err)
```

#### GetPeer

Retrieves a BGP peer by name.

**Signature:**
```go
func (c *APIClient) GetPeer(name string) (*models.BGPPeer, error)
```

**Parameters:**
- `name` - Peer name

**Returns:**
- `*models.BGPPeer` - Peer configuration
- `error` - Error if retrieval fails

**Example:**
```go
peer, err := ctx.Client.GetPeer("peer1")
require.NoError(t, err)
assert.Equal(t, "192.168.1.1", peer.RemoteIP)
```

#### ListPeers

Lists all BGP peers.

**Signature:**
```go
func (c *APIClient) ListPeers() ([]*models.BGPPeer, error)
```

**Returns:**
- `[]*models.BGPPeer` - List of peers
- `error` - Error if listing fails

**Example:**
```go
peers, err := ctx.Client.ListPeers()
require.NoError(t, err)
assert.Greater(t, len(peers), 0)
```

#### UpdatePeer

Updates an existing BGP peer.

**Signature:**
```go
func (c *APIClient) UpdatePeer(peer *models.BGPPeer) error
```

**Parameters:**
- `peer` - Updated peer configuration

**Returns:**
- `error` - Error if update fails

**Example:**
```go
peer.Description = "Updated description"
err := ctx.Client.UpdatePeer(peer)
require.NoError(t, err)
```

#### DeletePeer

Deletes a BGP peer.

**Signature:**
```go
func (c *APIClient) DeletePeer(name string) error
```

**Parameters:**
- `name` - Peer name

**Returns:**
- `error` - Error if deletion fails

**Example:**
```go
err := ctx.Client.DeletePeer("peer1")
require.NoError(t, err)
```

#### EnablePeer

Enables a BGP peer.

**Signature:**
```go
func (c *APIClient) EnablePeer(name string) error
```

**Parameters:**
- `name` - Peer name

**Returns:**
- `error` - Error if operation fails

**Example:**
```go
err := ctx.Client.EnablePeer("peer1")
require.NoError(t, err)
```

#### DisablePeer

Disables a BGP peer.

**Signature:**
```go
func (c *APIClient) DisablePeer(name string) error
```

**Parameters:**
- `name` - Peer name

**Returns:**
- `error` - Error if operation fails

**Example:**
```go
err := ctx.Client.DisablePeer("peer1")
require.NoError(t, err)
```

### Session Management

#### GetSession

Retrieves BGP session state for a peer.

**Signature:**
```go
func (c *APIClient) GetSession(peerName string) (*models.BGPSession, error)
```

**Parameters:**
- `peerName` - Peer name

**Returns:**
- `*models.BGPSession` - Session state
- `error` - Error if retrieval fails

**Example:**
```go
session, err := ctx.Client.GetSession("peer1")
require.NoError(t, err)
assert.Equal(t, "Established", session.State)
```

#### ListSessions

Lists all BGP sessions.

**Signature:**
```go
func (c *APIClient) ListSessions() ([]*models.BGPSession, error)
```

**Returns:**
- `[]*models.BGPSession` - List of sessions
- `error` - Error if listing fails

**Example:**
```go
sessions, err := ctx.Client.ListSessions()
require.NoError(t, err)
assert.Greater(t, len(sessions), 0)
```

#### ResetSession

Resets a BGP session (soft or hard reset).

**Signature:**
```go
func (c *APIClient) ResetSession(peerName string, hard bool) error
```

**Parameters:**
- `peerName` - Peer name
- `hard` - True for hard reset, false for soft reset

**Returns:**
- `error` - Error if reset fails

**Example:**
```go
// Soft reset
err := ctx.Client.ResetSession("peer1", false)
require.NoError(t, err)

// Hard reset
err = ctx.Client.ResetSession("peer1", true)
require.NoError(t, err)
```

### Configuration Management

#### GetConfiguration

Retrieves the current system configuration.

**Signature:**
```go
func (c *APIClient) GetConfiguration() (*models.Configuration, error)
```

**Returns:**
- `*models.Configuration` - System configuration
- `error` - Error if retrieval fails

**Example:**
```go
config, err := ctx.Client.GetConfiguration()
require.NoError(t, err)
assert.NotNil(t, config)
```

#### UpdateConfiguration

Updates the system configuration.

**Signature:**
```go
func (c *APIClient) UpdateConfiguration(config *models.Configuration) error
```

**Parameters:**
- `config` - New configuration

**Returns:**
- `error` - Error if update fails

**Example:**
```go
config.LocalASN = 65000
err := ctx.Client.UpdateConfiguration(config)
require.NoError(t, err)
```

#### ValidateConfiguration

Validates a configuration without applying it.

**Signature:**
```go
func (c *APIClient) ValidateConfiguration(config *models.Configuration) error
```

**Parameters:**
- `config` - Configuration to validate

**Returns:**
- `error` - Error if validation fails

**Example:**
```go
err := ctx.Client.ValidateConfiguration(config)
assert.NoError(t, err)
```

#### BackupConfiguration

Creates a backup of the current configuration.

**Signature:**
```go
func (c *APIClient) BackupConfiguration(description string) (*models.ConfigBackup, error)
```

**Parameters:**
- `description` - Backup description

**Returns:**
- `*models.ConfigBackup` - Backup information
- `error` - Error if backup fails

**Example:**
```go
backup, err := ctx.Client.BackupConfiguration("Pre-upgrade backup")
require.NoError(t, err)
assert.NotEmpty(t, backup.ID)
```

#### RestoreConfiguration

Restores a configuration from backup.

**Signature:**
```go
func (c *APIClient) RestoreConfiguration(backupID string) error
```

**Parameters:**
- `backupID` - Backup identifier

**Returns:**
- `error` - Error if restore fails

**Example:**
```go
err := ctx.Client.RestoreConfiguration(backup.ID)
require.NoError(t, err)
```

### Alert Management

#### GetAlerts

Retrieves all alerts.

**Signature:**
```go
func (c *APIClient) GetAlerts() ([]*models.Alert, error)
```

**Returns:**
- `[]*models.Alert` - List of alerts
- `error` - Error if retrieval fails

**Example:**
```go
alerts, err := ctx.Client.GetAlerts()
require.NoError(t, err)
assert.Greater(t, len(alerts), 0)
```

#### GetAlert

Retrieves a specific alert by ID.

**Signature:**
```go
func (c *APIClient) GetAlert(id string) (*models.Alert, error)
```

**Parameters:**
- `id` - Alert identifier

**Returns:**
- `*models.Alert` - Alert details
- `error` - Error if retrieval fails

**Example:**
```go
alert, err := ctx.Client.GetAlert("alert-123")
require.NoError(t, err)
assert.Equal(t, "PeerDown", alert.Type)
```

#### AcknowledgeAlert

Acknowledges an alert.

**Signature:**
```go
func (c *APIClient) AcknowledgeAlert(id string) error
```

**Parameters:**
- `id` - Alert identifier

**Returns:**
- `error` - Error if acknowledgment fails

**Example:**
```go
err := ctx.Client.AcknowledgeAlert("alert-123")
require.NoError(t, err)
```

#### ClearAlert

Clears an alert.

**Signature:**
```go
func (c *APIClient) ClearAlert(id string) error
```

**Parameters:**
- `id` - Alert identifier

**Returns:**
- `error` - Error if clearing fails

**Example:**
```go
err := ctx.Client.ClearAlert("alert-123")
require.NoError(t, err)
```

---

## Database Utilities

### QueryRow

Executes a query that returns a single row.

**Signature:**
```go
func (ctx *TestContext) QueryRow(query string, args ...interface{}) *sql.Row
```

**Parameters:**
- `query` - SQL query
- `args` - Query arguments

**Returns:**
- `*sql.Row` - Query result

**Example:**
```go
var count int
err := ctx.DB.QueryRow("SELECT COUNT(*) FROM bgp_peers").Scan(&count)
require.NoError(t, err)
assert.Greater(t, count, 0)
```

### Query

Executes a query that returns multiple rows.

**Signature:**
```go
func (ctx *TestContext) Query(query string, args ...interface{}) (*sql.Rows, error)
```

**Parameters:**
- `query` - SQL query
- `args` - Query arguments

**Returns:**
- `*sql.Rows` - Query results
- `error` - Error if query fails

**Example:**
```go
rows, err := ctx.DB.Query("SELECT name, remote_ip FROM bgp_peers")
require.NoError(t, err)
defer rows.Close()

for rows.Next() {
    var name, ip string
    err := rows.Scan(&name, &ip)
    require.NoError(t, err)
    t.Logf("Peer: %s, IP: %s", name, ip)
}
```

### Exec

Executes a query that doesn't return rows.

**Signature:**
```go
func (ctx *TestContext) Exec(query string, args ...interface{}) (sql.Result, error)
```

**Parameters:**
- `query` - SQL query
- `args` - Query arguments

**Returns:**
- `sql.Result` - Execution result
- `error` - Error if execution fails

**Example:**
```go
result, err := ctx.DB.Exec("DELETE FROM bgp_peers WHERE name = ?", "test-peer")
require.NoError(t, err)

rowsAffected, err := result.RowsAffected()
require.NoError(t, err)
assert.Equal(t, int64(1), rowsAffected)
```

### Database Helper Functions

#### CountRows

Counts rows in a table.

**Signature:**
```go
func CountRows(ctx *TestContext, table string) (int, error)
```

**Parameters:**
- `ctx` - Test context
- `table` - Table name

**Returns:**
- `int` - Row count
- `error` - Error if count fails

**Example:**
```go
count, err := testutil.CountRows(ctx, "bgp_peers")
require.NoError(t, err)
assert.Equal(t, 5, count)
```

#### RowExists

Checks if a row exists.

**Signature:**
```go
func RowExists(ctx *TestContext, table, column, value string) (bool, error)
```

**Parameters:**
- `ctx` - Test context
- `table` - Table name
- `column` - Column name
- `value` - Value to check

**Returns:**
- `bool` - True if row exists
- `error` - Error if check fails

**Example:**
```go
exists, err := testutil.RowExists(ctx, "bgp_peers", "name", "peer1")
require.NoError(t, err)
assert.True(t, exists)
```

#### ClearTable

Clears all rows from a table.

**Signature:**
```go
func ClearTable(ctx *TestContext, table string) error
```

**Parameters:**
- `ctx` - Test context
- `table` - Table name

**Returns:**
- `error` - Error if clear fails

**Example:**
```go
err := testutil.ClearTable(ctx, "bgp_peers")
require.NoError(t, err)
```

---

## Fixture Loaders

### LoadPeerFixture

Loads a BGP peer fixture from a YAML file.

**Signature:**
```go
func LoadPeerFixture(t *testing.T, path string) *models.BGPPeer
```

**Parameters:**
- `t` - Testing instance
- `path` - Fixture file path

**Returns:**
- `*models.BGPPeer` - Loaded peer configuration

**Example:**
```go
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
assert.NotNil(t, peer)
assert.Equal(t, "peer1", peer.Name)
```

### LoadUserFixture

Loads a user fixture from a YAML file.

**Signature:**
```go
func LoadUserFixture(t *testing.T, path string) *models.User
```

**Parameters:**
- `t` - Testing instance
- `path` - Fixture file path

**Returns:**
- `*models.User` - Loaded user

**Example:**
```go
user := testutil.LoadUserFixture(t, "fixtures/users/admin_user.yaml")
assert.NotNil(t, user)
assert.Equal(t, "admin", user.Username)
```

### LoadSessionFixture

Loads a BGP session fixture from a YAML file.

**Signature:**
```go
func LoadSessionFixture(t *testing.T, path string) *models.BGPSession
```

**Parameters:**
- `t` - Testing instance
- `path` - Fixture file path

**Returns:**
- `*models.BGPSession` - Loaded session state

**Example:**
```go
session := testutil.LoadSessionFixture(t, "fixtures/sessions/established_session.yaml")
assert.NotNil(t, session)
assert.Equal(t, "Established", session.State)
```

### LoadAlertFixture

Loads an alert fixture from a YAML file.

**Signature:**
```go
func LoadAlertFixture(t *testing.T, path string) *models.Alert
```

**Parameters:**
- `t` - Testing instance
- `path` - Fixture file path

**Returns:**
- `*models.Alert` - Loaded alert

**Example:**
```go
alert := testutil.LoadAlertFixture(t, "fixtures/alerts/peer_down_alert.yaml")
assert.NotNil(t, alert)
assert.Equal(t, "PeerDown", alert.Type)
```

### LoadConfigFixture

Loads a configuration fixture from a YAML file.

**Signature:**
```go
func LoadConfigFixture(t *testing.T, path string) *models.Configuration
```

**Parameters:**
- `t` - Testing instance
- `path` - Fixture file path

**Returns:**
- `*models.Configuration` - Loaded configuration

**Example:**
```go
config := testutil.LoadConfigFixture(t, "fixtures/config/backup_description.yaml")
assert.NotNil(t, config)
```

---

## Test Assertions

### Standard Assertions

Using `github.com/stretchr/testify/assert`:

```go
// Equality
assert.Equal(t, expected, actual)
assert.NotEqual(t, expected, actual)

// Nil checks
assert.Nil(t, object)
assert.NotNil(t, object)

// Boolean
assert.True(t, condition)
assert.False(t, condition)

// Strings
assert.Contains(t, str, substring)
assert.NotContains(t, str, substring)
assert.Empty(t, str)
assert.NotEmpty(t, str)

// Numbers
assert.Greater(t, actual, expected)
assert.GreaterOrEqual(t, actual, expected)
assert.Less(t, actual, expected)
assert.LessOrEqual(t, actual, expected)

// Errors
assert.Error(t, err)
assert.NoError(t, err)
assert.EqualError(t, err, "expected error message")

// Collections
assert.Len(t, collection, expectedLength)
assert.ElementsMatch(t, expected, actual)
```

### Required Assertions

Using `github.com/stretchr/testify/require` (stops test on failure):

```go
// Use require for critical checks
require.NoError(t, err)
require.NotNil(t, object)
require.Equal(t, expected, actual)
```

### Custom Assertions

#### AssertPeerExists

Asserts that a peer exists in the database.

**Signature:**
```go
func AssertPeerExists(t *testing.T, ctx *TestContext, peerName string)
```

**Parameters:**
- `t` - Testing instance
- `ctx` - Test context
- `peerName` - Peer name

**Example:**
```go
testutil.AssertPeerExists(t, ctx, "peer1")
```

#### AssertSessionState

Asserts that a session is in a specific state.

**Signature:**
```go
func AssertSessionState(t *testing.T, ctx *TestContext, peerName, expectedState string)
```

**Parameters:**
- `t` - Testing instance
- `ctx` - Test context
- `peerName` - Peer name
- `expectedState` - Expected BGP state

**Example:**
```go
testutil.AssertSessionState(t, ctx, "peer1", "Established")
```

#### AssertAlertExists

Asserts that an alert exists.

**Signature:**
```go
func AssertAlertExists(t *testing.T, ctx *TestContext, alertType string)
```

**Parameters:**
- `t` - Testing instance
- `ctx` - Test context
- `alertType` - Alert type

**Example:**
```go
testutil.AssertAlertExists(t, ctx, "PeerDown")
```

---

## Mock FRR Server API

### HTTP Debug Interface

The Mock FRR server provides an HTTP interface for testing and debugging.

**Base URL**: `http://localhost:51051`

#### Health Check

```bash
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "uptime": "1h30m"
}
```

#### Statistics

```bash
GET /stats
```

**Response:**
```json
{
  "peers_count": 5,
  "sessions_established": 3,
  "total_requests": 150
}
```

#### List Peers

```bash
GET /peers
```

**Response:**
```json
[
  {
    "ip_address": "192.168.1.1",
    "asn": 65001,
    "remote_asn": 65000
  }
]
```

#### Add Peer

```bash
POST /peers/add
Content-Type: application/json

{
  "IPAddress": "192.168.1.1",
  "ASN": 65000,
  "RemoteASN": 65001
}
```

#### Remove Peer

```bash
POST /peers/remove
Content-Type: application/json

{
  "ip_address": "192.168.1.1"
}
```

#### Get Session State

```bash
GET /sessions/state?ip=192.168.1.1
```

**Response:**
```json
{
  "ip_address": "192.168.1.1",
  "state": "Established",
  "uptime": 3600,
  "prefixes_received": 100,
  "prefixes_sent": 50
}
```

#### Get All Sessions

```bash
GET /sessions
```

**Response:**
```json
[
  {
    "ip_address": "192.168.1.1",
    "state": "Established",
    "uptime": 3600
  }
]
```

---

## Helper Functions

### WaitForCondition

Waits for a condition to become true with timeout.

**Signature:**
```go
func WaitForCondition(t *testing.T, timeout time.Duration, condition func() bool) error
```

**Parameters:**
- `t` - Testing instance
- `timeout` - Maximum wait time
- `condition` - Function that returns true when condition is met

**Returns:**
- `error` - Error if timeout occurs

**Example:**
```go
err := testutil.WaitForCondition(t, 10*time.Second, func() bool {
    session, _ := ctx.Client.GetSession("peer1")
    return session.State == "Established"
})
require.NoError(t, err)
```

### WaitForSessionState

Waits for a BGP session to reach a specific state.

**Signature:**
```go
func WaitForSessionState(t *testing.T, ctx *TestContext, peerName, state string, timeout time.Duration) error
```

**Parameters:**
- `t` - Testing instance
- `ctx` - Test context
- `peerName` - Peer name
- `state` - Expected BGP state
- `timeout` - Maximum wait time

**Returns:**
- `error` - Error if timeout occurs

**Example:**
```go
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
require.NoError(t, err)
```

### RetryOperation

Retries an operation with exponential backoff.

**Signature:**
```go
func RetryOperation(t *testing.T, maxRetries int, operation func() error) error
```

**Parameters:**
- `t` - Testing instance
- `maxRetries` - Maximum number of retries
- `operation` - Function to retry

**Returns:**
- `error` - Error if all retries fail

**Example:**
```go
err := testutil.RetryOperation(t, 3, func() error {
    return ctx.Client.CreatePeer(peer)
})
require.NoError(t, err)
```

### GenerateRandomString

Generates a random string for unique test data.

**Signature:**
```go
func GenerateRandomString(length int) string
```

**Parameters:**
- `length` - String length

**Returns:**
- `string` - Random string

**Example:**
```go
uniqueName := "peer-" + testutil.GenerateRandomString(8)
```

### GenerateRandomIP

Generates a random IP address for testing.

**Signature:**
```go
func GenerateRandomIP() string
```

**Returns:**
- `string` - Random IP address (RFC 1918 range)

**Example:**
```go
ip := testutil.GenerateRandomIP()
// Returns something like "192.168.100.42"
```

### Sleep

Sleeps for a specified duration (use sparingly, prefer WaitForCondition).

**Signature:**
```go
func Sleep(duration time.Duration)
```

**Parameters:**
- `duration` - Sleep duration

**Example:**
```go
testutil.Sleep(1 * time.Second)
```

---

## Error Handling

### Error Types

```go
// ErrNotFound - Resource not found
var ErrNotFound = errors.New("resource not found")

// ErrUnauthorized - Authentication required
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden - Insufficient permissions
var ErrForbidden = errors.New("forbidden")

// ErrBadRequest - Invalid request
var ErrBadRequest = errors.New("bad request")

// ErrConflict - Resource conflict
var ErrConflict = errors.New("conflict")

// ErrTimeout - Operation timeout
var ErrTimeout = errors.New("timeout")
```

### Error Checking

```go
// Check for specific error
if errors.Is(err, testutil.ErrNotFound) {
    // Handle not found
}

// Check error message
if err != nil && strings.Contains(err.Error(), "already exists") {
    // Handle duplicate
}
```

---

## Configuration

### Test Configuration Structure

```go
type Config struct {
    Server struct {
        Host string
        Port int
    }
    Database struct {
        Path string
    }
    FRR struct {
        GRPCHost string
        GRPCPort int
    }
    Auth struct {
        JWTSecret     string
        TokenExpiry   time.Duration
        RefreshExpiry time.Duration
    }
    Testing struct {
        Timeout          time.Duration
        CleanupOnSuccess bool
        LogLevel         string
        Parallel         bool
    }
}
```

### Loading Configuration

```go
config, err := testutil.LoadConfig("config/test-config.yaml")
require.NoError(t, err)
```

---

## Best Practices

1. **Always use defer for cleanup**
   ```go
   ctx := testutil.NewTestContext(t)
   defer ctx.Cleanup()
   ```

2. **Use require for critical checks**
   ```go
   require.NoError(t, err)  // Stops test if fails
   assert.Equal(t, a, b)    // Continues test if fails
   ```

3. **Load fixtures for test data**
   ```go
   peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
   ```

4. **Wait for async operations**
   ```go
   err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
   ```

5. **Use descriptive assertions**
   ```go
   assert.Equal(t, "Established", session.State, "Session should be established after 5 seconds")
   ```

---

**Document Version**: 1.0.0  
**Last Updated**: November 10, 2025  
**Maintained By**: FlintRoute Development Team