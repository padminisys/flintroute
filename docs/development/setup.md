# Development Setup Guide

## Table of Contents
- [Prerequisites](#prerequisites)
- [System Requirements](#system-requirements)
- [Go Installation](#go-installation)
- [Node.js Installation](#nodejs-installation)
- [Docker Setup](#docker-setup)
- [IDE Setup](#ide-setup)
- [Project Structure](#project-structure)
- [Environment Configuration](#environment-configuration)
- [Running Development Servers](#running-development-servers)
- [Hot Reload Setup](#hot-reload-setup)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before starting, ensure you have:
- Debian 12 (Bookworm) or compatible Linux distribution
- Root or sudo access
- Stable internet connection
- At least 8GB RAM and 20GB free disk space

---

## System Requirements

### Minimum Requirements
- **OS**: Debian 12 (Bookworm), Ubuntu 22.04+, or compatible
- **CPU**: 2 cores
- **RAM**: 4GB
- **Disk**: 10GB free space

### Recommended Requirements
- **OS**: Debian 12 (Bookworm)
- **CPU**: 4 cores
- **RAM**: 8GB
- **Disk**: 20GB free space

---

## Go Installation

### Install Go 1.21+

```bash
# Remove any existing Go installation
sudo rm -rf /usr/local/go

# Download Go 1.21.6 (or latest 1.21.x)
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# Extract to /usr/local
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# Clean up
rm go1.21.6.linux-amd64.tar.gz
```

### Configure Go Environment

Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Go environment
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Apply changes:

```bash
source ~/.bashrc
```

### Verify Installation

```bash
go version
# Expected output: go version go1.21.6 linux/amd64

go env GOPATH
# Expected output: /home/yourusername/go
```

### Install Go Tools

```bash
# Install common Go development tools
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

---

## Node.js Installation

### Install Node.js 18+ via NodeSource

```bash
# Install prerequisites
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg

# Add NodeSource repository
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# Install Node.js
sudo apt-get install -y nodejs

# Verify installation
node --version
# Expected output: v18.x.x

npm --version
# Expected output: 9.x.x or higher
```

### Alternative: Install via nvm (Recommended for Multiple Versions)

```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Load nvm
source ~/.bashrc

# Install Node.js 18
nvm install 18
nvm use 18
nvm alias default 18

# Verify
node --version
npm --version
```

### Configure npm

```bash
# Set npm prefix to avoid permission issues
mkdir -p ~/.npm-global
npm config set prefix '~/.npm-global'

# Add to PATH in ~/.bashrc
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

### Install Global npm Packages

```bash
# Install essential development tools
npm install -g typescript
npm install -g vite
npm install -g eslint
npm install -g prettier
```

---

## Docker Setup

### Install Docker

```bash
# Update package index
sudo apt-get update

# Install prerequisites
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Add Docker's official GPG key
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | \
    sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Set up repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/debian \
  $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Verify installation
sudo docker --version
sudo docker compose version
```

### Configure Docker for Non-Root User

```bash
# Create docker group (if not exists)
sudo groupadd docker

# Add your user to docker group
sudo usermod -aG docker $USER

# Apply group changes (logout/login or use newgrp)
newgrp docker

# Verify non-root access
docker run hello-world
```

### Install Docker Compose (Standalone)

```bash
# Download Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-$(uname -s)-$(uname -m)" \
    -o /usr/local/bin/docker-compose

# Make executable
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker-compose --version
```

---

## IDE Setup

### VSCode (Recommended)

#### Install VSCode

```bash
# Download and install VSCode
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] \
    https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
rm -f packages.microsoft.gpg

sudo apt-get update
sudo apt-get install -y code
```

#### Install VSCode Extensions

```bash
# Go extensions
code --install-extension golang.go
code --install-extension ms-vscode.go

# TypeScript/React extensions
code --install-extension dbaeumer.vscode-eslint
code --install-extension esbenp.prettier-vscode
code --install-extension dsznajder.es7-react-js-snippets
code --install-extension bradlc.vscode-tailwindcss

# Docker extensions
code --install-extension ms-azuretools.vscode-docker

# Git extensions
code --install-extension eamodio.gitlens

# General productivity
code --install-extension editorconfig.editorconfig
code --install-extension gruntfuggly.todo-tree
```

#### VSCode Settings

Create `.vscode/settings.json` in project root:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[go]": {
    "editor.defaultFormatter": "golang.go"
  }
}
```

### Alternative: GoLand / WebStorm

If using JetBrains IDEs:

```bash
# Download from https://www.jetbrains.com/go/
# Or use snap
sudo snap install goland --classic
sudo snap install webstorm --classic
```

---

## Project Structure

### Clone Repository

```bash
# Clone the repository
git clone https://github.com/padminisys/flintroute.git
cd flintroute

# Verify structure
tree -L 2 -d
```

### Expected Directory Structure

```
flintroute/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   │   └── server/        # Main server application
│   ├── internal/           # Private application code
│   │   ├── api/           # HTTP/WebSocket handlers
│   │   ├── auth/          # Authentication & authorization
│   │   ├── config/        # Configuration management
│   │   ├── frr/           # FRR gRPC client
│   │   ├── models/        # Data models
│   │   ├── repository/    # Data access layer
│   │   └── service/       # Business logic
│   ├── pkg/               # Public libraries
│   ├── go.mod             # Go module definition
│   └── go.sum             # Go dependencies
├── frontend/               # React frontend application
│   ├── src/               # Source code
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── store/         # Redux store
│   │   ├── services/      # API services
│   │   ├── hooks/         # Custom React hooks
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   ├── package.json       # npm dependencies
│   └── vite.config.ts     # Vite configuration
├── docs/                   # Documentation
├── scripts/                # Build and deployment scripts
├── docker/                 # Docker configurations
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── docker-compose.yml
├── .github/                # GitHub Actions workflows
├── Makefile               # Build automation
└── README.md              # Project overview
```

### Initialize Project Structure

If starting from scratch:

```bash
# Create backend structure
mkdir -p backend/{cmd/server,internal/{api,auth,config,frr,models,repository,service},pkg}

# Create frontend structure
mkdir -p frontend/{src/{components,pages,store,services,hooks,utils},public}

# Create other directories
mkdir -p {docs,scripts,docker,.github/workflows}
```

---

## Environment Configuration

### Backend Configuration

Create `backend/.env.development`:

```bash
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=development

# Database Configuration
DB_TYPE=sqlite
DB_PATH=./data/flintroute.db

# FRR Configuration
FRR_GRPC_HOST=localhost
FRR_GRPC_PORT=50051
FRR_GRPC_TLS=false

# JWT Configuration
JWT_SECRET=your-development-secret-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# WebSocket
WS_PING_INTERVAL=30s
WS_PONG_TIMEOUT=60s
```

### Frontend Configuration

Create `frontend/.env.development`:

```bash
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws

# Application Configuration
VITE_APP_NAME=FlintRoute
VITE_APP_VERSION=0.1.0-dev

# Feature Flags
VITE_ENABLE_DEBUG=true
VITE_ENABLE_MOCK_DATA=false
```

### Create Configuration Files

```bash
# Backend
cd backend
cat > .env.development << 'EOF'
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=development
DB_TYPE=sqlite
DB_PATH=./data/flintroute.db
FRR_GRPC_HOST=localhost
FRR_GRPC_PORT=50051
FRR_GRPC_TLS=false
JWT_SECRET=dev-secret-change-me
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
CORS_ALLOWED_ORIGINS=http://localhost:5173
LOG_LEVEL=debug
LOG_FORMAT=json
EOF

# Frontend
cd ../frontend
cat > .env.development << 'EOF'
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws
VITE_APP_NAME=FlintRoute
VITE_APP_VERSION=0.1.0-dev
VITE_ENABLE_DEBUG=true
EOF

cd ..
```

---

## Running Development Servers

### Install Dependencies

#### Backend Dependencies

```bash
cd backend

# Initialize Go module (if not exists)
go mod init github.com/padminisys/flintroute

# Install dependencies
go mod tidy

# Download dependencies
go mod download
```

#### Frontend Dependencies

```bash
cd frontend

# Install npm dependencies
npm install

# Or use yarn
yarn install
```

### Start Backend Server

```bash
cd backend

# Run with hot reload using air
go install github.com/cosmtrek/air@latest

# Create air configuration
cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF

# Start with hot reload
air

# Or run directly
go run cmd/server/main.go
```

### Start Frontend Server

```bash
cd frontend

# Start Vite dev server
npm run dev

# Or with yarn
yarn dev

# Server will start on http://localhost:5173
```

### Start Both Servers (Recommended)

Create a `Makefile` in project root:

```makefile
.PHONY: dev dev-backend dev-frontend install clean

# Start both backend and frontend in development mode
dev:
	@echo "Starting FlintRoute development servers..."
	@make -j2 dev-backend dev-frontend

# Start backend with hot reload
dev-backend:
	@echo "Starting backend server..."
	@cd backend && air || go run cmd/server/main.go

# Start frontend with hot reload
dev-frontend:
	@echo "Starting frontend server..."
	@cd frontend && npm run dev

# Install all dependencies
install:
	@echo "Installing backend dependencies..."
	@cd backend && go mod download
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf backend/tmp
	@rm -rf frontend/dist
	@rm -rf frontend/node_modules/.vite
```

Run both servers:

```bash
# Install dependencies first
make install

# Start development servers
make dev
```

### Access the Application

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api/v1
- **API Documentation**: http://localhost:8080/swagger (if enabled)

---

## Hot Reload Setup

### Backend Hot Reload (Air)

Already configured above. Air watches for file changes and automatically rebuilds.

**Features:**
- Automatic rebuild on `.go` file changes
- Fast incremental compilation
- Error logging to `build-errors.log`

### Frontend Hot Reload (Vite)

Vite provides built-in hot module replacement (HMR).

**Configure Vite** (`frontend/vite.config.ts`):

```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    host: true,
    hmr: {
      overlay: true
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true
      }
    }
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: true
  }
})
```

### Database Hot Reload

For SQLite development database:

```bash
# Watch for schema changes
cd backend
go install github.com/pressly/goose/v3/cmd/goose@latest

# Auto-migrate on changes
watch -n 5 'goose -dir migrations sqlite3 ./data/flintroute.db up'
```

---

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill process
sudo kill -9 <PID>

# Or change port in .env.development
```

#### Go Module Issues

```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify dependencies
go mod verify
```

#### npm Install Failures

```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules and package-lock.json
rm -rf node_modules package-lock.json

# Reinstall
npm install
```

#### Docker Permission Denied

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Logout and login again
# Or use newgrp
newgrp docker
```

#### FRR Connection Failed

```bash
# Check FRR is running
sudo systemctl status frr

# Check gRPC port is open
sudo netstat -tlnp | grep 50051

# Test gRPC connection
grpcurl -plaintext localhost:50051 list
```

### Debug Mode

#### Backend Debug

```bash
# Run with delve debugger
cd backend
dlv debug cmd/server/main.go

# Or attach to running process
dlv attach <PID>
```

#### Frontend Debug

```bash
# Enable React DevTools
# Install browser extension: React Developer Tools

# Enable Redux DevTools
# Install browser extension: Redux DevTools

# Enable verbose logging
VITE_ENABLE_DEBUG=true npm run dev
```

### Performance Issues

#### Backend Performance

```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.

# Profile memory usage
go test -memprofile=mem.prof -bench=.

# Analyze profiles
go tool pprof cpu.prof
```

#### Frontend Performance

```bash
# Analyze bundle size
npm run build
npm run analyze

# Check for memory leaks
# Use Chrome DevTools > Memory > Take Heap Snapshot
```

---

## Next Steps

1. **Set up FRR**: Follow [FRR Installation Guide](frr-installation.md)
2. **Run Tests**: Follow [Testing Guide](testing.md)
3. **Start Development**: Review [Phase 1 Implementation Plan](../phase1/implementation-plan.md)
4. **Contributing**: Read [Contributing Guidelines](../../CONTRIBUTING.md)

---

## Quick Reference

### Essential Commands

```bash
# Start development
make dev

# Run tests
make test

# Build for production
make build

# Clean artifacts
make clean

# Format code
make fmt

# Lint code
make lint
```

### Directory Navigation

```bash
# Backend
cd backend

# Frontend
cd frontend

# Documentation
cd docs

# Scripts
cd scripts
```

### Useful Aliases

Add to `~/.bashrc`:

```bash
# FlintRoute aliases
alias fr-dev='cd ~/flintroute && make dev'
alias fr-test='cd ~/flintroute && make test'
alias fr-backend='cd ~/flintroute/backend'
alias fr-frontend='cd ~/flintroute/frontend'
```

---

**Last Updated**: 2024-01-15  
**Version**: 0.1.0-alpha