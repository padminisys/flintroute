# FlintRoute - BGP Management System

A modern web-based BGP management system built with Go and React, designed for managing FRRouting (FRR) BGP configurations.

## Features

### Phase 1 (MVP)
- ✅ **Authentication**: JWT-based authentication with refresh tokens
- ✅ **BGP Peer Management**: Create, read, update, and delete BGP peers
- ✅ **Session Monitoring**: Real-time BGP session state monitoring
- ✅ **Configuration Management**: Backup and restore FRR configurations
- ✅ **Alerting System**: Peer state change alerts with acknowledgment
- ✅ **Real-time Updates**: WebSocket-based live updates
- ✅ **RESTful API**: Complete API for all operations

## Architecture

### Backend (Go)
- **Framework**: Gin (HTTP), Gorilla WebSocket
- **Database**: SQLite with GORM
- **Authentication**: JWT tokens
- **FRR Integration**: gRPC client (stub implementation)
- **Logging**: Structured logging with zap

### Frontend (React + TypeScript)
- **Framework**: React 18 with Vite
- **State Management**: Redux Toolkit
- **UI Library**: Material-UI (MUI)
- **API Client**: Axios with interceptors
- **Real-time**: WebSocket client

## Project Structure

```
flintroute/
├── cmd/
│   └── flintroute/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/                        # HTTP/WebSocket handlers
│   ├── auth/                       # JWT authentication
│   ├── bgp/                        # BGP management service
│   ├── config/                     # Configuration management
│   ├── database/                   # Database layer
│   ├── frr/                        # FRR gRPC client
│   ├── models/                     # Data models
│   └── websocket/                  # WebSocket server
├── frontend/                       # React application
│   ├── src/
│   │   ├── components/            # React components
│   │   ├── pages/                 # Page components
│   │   ├── services/              # API services
│   │   ├── store/                 # Redux store
│   │   └── App.tsx
│   ├── package.json
│   └── vite.config.ts
├── configs/
│   ├── config.example.yaml        # Example configuration
│   └── frr/                       # FRR test configs
├── docs/                          # Documentation
├── docker-compose.yml             # FRR test environment
├── Makefile                       # Build automation
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Docker and Docker Compose (for FRR testing)
- Git

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/padminisys/flintroute.git
cd flintroute
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Configure the Application

```bash
cp configs/config.example.yaml configs/config.yaml
# Edit configs/config.yaml with your settings
```

### 4. Build the Application

```bash
make build
```

### 5. Run the Application

**Option A: Run backend and frontend separately**

Terminal 1 (Backend):
```bash
make dev-backend
```

Terminal 2 (Frontend):
```bash
make dev-frontend
```

**Option B: Run both together**
```bash
make dev
```

### 6. Access the Application

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Default credentials: `admin` / `admin`

## Development

### Backend Development

```bash
# Run with hot reload (requires air)
make dev-hot

# Run tests
make test

# Build binary
make build
```

### Frontend Development

```bash
cd frontend

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Testing with FRR

Start the FRR test environment:

```bash
make docker-up
```

Stop the FRR test environment:

```bash
make docker-down
```

## API Documentation

### Authentication

```bash
# Login
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin"
}

# Refresh token
POST /api/v1/auth/refresh
{
  "refresh_token": "your-refresh-token"
}

# Logout
POST /api/v1/auth/logout
```

### BGP Peers

```bash
# List all peers
GET /api/v1/bgp/peers

# Get specific peer
GET /api/v1/bgp/peers/:id

# Create peer
POST /api/v1/bgp/peers
{
  "name": "Peer1",
  "ip_address": "192.168.1.1",
  "asn": 65001,
  "remote_asn": 65002,
  "enabled": true
}

# Update peer
PUT /api/v1/bgp/peers/:id

# Delete peer
DELETE /api/v1/bgp/peers/:id
```

### BGP Sessions

```bash
# List all sessions
GET /api/v1/bgp/sessions

# Get specific session
GET /api/v1/bgp/sessions/:id
```

### Configuration

```bash
# List configuration versions
GET /api/v1/config/versions

# Backup current configuration
POST /api/v1/config/backup
{
  "description": "Before maintenance"
}

# Restore configuration
POST /api/v1/config/restore/:id
```

### Alerts

```bash
# List alerts
GET /api/v1/alerts?acknowledged=false&severity=warning

# Acknowledge alert
POST /api/v1/alerts/:id/acknowledge
```

### WebSocket

```bash
# Connect to WebSocket
WS /api/v1/ws

# Message types:
# - session_update: BGP session state changes
# - peer_update: BGP peer configuration changes
# - alert: New alerts
```

## Configuration

### Backend Configuration (configs/config.yaml)

```yaml
server:
  host: 0.0.0.0
  port: 8080

database:
  path: ./data/flintroute.db

frr:
  grpc_host: localhost
  grpc_port: 50051

auth:
  jwt_secret: your-secret-key-here
  token_expiry: 15m
  refresh_expiry: 168h  # 7 days
```

### Frontend Configuration (frontend/.env)

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080
```

## Security

- Change the default admin password immediately after first login
- Use a strong JWT secret in production
- Enable HTTPS in production
- Restrict CORS origins in production
- Keep dependencies up to date

## Troubleshooting

### Backend won't start
- Check if port 8080 is available
- Verify database path is writable
- Check logs for detailed error messages

### Frontend won't connect to backend
- Verify backend is running
- Check CORS settings
- Verify API URL in frontend/.env

### WebSocket connection fails
- Ensure you're authenticated
- Check WebSocket URL configuration
- Verify firewall settings

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

See [docs/phase1/roadmap.md](docs/phase1/roadmap.md) for the complete roadmap.

### Upcoming Features (Phase 2+)
- Route policy management
- Prefix list management
- Community management
- BGP route visualization
- Multi-router support
- RBAC (Role-Based Access Control)
- Audit logging
- Metrics and monitoring
- Backup scheduling

## Support

For issues and questions:
- GitHub Issues: https://github.com/padminisys/flintroute/issues
- Documentation: [docs/](docs/)

## Acknowledgments

- FRRouting project for the excellent routing software
- The Go and React communities for amazing tools and libraries