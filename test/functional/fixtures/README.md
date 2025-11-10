# Test Fixtures

This directory contains test data fixtures used across the functional test suite.

## Overview

Fixtures provide consistent, reusable test data for:
- BGP peer configurations
- User accounts and credentials
- Session states
- Configuration templates
- Error scenarios

## Directory Structure

```
fixtures/
├── peers/
│   ├── valid/          # Valid peer configurations
│   └── invalid/        # Invalid configurations for error testing
├── users/              # User account fixtures
└── sessions/           # BGP session state fixtures
```

## Fixture Formats

### JSON Format
Primary format for structured data:

```json
{
  "name": "peer1",
  "description": "Test BGP peer",
  "remote_asn": 65001,
  "remote_ip": "192.168.1.1",
  "local_ip": "192.168.1.2"
}
```

### YAML Format
Alternative format for configuration-style data:

```yaml
name: peer1
description: Test BGP peer
remote_asn: 65001
remote_ip: 192.168.1.1
local_ip: 192.168.1.2
```

## Peer Fixtures

### Valid Peers (`peers/valid/`)

**peer1.json** - Basic IPv4 peer
```json
{
  "name": "peer1",
  "description": "Basic IPv4 BGP peer",
  "remote_asn": 65001,
  "remote_ip": "192.168.1.1",
  "local_ip": "192.168.1.2",
  "peer_group": "external",
  "enabled": true
}
```

**peer-ipv6.json** - IPv6 peer
```json
{
  "name": "peer-ipv6",
  "description": "IPv6 BGP peer",
  "remote_asn": 65002,
  "remote_ip": "2001:db8::1",
  "local_ip": "2001:db8::2",
  "peer_group": "external",
  "enabled": true
}
```

**peer-with-password.json** - Authenticated peer
```json
{
  "name": "peer-auth",
  "description": "Peer with MD5 authentication",
  "remote_asn": 65003,
  "remote_ip": "192.168.2.1",
  "local_ip": "192.168.2.2",
  "password": "test-password-123",
  "enabled": true
}
```

### Invalid Peers (`peers/invalid/`)

**peer-invalid-ip.json** - Invalid IP address
```json
{
  "name": "peer-bad-ip",
  "remote_asn": 65001,
  "remote_ip": "999.999.999.999",
  "local_ip": "192.168.1.2"
}
```

**peer-invalid-asn.json** - Invalid ASN
```json
{
  "name": "peer-bad-asn",
  "remote_asn": 4294967296,
  "remote_ip": "192.168.1.1",
  "local_ip": "192.168.1.2"
}
```

**peer-missing-fields.json** - Missing required fields
```json
{
  "name": "peer-incomplete",
  "remote_asn": 65001
}
```

## User Fixtures

### users/admin.json
```json
{
  "username": "admin",
  "password": "admin123",
  "email": "admin@example.com",
  "role": "admin",
  "enabled": true
}
```

### users/operator.json
```json
{
  "username": "operator",
  "password": "operator123",
  "email": "operator@example.com",
  "role": "operator",
  "enabled": true
}
```

### users/viewer.json
```json
{
  "username": "viewer",
  "password": "viewer123",
  "email": "viewer@example.com",
  "role": "viewer",
  "enabled": true
}
```

## Session Fixtures

### sessions/established.json
```json
{
  "peer_name": "peer1",
  "state": "Established",
  "uptime": 3600,
  "prefixes_received": 100,
  "prefixes_sent": 50
}
```

### sessions/idle.json
```json
{
  "peer_name": "peer2",
  "state": "Idle",
  "uptime": 0,
  "prefixes_received": 0,
  "prefixes_sent": 0
}
```

### sessions/active.json
```json
{
  "peer_name": "peer3",
  "state": "Active",
  "uptime": 30,
  "prefixes_received": 0,
  "prefixes_sent": 0
}
```

## Using Fixtures in Tests

### Loading Fixtures

```go
import "github.com/yourusername/flintroute/test/functional/pkg/testutil"

// Load peer fixture
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/peer1.json")

// Load user fixture
user := testutil.LoadUserFixture(t, "fixtures/users/admin.json")

// Load session fixture
session := testutil.LoadSessionFixture(t, "fixtures/sessions/established.json")
```

### Creating Test Data

```go
// Use fixture as template
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/peer1.json")
peer.Name = "test-peer-" + uuid.New().String()
peer.RemoteIP = "192.168.100.1"

// Create in test database
err := ctx.Client.CreatePeer(peer)
require.NoError(t, err)
```

### Parameterized Tests

```go
func TestPeerValidation(t *testing.T) {
    fixtures := []string{
        "fixtures/peers/invalid/peer-invalid-ip.json",
        "fixtures/peers/invalid/peer-invalid-asn.json",
        "fixtures/peers/invalid/peer-missing-fields.json",
    }
    
    for _, fixture := range fixtures {
        t.Run(fixture, func(t *testing.T) {
            peer := testutil.LoadPeerFixture(t, fixture)
            err := ctx.Client.CreatePeer(peer)
            assert.Error(t, err)
        })
    }
}
```

## Fixture Guidelines

### Naming Conventions
- Use descriptive names: `peer-ipv6.json`, not `p2.json`
- Include scenario in name: `peer-invalid-ip.json`
- Use lowercase with hyphens: `peer-with-password.json`

### Data Consistency
- Use realistic but non-production values
- Keep ASNs in private range (64512-65534)
- Use RFC 1918 private IPs (192.168.x.x, 10.x.x.x)
- Use documentation IPv6 prefix (2001:db8::/32)

### Maintenance
- Update fixtures when models change
- Add new fixtures for new test scenarios
- Remove obsolete fixtures
- Document fixture purpose in comments

### Security
- Never use production credentials
- Use obvious test passwords
- Don't commit sensitive data
- Rotate test credentials periodically

## Creating New Fixtures

### Template
```json
{
  "_comment": "Description of this fixture and its purpose",
  "name": "fixture-name",
  "field1": "value1",
  "field2": "value2"
}
```

### Validation
```bash
# Validate JSON syntax
jq . fixtures/peers/valid/new-peer.json

# Validate against schema (if available)
jsonschema -i fixtures/peers/valid/new-peer.json schema/peer.json
```

### Testing
```go
func TestNewFixture(t *testing.T) {
    // Ensure fixture loads correctly
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/new-peer.json")
    assert.NotNil(t, peer)
    
    // Validate fixture data
    assert.NotEmpty(t, peer.Name)
    assert.NotEmpty(t, peer.RemoteIP)
}
```

## Fixture Organization

### By Feature
Group fixtures by the feature they test:
- `peers/valid/` - Valid peer configurations
- `peers/invalid/` - Invalid configurations
- `peers/edge-cases/` - Edge case scenarios

### By Scenario
Group fixtures by test scenario:
- `workflows/peer-lifecycle/` - Complete peer lifecycle
- `workflows/session-recovery/` - Session recovery scenarios
- `workflows/configuration-update/` - Configuration updates

### By Complexity
Group fixtures by complexity:
- `simple/` - Basic, minimal fixtures
- `complex/` - Full-featured fixtures
- `realistic/` - Production-like fixtures

## Best Practices

1. **Reusability**: Create generic fixtures that can be customized
2. **Clarity**: Use descriptive names and include comments
3. **Completeness**: Include all required fields
4. **Validity**: Ensure valid fixtures are actually valid
5. **Coverage**: Cover common cases and edge cases
6. **Maintenance**: Keep fixtures up-to-date with code changes
7. **Documentation**: Document fixture purpose and usage

## Troubleshooting

### Fixture Not Found
```
fixture not found: fixtures/peers/valid/peer1.json
```
- Verify file path is correct
- Check file exists: `ls -la fixtures/peers/valid/`
- Ensure working directory is correct

### Invalid JSON
```
invalid character '}' looking for beginning of object key string
```
- Validate JSON: `jq . fixture.json`
- Check for trailing commas
- Verify proper escaping

### Schema Mismatch
```
fixture does not match expected schema
```
- Compare with current model definition
- Update fixture to match schema
- Check for deprecated fields