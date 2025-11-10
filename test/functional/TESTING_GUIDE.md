
# FlintRoute Functional Testing Guide

**Complete Guide to Testing FlintRoute's BGP Management System**

Version: 1.0.0  
Last Updated: November 10, 2025

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Quick Start](#quick-start)
4. [Prerequisites](#prerequisites)
5. [Running Tests](#running-tests)
6. [Writing Tests](#writing-tests)
7. [Fixtures](#fixtures)
8. [Debugging](#debugging)
9. [CI/CD Integration](#cicd-integration)
10. [Best Practices](#best-practices)

---

## Overview

### What is the Functional Testing Framework?

The FlintRoute Functional Testing Framework is a comprehensive end-to-end testing solution designed to validate the complete behavior of the FlintRoute BGP management system. Unlike unit tests that verify individual components in isolation, functional tests validate the entire system working together, including:

- **API Endpoints** - REST API request/response validation
- **Authentication & Authorization** - User login, token management, role-based access
- **BGP Operations** - Peer management, session state tracking
- **Database Operations** - Data persistence and retrieval
- **FRR Integration** - Communication with FRRouting via gRPC
- **WebSocket Communication** - Real-time updates and notifications
- **Error Handling** - Graceful degradation and error recovery

### Key Features

✅ **Realistic Testing Environment** - Mock FRR server simulates production behavior  
✅ **Comprehensive Coverage** - 7 test suites covering all major features  
✅ **Isolated Test Execution** - Each test runs in a clean environment  
✅ **Flexible Configuration** - YAML-based configuration for different scenarios  
✅ **Detailed Reporting** - JSON, HTML, and text reports  
✅ **CI/CD Ready** - Automated execution in continuous integration pipelines  
✅ **Easy Debugging** - Verbose logging and artifact preservation on failure

### Test Coverage

The framework includes **7 comprehensive test suites**:

1. **Authentication** (01_authentication/) - User login, JWT tokens, permissions
2. **Peer Management** (02_peer_management/) - CRUD operations for BGP peers
3. **Session Management** (03_session_management/) - BGP session state tracking
4. **Configuration** (04_configuration/) - System configuration and validation
5. **Alerts** (05_alerts/) - Alert generation and notification handling
6. **Error Handling** (06_error_handling/) - Error scenarios and recovery
7. **Workflows** (07_workflows/) - Complete end-to-end user workflows

---

## Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Test Runner (run-tests.sh)                │
│  • Environment Setup  • Test Execution  • Report Generation  │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
┌────────────────┐ ┌──────────┐ ┌──────────────┐
│  Test Suites   │ │ Fixtures │ │   Scripts    │
│  (Go Tests)    │ │  (YAML)  │ │   (Bash)     │
└────────┬───────┘ └────┬─────┘ └──────┬───────┘
         │              │               │
         └──────────────┼───────────────┘
                        │
         ┌──────────────┴──────────────┐
         │                             │
         ▼                             ▼
┌─────────────────┐          ┌──────────────────┐
│  FlintRoute API │          │  Mock FRR Server │
│  (Under Test)   │◄────────►│  (gRPC/HTTP)     │
└─────────────────┘          └──────────────────┘
         │
         ▼
┌─────────────────┐
│  SQLite DB      │
│  (Test Data)    │
└─────────────────┘
```

### Component Overview

#### 1. Test Runner (`run-tests.sh`)
- **Purpose**: Orchestrates the entire test lifecycle
- **Responsibilities**:
  - Environment setup and teardown
  - Test execution with configurable parameters
  - Report generation (JSON, HTML, text)
  - Cleanup on success/failure
- **Location**: [`test/functional/run-tests.sh`](run-tests.sh:1)

#### 2. Mock FRR Server
- **Purpose**: Simulates FRRouting's gRPC interface
- **Features**:
  - Realistic BGP session state transitions
  - Configurable delays and error injection
  - Dual interface (gRPC + HTTP debug)
  - Thread-safe concurrent operations
- **Location**: [`test/functional/cmd/mock-frr-server/`](cmd/mock-frr-server/)
- **Documentation**: [Mock FRR Server README](cmd/mock-frr-server/README.md:1)

#### 3. Test Suites
- **Purpose**: Validate specific features and workflows
- **Organization**: Numbered for sequential execution
- **Location**: [`test/functional/tests/`](tests/)
- **Language**: Go with testify assertions

#### 4. Fixtures
- **Purpose**: Provide consistent, reusable test data
- **Formats**: YAML (primary), JSON (legacy)
- **Categories**: Peers, users, sessions, alerts, configurations
- **Location**: [`test/functional/fixtures/`](fixtures/)
- **Documentation**: [Fixtures README](fixtures/README.md:1)

#### 5. Configuration
- **Purpose**: Control test behavior and environment
- **Files**:
  - [`test-config.yaml`](config/test-config.yaml:1) - Main test configuration
  - [`mock-frr-config.yaml`](config/mock-frr-config.yaml:1) - Mock server settings
  - [`logging-config.yaml`](config/logging-config.yaml:1) - Logging configuration
- **Location**: [`test/functional/config/`](config/)

#### 6. Scripts
- **Purpose**: Automate common tasks
- **Examples**:
  - [`check-prerequisites.sh`](scripts/check-prerequisites.sh:1) - Verify system requirements
  - [`setup-env.sh`](scripts/setup-env.sh:1) - Initialize test environment
  - [`cleanup-all.sh`](scripts/cleanup-all.sh:1) - Clean all artifacts
- **Location**: [`test/functional/scripts/`](scripts/)

### Data Flow

```
1. Test Initialization
   ├─ Load configuration (test-config.yaml)
   ├─ Start Mock FRR Server (port 50051)
   ├─ Initialize test database (tmp/test.db)
   └─ Create API client

2. Test Execution
   ├─ Load fixtures (YAML data)
   ├─ Make API requests
   ├─ Verify responses
   ├─ Check database state
   └─ Validate FRR interactions

3. Test Cleanup
   ├─ Stop Mock FRR Server
   ├─ Clean database
   ├─ Archive logs
   └─ Generate reports
```

---

## Quick Start

### Get Started in 5 Minutes

```bash
# 1. Navigate to functional test directory
cd test/functional

# 2. Check prerequisites
./scripts/check-prerequisites.sh

# 3. Run all tests (clean environment)
./run-clean.sh

# 4. View results
cat results/test-summary.txt
```

### Your First Test Run

```bash
# Run a specific test suite
./run-tests.sh --pattern ./tests/01_authentication/...

# Run with verbose output
./run-tests.sh --verbose

# Run with debug logging
./run-tests.sh --log-level debug

# Run without cleanup (for debugging)
./run-tests.sh --no-cleanup
```

### Understanding Test Output

```
[INFO] FlintRoute Functional Test Runner
[INFO] ==================================
[INFO] Using config file: config/test-config.yaml
[INFO] Setting up test environment...
[✓] Environment setup completed
[INFO] Running tests with pattern: ./...
[INFO] Log level: info
=== RUN   TestAuthentication
=== RUN   TestAuthentication/Login_Success
--- PASS: TestAuthentication/Login_Success (0.05s)
=== RUN   TestAuthentication/Login_Invalid_Credentials
--- PASS: TestAuthentication/Login_Invalid_Credentials (0.03s)
--- PASS: TestAuthentication (0.08s)
[SUCCESS] All tests passed
```

---

## Prerequisites

### Required Software

| Software | Minimum Version | Purpose |
|----------|----------------|---------|
| **Go** | 1.21+ | Test execution and compilation |
| **Make** | Any | Build automation |
| **SQLite3** | 3.x | Test database |
| **Git** | 2.x | Version control |

### Optional Tools

| Tool | Purpose |
|------|---------|
| **curl** | Manual API testing |
| **jq** | JSON processing |
| **protoc** | Rebuilding proto files |
| **lsof** | Port checking |

### System Requirements

- **Disk Space**: 100MB minimum
- **Memory**: 512MB minimum
- **Ports**: 8080, 50051, 51051 must be available
- **OS**: Linux, macOS, or WSL2 on Windows

### Checking Prerequisites

Run the prerequisite checker to verify your system:

```bash
cd test/functional
./scripts/check-prerequisites.sh
```

**Expected Output:**
```
[INFO] FlintRoute Test Prerequisites Check
[INFO] ====================================

[✓] Go version 1.21.5 (>= 1.21 required)
[✓] Port 8080 is available
[✓] Port 50051 is available
[✓] Port 51051 is available
[✓] Disk space: 5000MB available (>= 100MB required)
[✓] Write permission: ./logs
[✓] Write permission: ./results
[✓] Write permission: ./tmp
[✓] Required command: make
[✓] Required command: git
[SUCCESS] All prerequisites met
```

### Installation Guide

#### Installing Go (if needed)

**Linux:**
```bash
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**macOS:**
```bash
brew install go@1.21
```

#### Installing SQLite3

**Linux:**
```bash
sudo apt-get install sqlite3  # Debian/Ubuntu
sudo yum install sqlite       # RHEL/CentOS
```

**macOS:**
```bash
brew install sqlite3
```

---

## Running Tests

### Basic Test Execution

#### Run All Tests
```bash
./run-tests.sh
```

#### Run Specific Test Suite
```bash
./run-tests.sh --pattern ./tests/01_authentication/...
```

#### Run Single Test
```bash
cd tests/01_authentication
go test -v -run TestLogin
```

### Advanced Options

#### Clean Run (Recommended)
Performs full cleanup before running tests:
```bash
./run-clean.sh
```

#### Verbose Output
```bash
./run-tests.sh --verbose
```

#### Debug Logging
```bash
./run-tests.sh --log-level debug
```

#### Keep Artifacts on Failure
```bash
./run-tests.sh --no-cleanup
```

#### Custom Configuration
```bash
./run-tests.sh --config config/custom-config.yaml
```

### Test Runner Options

```
Usage: ./run-tests.sh [OPTIONS]

OPTIONS:
    --pattern PATTERN    Run tests matching pattern (default: ./...)
    --config FILE        Use specific config file (default: config/test-config.yaml)
    --log-level LEVEL    Set log level: debug|info|warn|error (default: info)
    --no-cleanup         Don't cleanup on success
    --verbose            Verbose output
    --help               Show help message

EXAMPLES:
    # Run all tests
    ./run-tests.sh

    # Run authentication tests only
    ./run-tests.sh --pattern ./tests/01_authentication/...

    # Run with debug logging
    ./run-tests.sh --log-level debug --verbose

    # Run without cleanup (for debugging)
    ./run-tests.sh --no-cleanup
```

### Test Execution Flow

```
1. Pre-Flight Checks
   ├─ Validate configuration file
   ├─ Check port availability
   └─ Verify write permissions

2. Environment Setup
   ├─ Create directories (logs/, results/, tmp/)
   ├─ Initialize test database
   ├─ Start Mock FRR Server
   └─ Wait for services to be ready

3. Test Execution
   ├─ Run tests with go test
   ├─ Capture JSON output
   ├─ Stream logs to console
   └─ Handle test failures

4. Report Generation
   ├─ Parse JSON results
   ├─ Generate HTML report (if tool available)
   ├─ Create text summary
   └─ Archive artifacts

5. Cleanup
   ├─ Stop Mock FRR Server
   ├─ Clean temporary files (if successful)
   ├─ Archive logs
   └─ Exit with appropriate code
```

### Understanding Exit Codes

| Exit Code | Meaning | Action |
|-----------|---------|--------|
| **0** | All tests passed | Success - continue |
| **1** | Tests failed | Review failures, fix issues |
| **2** | Setup/environment error | Check prerequisites, configuration |

---

## Writing Tests

### Test Structure

Every test follows this standard structure:

```go
package authentication_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFeature(t *testing.T) {
    // 1. Setup - Initialize test context
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // 2. Arrange - Prepare test data
    user := testutil.LoadUserFixture(t, "fixtures/users/admin_user.yaml")
    
    // 3. Act - Execute the operation
    token, err := ctx.Client.Login(user.Username, user.Password)
    
    // 4. Assert - Verify results
    require.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.True(t, ctx.Client.IsAuthenticated())
}
```

### Test Organization

#### Directory Structure
```
tests/
├── 01_authentication/
│   ├── login_test.go
│   ├── logout_test.go
│   └── permissions_test.go
├── 02_peer_management/
│   ├── create_peer_test.go
│   ├── update_peer_test.go
│   └── delete_peer_test.go
└── ...
```

#### Naming Conventions

- **Files**: `feature_test.go` (e.g., `login_test.go`)
- **Functions**: `TestFeatureName` (e.g., `TestLoginSuccess`)
- **Subtests**: Use `t.Run()` for scenarios

```go
func TestPeerManagement(t *testing.T) {
    t.Run("Create_Valid_Peer", func(t *testing.T) {
        // Test creating a valid peer
    })
    
    t.Run("Create_Invalid_IP", func(t *testing.T) {
        // Test error handling for invalid IP
    })
    
    t.Run("Update_Peer_Configuration", func(t *testing.T) {
        // Test updating peer settings
    })
}
```

### Using Test Context

The test context provides a complete testing environment:

```go
// Create test context
ctx := testutil.NewTestContext(t)
defer ctx.Cleanup()

// Access components
ctx.Client      // API client
ctx.DB          // Database connection
ctx.Config      // Test configuration
ctx.Logger      // Test logger
ctx.MockFRR     // Mock FRR server reference
```

### Making API Requests

```go
// Authentication
token, err := ctx.Client.Login("admin", "password")
require.NoError(t, err)

// Create peer
peer := &models.BGPPeer{
    Name:      "test-peer",
    RemoteIP:  "192.168.1.1",
    RemoteASN: 65001,
}
err = ctx.Client.CreatePeer(peer)
require.NoError(t, err)

// Get peer
retrieved, err := ctx.Client.GetPeer(peer.Name)
require.NoError(t, err)
assert.Equal(t, peer.RemoteIP, retrieved.RemoteIP)

// Update peer
peer.Description = "Updated description"
err = ctx.Client.UpdatePeer(peer)
require.NoError(t, err)

// Delete peer
err = ctx.Client.DeletePeer(peer.Name)
require.NoError(t, err)
```

### Database Assertions

```go
// Verify data in database
var count int
err := ctx.DB.QueryRow("SELECT COUNT(*) FROM bgp_peers WHERE name = ?", "test-peer").Scan(&count)
require.NoError(t, err)
assert.Equal(t, 1, count)

// Verify peer state
var state string
err = ctx.DB.QueryRow("SELECT state FROM bgp_sessions WHERE peer_name = ?", "test-peer").Scan(&state)
require.NoError(t, err)
assert.Equal(t, "Established", state)
```

### Error Testing

```go
func TestInvalidPeerCreation(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Test with invalid IP
    peer := &models.BGPPeer{
        Name:      "bad-peer",
        RemoteIP:  "999.999.999.999",  // Invalid IP
        RemoteASN: 65001,
    }
    
    err := ctx.Client.CreatePeer(peer)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid IP address")
}
```

### Parameterized Tests

```go
func TestPeerValidation(t *testing.T) {
    testCases := []struct {
        name        string
        fixture     string
        expectError bool
        errorMsg    string
    }{
        {
            name:        "Valid_Basic_Peer",
            fixture:     "fixtures/peers/valid/basic_peer.yaml",
            expectError: false,
        },
        {
            name:        "Invalid_IP_Format",
            fixture:     "fixtures/peers/invalid/invalid_ip_format.yaml",
            expectError: true,
            errorMsg:    "invalid IP address",
        },
        {
            name:        "Missing_ASN",
            fixture:     "fixtures/peers/invalid/missing_asn.yaml",
            expectError: true,
            errorMsg:    "ASN is required",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctx := testutil.NewTestContext(t)
            defer ctx.Cleanup()
            
            peer := testutil.LoadPeerFixture(t, tc.fixture)
            err := ctx.Client.CreatePeer(peer)
            
            if tc.expectError {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tc.errorMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Async Operations

```go
func TestSessionEstablishment(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Create peer
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    err := ctx.Client.CreatePeer(peer)
    require.NoError(t, err)
    
    // Wait for session to establish (with timeout)
    timeout := time.After(10 * time.Second)
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            t.Fatal("Timeout waiting for session to establish")
        case <-ticker.C:
            session, err := ctx.Client.GetSession(peer.Name)
            require.NoError(t, err)
            
            if session.State == "Established" {
                assert.Greater(t, session.PrefixesReceived, 0)
                return
            }
        }
    }
}
```

---

## Fixtures

### What are Fixtures?

Fixtures are pre-defined test data files that provide consistent, reusable data for tests. They eliminate the need to hardcode test data in test files and make tests more maintainable.

### Fixture Types

#### 1. Peer Fixtures
**Location**: [`fixtures/peers/`](fixtures/peers/)

**Valid Peers** (`fixtures/peers/valid/`):
- [`basic_peer.yaml`](fixtures/peers/valid/basic_peer.yaml:1) - Simple IPv4 peer
- [`peer_with_password.yaml`](fixtures/peers/valid/peer_with_password.yaml:1) - Authenticated peer
- [`peer_with_multihop.yaml`](fixtures/peers/valid/peer_with_multihop.yaml:1) - Multi-hop BGP
- [`peer_full_config.yaml`](fixtures/peers/valid/peer_full_config.yaml:1) - All options configured

**Invalid Peers** (`fixtures/peers/invalid/`):
- [`invalid_ip_format.yaml`](fixtures/peers/invalid/invalid_ip_format.yaml:1) - Malformed IP
- [`missing_asn.yaml`](fixtures/peers/invalid/missing_asn.yaml:1) - Required field missing
- [`negative_asn.yaml`](fixtures/peers/invalid/negative_asn.yaml:1) - Invalid ASN value

#### 2. User Fixtures
**Location**: [`fixtures/users/`](fixtures/users/)

- [`admin_user.yaml`](fixtures/users/admin_user.yaml:1) - Administrator account
- [`regular_user.yaml`](fixtures/users/regular_user.yaml:1) - Standard user
- [`inactive_user.yaml`](fixtures/users/inactive_user.yaml:1) - Disabled account

#### 3. Session Fixtures
**Location**: [`fixtures/sessions/`](fixtures/sessions/)

- [`established_session.yaml`](fixtures/sessions/established_session.yaml:1) - Active BGP session
- [`idle_session.yaml`](fixtures/sessions/idle_session.yaml:1) - Idle state
- [`active_session.yaml`](fixtures/sessions/active_session.yaml:1) - Connecting state

#### 4. Alert Fixtures
**Location**: [`fixtures/alerts/`](fixtures/alerts/)

- [`peer_down_alert.yaml`](fixtures/alerts/peer_down_alert.yaml:1) - Peer disconnection
- [`peer_up_alert.yaml`](fixtures/alerts/peer_up_alert.yaml:1) - Peer connection
- [`max_prefixes_alert.yaml`](fixtures/alerts/max_prefixes_alert.yaml:1) - Prefix limit exceeded

### Loading Fixtures

```go
// Load peer fixture
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")

// Load user fixture
user := testutil.LoadUserFixture(t, "fixtures/users/admin_user.yaml")

// Load session fixture
session := testutil.LoadSessionFixture(t, "fixtures/sessions/established_session.yaml")

// Load alert fixture
alert := testutil.LoadAlertFixture(t, "fixtures/alerts/peer_down_alert.yaml")
```

### Creating Custom Fixtures

#### Peer Fixture Example

```yaml
# fixtures/peers/valid/my_custom_peer.yaml
name: custom-peer-1
description: My custom test peer
remote_ip: 192.168.100.1
remote_asn: 65100
local_ip: 192.168.100.2
local_asn: 65000
enabled: true
password: test-password
multihop: 2
update_source: eth0
route_map_in: IMPORT-POLICY
route_map_out: EXPORT-POLICY
prefix_list_in: PREFIX-IN
prefix_list_out: PREFIX-OUT
max_prefixes: 1000
local_preference: 100
```

#### User Fixture Example

```yaml
# fixtures/users/my_test_user.yaml
username: testuser
password: testpass123
email: testuser@example.com
role: operator
enabled: true
```

### Fixture Best Practices

1. **Use Descriptive Names**: `peer_with_authentication.yaml` not `peer2.yaml`
2. **Keep Data Realistic**: Use RFC 1918 IPs, private ASNs
3. **Document Purpose**: Add comments explaining the fixture's use case
4. **Maintain Consistency**: Follow existing fixture patterns
5. **Version Control**: Commit fixtures with tests
6. **Avoid Secrets**: Never use production credentials

---

## Debugging

### Debugging Failed Tests

#### 1. Review Test Output

```bash
# Run with verbose output
./run-tests.sh --verbose

# Run with debug logging
./run-tests.sh --log-level debug
```

#### 2. Preserve Artifacts

```bash
# Keep all artifacts on failure
./run-tests.sh --no-cleanup
```

This preserves:
- Test database (`tmp/test.db`)
- Log files (`logs/`)
- Test results (`results/`)

#### 3. Examine Logs

```bash
# View test execution log
cat logs/test-execution.log

# View mock FRR server log
cat logs/mock-frr-server.log

# View API server log
cat logs/api-server.log
```

#### 4. Inspect Database

```bash
# Open test database
sqlite3 tmp/test.db

# Query peers
SELECT * FROM bgp_peers;

# Query sessions
SELECT * FROM bgp_sessions;

# Query users
SELECT * FROM users;
```

#### 5. Test Mock FRR Server

```bash
# Check if server is running
curl http://localhost:51051/health

# View server statistics
curl http://localhost:51051/stats

# List all peers
curl http://localhost:51051/peers

# Get session state
curl "http://localhost:51051/sessions/state?ip=192.168.1.1"
```

### Common Issues and Solutions

#### Issue: Tests Hang

**Symptoms**: Tests don't complete, no output

**Causes**:
- Mock FRR server not running
- Port conflicts
- Network connectivity issues

**Solutions**:
```bash
# Check if mock server is running
ps aux | grep mock-frr-server

# Check port availability
lsof -i :50051
lsof -i :51051

# Restart mock server
cd cmd/mock-frr-server
./mock-frr-server &
```

#### Issue: Database Errors

**Symptoms**: `database locked`, `unable to open database`

**Causes**:
- Leftover database from previous run
- Insufficient permissions
- Concurrent access

**Solutions**:
```bash
# Clean database
./scripts/cleanup-db.sh

# Check permissions
ls -la tmp/

# Ensure single test instance
ps aux | grep "go test"
```

#### Issue: Authentication Failures

**Symptoms**: `401 Unauthorized`, `invalid token`

**Causes**:
- JWT secret mismatch
- Token expiry
- Incorrect credentials in fixtures

**Solutions**:
```bash
# Verify JWT secret in config
grep jwt_secret config/test-config.yaml

# Check user fixtures
cat fixtures/users/admin_user.yaml

# Test login manually
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

#### Issue: Port Already in Use

**Symptoms**: `bind: address already in use`

**Causes**:
- Previous test run didn't clean up
- Another service using the port

**Solutions**:
```bash
# Find process using port
lsof -i :8080
lsof -i :50051

# Kill process
kill -9 <PID>

# Or use cleanup script
./scripts/cleanup-all.sh
```

### Debugging Techniques

#### Enable Verbose Logging

```go
// In test file
func TestWithDebug(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    ctx.Logger.SetLevel(logrus.DebugLevel)
    defer ctx.Cleanup()
    
    // Your test code
}
```

#### Add Debug Prints

```go
func TestDebug(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    
    // Debug print
    t.Logf("Creating peer: %+v", peer)
    
    err := ctx.Client.CreatePeer(peer)
    require.NoError(t, err)
    
    // Debug print
    t.Logf("Peer created successfully")
}
```

#### Use Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug specific test
cd tests/01_authentication
dlv test -- -test.run TestLogin

# Set breakpoint and run
(dlv) break login_test.go:25
(dlv) continue
```

#### Capture HTTP Traffic

```bash
# Use tcpdump to capture traffic
sudo tcpdump -i lo -w test-traffic.pcap port 8080

# Analyze with wireshark
wireshark test-traffic.pcap
```

---

## CI/CD Integration

### GitHub Actions

Create `.github/workflows/functional-tests.yml`:

```yaml
name: Functional Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  functional-tests:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Check prerequisites
      run: |
        cd test/functional
        ./scripts/check-prerequisites.sh
    
    - name: Run functional tests
      run: |
        cd test/functional
        ./run-clean.sh
    
    - name: Upload test results
      if: always()
      uses: actions/upload-artifact@v3
      with:
        name: test-results
        path: test/functional/results/
    
    - name: Upload logs
      if: failure()
      uses: actions/upload-artifact@v3
      with:
        name: test-logs
        path: test/functional/logs/
```

### GitLab CI

Create `.gitlab-ci.yml`:

```yaml
functional-tests:
  stage: test
  image: golang:1.21
  
  before_script:
    - cd test/functional
    - ./scripts/check-prerequisites.sh
  
  script:
    - ./run-clean.sh
  
  artifacts:
    when: always
    paths:
      - test/functional/results/
    reports:
      junit: test/functional/results/test-results-*.xml
  
  only:
    - main
    - develop
    - merge_requests
```

### Jenkins

Create `Jenkinsfile`:

```groovy
pipeline {
    agent any
    
    stages {
        stage('Prerequisites') {
            steps {
                sh '''
                    cd test/functional
                    ./scripts/check-prerequisites.sh
                '''
            }
        }
        
        stage('Functional Tests') {
            steps {
                sh '''
                    cd test/functional
                    ./
run-clean.sh
                '''
            }
        }
        
        stage('Publish Results') {
            steps {
                junit 'test/functional/results/test-results-*.xml'
                archiveArtifacts artifacts: 'test/functional/results/**', allowEmptyArchive: true
            }
        }
    }
    
    post {
        failure {
            archiveArtifacts artifacts: 'test/functional/logs/**', allowEmptyArchive: true
        }
    }
}
```

### Docker-based CI

Create `Dockerfile.test`:

```dockerfile
FROM golang:1.21-alpine

# Install dependencies
RUN apk add --no-cache \
    bash \
    make \
    sqlite \
    git \
    curl

# Set working directory
WORKDIR /app

# Copy test files
COPY test/functional /app/test/functional

# Run tests
WORKDIR /app/test/functional
CMD ["./run-clean.sh"]
```

Run in CI:

```bash
docker build -f Dockerfile.test -t flintroute-tests .
docker run --rm flintroute-tests
```

### Test Reporting

#### JUnit XML Format

For CI integration, convert JSON results to JUnit XML:

```bash
# Install go-junit-report
go install github.com/jstemmer/go-junit-report@latest

# Convert results
go test -v ./... 2>&1 | go-junit-report > results/junit.xml
```

#### HTML Reports

Generate HTML reports for human review:

```bash
# Install go-test-report
go install github.com/vakenbolt/go-test-report@latest

# Generate HTML
cat results/test-results-*.json | go-test-report > results/report.html
```

---

## Best Practices

### Test Design Principles

#### 1. Independence
Each test should be completely independent:

```go
// ✅ Good - Independent test
func TestCreatePeer(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    err := ctx.Client.CreatePeer(peer)
    assert.NoError(t, err)
}

// ❌ Bad - Depends on previous test
func TestUpdatePeer(t *testing.T) {
    // Assumes peer from TestCreatePeer exists
    peer := &models.BGPPeer{Name: "test-peer"}
    err := ctx.Client.UpdatePeer(peer)
    assert.NoError(t, err)
}
```

#### 2. Clarity
Tests should be easy to understand:

```go
// ✅ Good - Clear intent
func TestLoginWithValidCredentials(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Arrange
    username := "admin"
    password := "admin123"
    
    // Act
    token, err := ctx.Client.Login(username, password)
    
    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, token)
}

// ❌ Bad - Unclear purpose
func TestLogin(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    x, y := ctx.Client.Login("admin", "admin123")
    assert.NoError(t, y)
}
```

#### 3. Completeness
Test both success and failure cases:

```go
func TestPeerCreation(t *testing.T) {
    // Test success case
    t.Run("Valid_Peer", func(t *testing.T) {
        ctx := testutil.NewTestContext(t)
        defer ctx.Cleanup()
        
        peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
        err := ctx.Client.CreatePeer(peer)
        assert.NoError(t, err)
    })
    
    // Test failure cases
    t.Run("Invalid_IP", func(t *testing.T) {
        ctx := testutil.NewTestContext(t)
        defer ctx.Cleanup()
        
        peer := testutil.LoadPeerFixture(t, "fixtures/peers/invalid/invalid_ip_format.yaml")
        err := ctx.Client.CreatePeer(peer)
        assert.Error(t, err)
    })
    
    t.Run("Duplicate_Name", func(t *testing.T) {
        ctx := testutil.NewTestContext(t)
        defer ctx.Cleanup()
        
        peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
        
        // Create once
        err := ctx.Client.CreatePeer(peer)
        require.NoError(t, err)
        
        // Try to create again
        err = ctx.Client.CreatePeer(peer)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "already exists")
    })
}
```

### Naming Conventions

#### Test Files
- Use `_test.go` suffix
- Match feature name: `login_test.go`, `peer_management_test.go`

#### Test Functions
- Start with `Test`
- Use descriptive names: `TestLoginWithValidCredentials`
- Use underscores in subtests: `Create_Valid_Peer`

#### Variables
- Use descriptive names: `expectedToken`, not `t`
- Use `ctx` for test context
- Use `err` for errors

### Error Handling

```go
// ✅ Good - Proper error handling
func TestFeature(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    result, err := ctx.Client.DoSomething()
    require.NoError(t, err, "DoSomething should not return error")
    assert.NotNil(t, result, "Result should not be nil")
}

// ❌ Bad - Ignoring errors
func TestFeature(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    result, _ := ctx.Client.DoSomething()
    assert.NotNil(t, result)
}
```

### Resource Cleanup

Always clean up resources:

```go
func TestWithCleanup(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()  // Always defer cleanup
    
    // Create resources
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    err := ctx.Client.CreatePeer(peer)
    require.NoError(t, err)
    
    // Test operations
    // ...
    
    // Cleanup happens automatically via defer
}
```

### Test Data Management

#### Use Fixtures
```go
// ✅ Good - Use fixtures
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")

// ❌ Bad - Hardcode data
peer := &models.BGPPeer{
    Name:      "test-peer",
    RemoteIP:  "192.168.1.1",
    RemoteASN: 65001,
}
```

#### Customize When Needed
```go
// Load fixture as template
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")

// Customize for this test
peer.Name = "unique-peer-" + uuid.New().String()
peer.RemoteIP = "192.168.100.1"
```

### Assertion Guidelines

#### Use Appropriate Assertions

```go
// For critical checks that should stop test
require.NoError(t, err)
require.NotNil(t, result)

// For non-critical checks
assert.Equal(t, expected, actual)
assert.Contains(t, str, substring)
assert.Greater(t, value, threshold)
```

#### Provide Context

```go
// ✅ Good - Descriptive message
assert.Equal(t, "Established", session.State, 
    "Session should be in Established state after 5 seconds")

// ❌ Bad - No context
assert.Equal(t, "Established", session.State)
```

### Performance Considerations

#### Use Timeouts

```go
func TestWithTimeout(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Set test timeout
    ctx.SetTimeout(30 * time.Second)
    
    // Your test code
}
```

#### Avoid Sleep

```go
// ❌ Bad - Fixed sleep
time.Sleep(5 * time.Second)

// ✅ Good - Poll with timeout
err := testutil.WaitForCondition(t, 10*time.Second, func() bool {
    session, _ := ctx.Client.GetSession("peer1")
    return session.State == "Established"
})
require.NoError(t, err)
```

### Documentation

#### Comment Complex Logic

```go
func TestComplexScenario(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Create peer and wait for initial connection
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    err := ctx.Client.CreatePeer(peer)
    require.NoError(t, err)
    
    // Wait for session to establish (BGP takes time)
    // Expected states: Idle -> Connect -> Active -> OpenSent -> OpenConfirm -> Established
    err = testutil.WaitForSessionState(t, ctx, peer.Name, "Established", 10*time.Second)
    require.NoError(t, err, "Session should establish within 10 seconds")
    
    // Verify metrics are populated
    session, err := ctx.Client.GetSession(peer.Name)
    require.NoError(t, err)
    assert.Greater(t, session.PrefixesReceived, 0, "Should have received prefixes")
}
```

#### Document Test Purpose

```go
// TestPeerLifecycle validates the complete lifecycle of a BGP peer:
// 1. Creation with valid configuration
// 2. Session establishment
// 3. Configuration updates
// 4. Graceful shutdown
// 5. Cleanup
func TestPeerLifecycle(t *testing.T) {
    // Test implementation
}
```

### Anti-Patterns to Avoid

#### ❌ Don't Share State Between Tests

```go
// Bad - Global state
var globalPeer *models.BGPPeer

func TestCreate(t *testing.T) {
    globalPeer = &models.BGPPeer{...}
}

func TestUpdate(t *testing.T) {
    // Uses globalPeer from TestCreate
}
```

#### ❌ Don't Test Implementation Details

```go
// Bad - Testing internal implementation
func TestInternalCache(t *testing.T) {
    cache := ctx.Client.GetInternalCache()
    assert.NotNil(t, cache)
}

// Good - Test behavior
func TestPeerRetrieval(t *testing.T) {
    peer, err := ctx.Client.GetPeer("peer1")
    assert.NoError(t, err)
    assert.NotNil(t, peer)
}
```

#### ❌ Don't Make Tests Too Complex

```go
// Bad - Too complex
func TestEverything(t *testing.T) {
    // 500 lines of test code testing multiple features
}

// Good - Focused tests
func TestPeerCreation(t *testing.T) {
    // 20 lines testing one thing
}

func TestPeerUpdate(t *testing.T) {
    // 20 lines testing one thing
}
```

---

## Appendix

### A. Directory Reference

```
test/functional/
├── cmd/                    # Executable commands
│   └── mock-frr-server/   # Mock FRR gRPC server
├── config/                 # Configuration files
│   ├── test-config.yaml
│   ├── mock-frr-config.yaml
│   └── logging-config.yaml
├── fixtures/               # Test data
│   ├── peers/
│   ├── users/
│   ├── sessions/
│   ├── alerts/
│   └── config/
├── logs/                   # Test execution logs
├── results/                # Test results and reports
├── scripts/                # Helper scripts
│   ├── check-prerequisites.sh
│   ├── setup-env.sh
│   ├── teardown-env.sh
│   ├── cleanup-*.sh
│   └── wait-for-service.sh
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
├── run-clean.sh           # Clean run script
└── run-retest.sh          # Rerun failed tests
```

### B. Configuration Reference

#### test-config.yaml

```yaml
server:
  host: 127.0.0.1
  port: 0  # Random available port
  
database:
  path: ./tmp/test.db
  
frr:
  grpc_host: localhost
  grpc_port: 50051
  
auth:
  jwt_secret: test-secret-key-for-functional-testing
  token_expiry: 5m
  refresh_expiry: 1h

testing:
  timeout: 30s
  cleanup_on_success: false
  log_level: debug
  parallel: false
```

#### mock-frr-config.yaml

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

### C. Command Reference

#### Test Execution

```bash
# Run all tests
./run-tests.sh

# Run specific suite
./run-tests.sh --pattern ./tests/01_authentication/...

# Run with options
./run-tests.sh --verbose --log-level debug --no-cleanup

# Clean run
./run-clean.sh

# Rerun failed tests
./run-retest.sh
```

#### Environment Management

```bash
# Check prerequisites
./scripts/check-prerequisites.sh

# Setup environment
./scripts/setup-env.sh

# Teardown environment
./scripts/teardown-env.sh

# Full cleanup
./scripts/cleanup-all.sh
```

#### Mock Server

```bash
# Start mock server
cd cmd/mock-frr-server
./mock-frr-server

# Test mock server
./test-server.sh

# Build mock server
make build

# Run tests
make test
```

### D. Troubleshooting Checklist

- [ ] Prerequisites met (`./scripts/check-prerequisites.sh`)
- [ ] Ports available (8080, 50051, 51051)
- [ ] Configuration files present
- [ ] Write permissions on logs/, results/, tmp/
- [ ] No leftover processes from previous runs
- [ ] Mock FRR server running
- [ ] Database not locked
- [ ] Sufficient disk space

### E. Additional Resources

- **Main Documentation**: [`test/functional/README.md`](README.md:1)
- **Mock FRR Server**: [`cmd/mock-frr-server/README.md`](cmd/mock-frr-server/README.md:1)
- **Fixtures Guide**: [`fixtures/README.md`](fixtures/README.md:1)
- **API Reference**: [`API_REFERENCE.md`](API_REFERENCE.md:1)
- **Quick Reference**: [`QUICK_REFERENCE.md`](QUICK_REFERENCE.md:1)
- **Contributing**: [`CONTRIBUTING.md`](CONTRIBUTING.md:1)
- **FAQ**: [`FAQ.md`](FAQ.md:1)

### F. Glossary

- **BGP**: Border Gateway Protocol - Routing protocol for the internet
- **FRR**: Free Range Routing - Open source routing software
- **gRPC**: Google Remote Procedure Call - High-performance RPC framework
- **JWT**: JSON Web Token - Authentication token format
- **Mock**: Simulated component for testing
- **Fixture**: Pre-defined test data
- **Assertion**: Test verification statement
- **Context**: Test execution environment

---

## Support

### Getting Help

- **Documentation**: Read this guide and related docs
- **Issues**: Check existing GitHub issues
- **Logs**: Review test execution logs
- **Community**: Ask in project discussions

### Reporting Issues

When reporting test failures, include:

1. Test command used
2. Full error output
3. Log files (if available)
4. System information
5. Steps to reproduce

### Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md:1) for guidelines on:
- Adding new tests
- Improving documentation
- Reporting bugs
- Submitting pull requests

---

**Document Version**: 1.0.0  
**Last Updated**: November 10, 2025  
**Maintained By**: FlintRoute Development Team