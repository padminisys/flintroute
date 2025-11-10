# Mock FRR gRPC Server

A mock implementation of the FRR (Free Range Routing) gRPC server for functional testing of FlintRoute.

## Overview

This mock server simulates FRR's gRPC interface, allowing realistic testing of FlintRoute without requiring an actual FRR instance. It maintains in-memory state of BGP peers and sessions, simulates session state transitions, and provides both gRPC and HTTP interfaces for testing.

## Features

- **BGP Peer Management**: Add, remove, and update BGP peer configurations
- **Session State Simulation**: Realistic BGP session state transitions (Idle → Connect → Active → OpenSent → OpenConfirm → Established)
- **Thread-Safe State Management**: Concurrent request handling with mutex protection
- **Configurable Behavior**: Control simulation delays and error injection via YAML configuration
- **Dual Interface**: gRPC for production-like testing, HTTP for debugging and manual testing
- **Comprehensive Logging**: Structured logging with configurable levels and file output

## Building

```bash
cd test/functional/cmd/mock-frr-server
go build -o mock-frr-server
```

## Running

### Using Default Configuration

```bash
./mock-frr-server
```

This uses the default configuration at `test/functional/config/mock-frr-config.yaml`.

### Using Custom Configuration

```bash
./mock-frr-server -config /path/to/config.yaml
```

## Configuration

The server is configured via YAML. See [`test/functional/config/mock-frr-config.yaml`](../../config/mock-frr-config.yaml) for the default configuration:

```yaml
server:
  host: localhost
  port: 50051              # gRPC port
  
simulation:
  session_state_delay: 100ms  # Delay between BGP state transitions
  error_injection: false      # Enable to simulate errors
  
logging:
  level: info                 # debug, info, warn, error
  file: ./logs/mock-frr-server.log
```

## Interfaces

### gRPC Interface (Port 50051)

The gRPC interface implements the FRR service with the following methods:

- `AddBGPPeer` - Add a BGP peer
- `RemoveBGPPeer` - Remove a BGP peer
- `UpdateBGPPeer` - Update peer configuration
- `GetBGPSessionState` - Get session state for a specific peer
- `GetAllBGPSessions` - Get all session states
- `GetRunningConfig` - Get mock FRR configuration

### HTTP Debug Interface (Port 51051)

The HTTP interface provides easy testing and debugging:

#### Health Check
```bash
curl http://localhost:51051/health
```

#### Get Statistics
```bash
curl http://localhost:51051/stats
```

#### List All Peers
```bash
curl http://localhost:51051/peers
```

#### Add a Peer
```bash
curl -X POST http://localhost:51051/peers/add \
  -H "Content-Type: application/json" \
  -d '{
    "IPAddress": "192.168.1.1",
    "ASN": 65000,
    "RemoteASN": 65001
  }'
```

#### Remove a Peer
```bash
curl -X POST http://localhost:51051/peers/remove \
  -H "Content-Type: application/json" \
  -d '{"ip_address": "192.168.1.1"}'
```

#### Get All Sessions
```bash
curl http://localhost:51051/sessions
```

#### Get Session State
```bash
curl "http://localhost:51051/sessions/state?ip=192.168.1.1"
```

#### Get Running Config
```bash
curl http://localhost:51051/config
```

## BGP Session State Simulation

When a peer is added, the server automatically simulates the BGP session establishment process:

1. **Idle** (initial state)
2. **Connect** (after `session_state_delay`)
3. **Active** (after `session_state_delay`)
4. **OpenSent** (after `session_state_delay`)
5. **OpenConfirm** (after `session_state_delay`)
6. **Established** (after `session_state_delay`)

Once established, the session includes simulated metrics:
- Prefixes Received: 100
- Prefixes Sent: 50
- Messages Received: 1000
- Messages Sent: 900
- Uptime: Calculated from state change time

## Error Injection

Enable error injection in the configuration to test error handling:

```yaml
simulation:
  error_injection: true
```

When enabled, all peer operations will return errors, allowing you to test error scenarios.

## State Management

The server maintains two types of state:

### Peer State
- IP Address
- ASN (Autonomous System Number)
- Remote ASN
- Password
- Multihop settings
- Update source
- Route maps (in/out)
- Prefix lists (in/out)
- Maximum prefixes
- Local preference

### Session State
- IP Address
- BGP State (Idle, Connect, Active, OpenSent, OpenConfirm, Established)
- Uptime (seconds)
- Prefixes Received/Sent
- Messages Received/Sent
- Last Error

## Thread Safety

All state operations are protected by `sync.RWMutex`, ensuring safe concurrent access from multiple gRPC/HTTP requests.

## Logging

Logs are written to both:
- **Console**: Human-readable format
- **File**: JSON format for parsing

Log levels: `debug`, `info`, `warn`, `error`

## Integration with FlintRoute Tests

The mock server is designed to be used in FlintRoute's functional tests:

1. Start the mock server before running tests
2. Configure FlintRoute to connect to `localhost:50051`
3. Run functional tests
4. Stop the mock server after tests complete

Example test workflow:
```bash
# Start mock server
./mock-frr-server &
MOCK_PID=$!

# Run tests
go test ./test/functional/tests/...

# Stop mock server
kill $MOCK_PID
```

## Architecture

```
main.go          - Entry point, signal handling, logger initialization
config.go        - Configuration loading and validation
state.go         - BGP state management (peers and sessions)
server.go        - gRPC and HTTP server implementation
proto/frr.proto  - Protocol buffer definitions (for reference)
```

## Development

### Adding New Features

1. Update `state.go` if new state is needed
2. Add methods to `server.go` for new operations
3. Update proto definitions if changing the interface
4. Add corresponding HTTP endpoints for debugging

### Testing

The mock server can be tested manually using the HTTP interface or by writing integration tests that use the gRPC interface.

## Troubleshooting

### Server Won't Start

- Check if port 50051 (gRPC) or 51051 (HTTP) is already in use
- Verify the configuration file path is correct
- Check log file permissions

### Sessions Not Establishing

- Verify `session_state_delay` is set appropriately
- Check if error injection is enabled
- Review server logs for errors

### Connection Refused

- Ensure the server is running
- Verify the correct host and port
- Check firewall settings

## License

This is part of the FlintRoute project. See the main project LICENSE file.