# Phase 1 Implementation Plan

## Table of Contents
- [Phase 1 Implementation Plan](#phase-1-implementation-plan)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Goals](#goals)
    - [Methodology](#methodology)
  - [Timeline Summary](#timeline-summary)
    - [Milestones](#milestones)
  - [Week 1-2: Foundation](#week-1-2-foundation)
    - [Objectives](#objectives)
    - [Day-by-Day Breakdown](#day-by-day-breakdown)
      - [Week 1](#week-1)
      - [Week 2](#week-2)
    - [Week 1-2 Deliverables](#week-1-2-deliverables)
    - [Week 1-2 Testing](#week-1-2-testing)
  - [Week 3-4: BGP Peer Management](#week-3-4-bgp-peer-management)
    - [Objectives](#objectives-1)
    - [Backend Implementation](#backend-implementation)
    - [Frontend Implementation](#frontend-implementation)
    - [Week 3-4 Deliverables](#week-3-4-deliverables)
    - [Week 3-4 Testing](#week-3-4-testing)
  - [Week 5-6: Real-time Session Monitoring](#week-5-6-real-time-session-monitoring)
    - [Objectives](#objectives-2)
    - [Backend Implementation](#backend-implementation-1)
    - [Frontend Implementation](#frontend-implementation-1)
    - [Week 5-6 Deliverables](#week-5-6-deliverables)
  - [Week 7-8: Configuration Management](#week-7-8-configuration-management)
    - [Objectives](#objectives-3)
    - [Backend Implementation](#backend-implementation-2)
    - [Frontend Implementation](#frontend-implementation-2)
    - [Week 7-8 Deliverables](#week-7-8-deliverables)
  - [Week 9-10: Alerting System](#week-9-10-alerting-system)
    - [Objectives](#objectives-4)
    - [Backend Implementation](#backend-implementation-3)
    - [Frontend Implementation](#frontend-implementation-3)
    - [Week 9-10 Deliverables](#week-9-10-deliverables)
  - [Week 11-12: Topology Visualization \& Polish](#week-11-12-topology-visualization--polish)
    - [Objectives](#objectives-5)
    - [Implementation](#implementation)
    - [Week 11-12 Deliverables](#week-11-12-deliverables)
  - [Risk Management](#risk-management)
    - [Identified Risks](#identified-risks)
    - [Contingency Plans](#contingency-plans)
  - [Success Criteria](#success-criteria)
    - [Functional Requirements](#functional-requirements)
    - [Performance Requirements](#performance-requirements)
    - [Quality Requirements](#quality-requirements)
    - [Deployment Requirements](#deployment-requirements)
    - [User Acceptance](#user-acceptance)
  - [Next Steps](#next-steps)

---

## Overview

This document provides a detailed, week-by-week implementation plan for FlintRoute Phase 1 MVP. The plan is designed to deliver a production-ready BGP management system in 12 weeks.

### Goals

- **Functional**: Complete BGP peer management with real-time monitoring
- **Quality**: 80%+ test coverage, production-ready code
- **Timeline**: 12 weeks from start to MVP release
- **Team**: 2-3 developers (1-2 backend, 1 frontend)

### Methodology

- **Agile/Scrum**: 2-week sprints with weekly demos
- **Test-Driven**: Write tests before implementation
- **Continuous Integration**: Automated testing on every commit
- **Incremental Delivery**: Working features at end of each sprint

---

## Timeline Summary

```
┌─────────────────────────────────────────────────────────────┐
│                    12-Week Timeline                          │
├─────────────────────────────────────────────────────────────┤
│ Week 1-2   │ Foundation & Setup                             │
│ Week 3-4   │ BGP Peer Management                            │
│ Week 5-6   │ Real-time Session Monitoring                   │
│ Week 7-8   │ Configuration Management                       │
│ Week 9-10  │ Alerting System                                │
│ Week 11-12 │ Topology Visualization & Polish                │
└─────────────────────────────────────────────────────────────┘
```

### Milestones

| Week | Milestone | Deliverable |
|------|-----------|-------------|
| 2 | Foundation Complete | Dev environment, basic auth, FRR connection |
| 4 | Peer Management | CRUD operations for BGP peers |
| 6 | Session Monitoring | Real-time session status updates |
| 8 | Config Management | Backup/restore functionality |
| 10 | Alerting | Peer down detection and notifications |
| 12 | MVP Release | Production-ready v0.1.0 |

---

## Week 1-2: Foundation

### Objectives

- Set up complete development environment
- Establish project structure and tooling
- Implement basic authentication
- Create FRR gRPC client
- Set up CI/CD pipeline

### Day-by-Day Breakdown

#### Week 1

**Day 1-2: Project Initialization**

Backend Tasks:
```bash
# Initialize Go project
mkdir -p backend/{cmd/server,internal/{api,auth,config,frr,models,repository,service},pkg}
cd backend
go mod init github.com/padminisys/flintroute

# Install core dependencies
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get google.golang.org/grpc
```

Frontend Tasks:
```bash
# Initialize React project with Vite
npm create vite@latest frontend -- --template react-ts
cd frontend
npm install

# Install core dependencies
npm install react-router-dom @tanstack/react-query
npm install @reduxjs/toolkit react-redux
npm install axios
npm install @mui/material @emotion/react @emotion/styled
```

Deliverables:
- [ ] Repository structure created
- [ ] Go modules initialized
- [ ] React app initialized
- [ ] Dependencies installed
- [ ] README.md updated

**Day 3-4: Development Environment**

Tasks:
- [ ] Set up Docker Compose for FRR
- [ ] Configure hot reload (Air for Go, Vite for React)
- [ ] Create `.env` files for configuration
- [ ] Set up VSCode workspace settings
- [ ] Document setup process

Docker Compose (`docker-compose.dev.yml`):
```yaml
version: '3.8'

services:
  frr:
    image: frrouting/frr:v8.4.0
    container_name: flintroute-frr-dev
    privileged: true
    ports:
      - "50051:50051"  # gRPC
      - "2605:2605"    # vtysh
    volumes:
      - ./docker/frr/frr.conf:/etc/frr/frr.conf
    networks:
      - flintroute-net

  postgres:
    image: postgres:15-alpine
    container_name: flintroute-db-dev
    environment:
      POSTGRES_DB: flintroute
      POSTGRES_USER: flintroute
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - flintroute-net

networks:
  flintroute-net:
    driver: bridge

volumes:
  postgres_data:
```

Deliverables:
- [ ] Docker Compose configured
- [ ] FRR container running with gRPC enabled
- [ ] Hot reload working for both backend and frontend
- [ ] Environment variables documented

**Day 5: Database Schema Design**

Create database models:

```go
// internal/models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    Username     string `gorm:"uniqueIndex;not null" json:"username"`
    Email        string `gorm:"uniqueIndex;not null" json:"email"`
    PasswordHash string `gorm:"not null" json:"-"`
    Role         string `gorm:"not null;default:'operator'" json:"role"`
    Active       bool   `gorm:"not null;default:true" json:"active"`
}

// internal/models/bgp_peer.go
type BGPPeer struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    
    RemoteAS    uint32 `gorm:"not null" json:"remote_as"`
    RemoteAddr  string `gorm:"uniqueIndex;not null" json:"remote_addr"`
    Description string `json:"description"`
    Enabled     bool   `gorm:"not null;default:true" json:"enabled"`
    
    // Session state (populated from FRR)
    State           string    `gorm:"-" json:"state"`
    Uptime          int64     `gorm:"-" json:"uptime"`
    RoutesReceived  int       `gorm:"-" json:"routes_received"`
    RoutesAdvertised int      `gorm:"-" json:"routes_advertised"`
    LastUpdate      time.Time `gorm:"-" json:"last_update"`
}

// internal/models/config_version.go
type ConfigVersion struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    
    Version     int    `gorm:"not null" json:"version"`
    Config      string `gorm:"type:text;not null" json:"config"`
    Description string `json:"description"`
    CreatedBy   uint   `gorm:"not null" json:"created_by"`
    User        User   `gorm:"foreignKey:CreatedBy" json:"user"`
}

// internal/models/audit_log.go
type AuditLog struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `json:"created_at"`
    
    UserID    uint   `gorm:"not null" json:"user_id"`
    User      User   `gorm:"foreignKey:UserID" json:"user"`
    Action    string `gorm:"not null" json:"action"`
    Resource  string `gorm:"not null" json:"resource"`
    ResourceID uint  `json:"resource_id"`
    Details   string `gorm:"type:text" json:"details"`
    IPAddress string `json:"ip_address"`
}
```

Deliverables:
- [ ] Database models defined
- [ ] Migrations created
- [ ] Database seeding script for development

#### Week 2

**Day 6-7: Authentication System**

Implement JWT-based authentication:

```go
// internal/auth/jwt.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   uint     `json:"user_id"`
    Username string   `json:"username"`
    Role     string   `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(getJWTSecret()))
}

// internal/api/auth_handler.go
package api

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"`
}

func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // Authenticate user
    user, err := h.authService.Authenticate(req.Username, req.Password)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Generate tokens
    accessToken, _ := auth.GenerateToken(user.ID, user.Username, user.Role)
    refreshToken, _ := auth.GenerateRefreshToken(user.ID)
    
    c.JSON(200, LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    900, // 15 minutes
    })
}
```

Frontend authentication:

```typescript
// src/services/auth.ts
import axios from 'axios'

interface LoginCredentials {
  username: string
  password: string
}

interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

export const authService = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    const response = await axios.post('/api/v1/auth/login', credentials)
    return response.data
  },

  async logout(): Promise<void> {
    await axios.post('/api/v1/auth/logout')
    localStorage.removeItem('access_token')
  },

  async refreshToken(): Promise<string> {
    const response = await axios.post('/api/v1/auth/refresh')
    return response.data.access_token
  },
}

// src/store/authSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { authService } from '@/services/auth'

export const login = createAsyncThunk(
  'auth/login',
  async (credentials: LoginCredentials) => {
    const response = await authService.login(credentials)
    localStorage.setItem('access_token', response.access_token)
    return response
  }
)

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: null,
    token: null,
    isAuthenticated: false,
    loading: false,
    error: null,
  },
  reducers: {
    logout: (state) => {
      state.user = null
      state.token = null
      state.isAuthenticated = false
      localStorage.removeItem('access_token')
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false
        state.token = action.payload.access_token
        state.isAuthenticated = true
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false
        state.error = action.error.message
      })
  },
})
```

Deliverables:
- [ ] JWT token generation and validation
- [ ] Login/logout endpoints
- [ ] Password hashing with bcrypt
- [ ] Frontend login page
- [ ] Token storage and refresh logic
- [ ] Protected routes

**Day 8-9: FRR gRPC Client**

Implement FRR gRPC client:

```go
// internal/frr/client.go
package frr

import (
    "context"
    "fmt"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type Client struct {
    conn   *grpc.ClientConn
    client pb.NorthboundClient
}

func NewClient(host string, port int) (*Client, error) {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to FRR: %w", err)
    }
    
    return &Client{
        conn:   conn,
        client: pb.NewNorthboundClient(conn),
    }, nil
}

func (c *Client) Close() error {
    return c.conn.Close()
}

func (c *Client) GetConfig(ctx context.Context) (string, error) {
    req := &pb.GetRequest{
        Type:     pb.GetRequest_STATE,
        Encoding: pb.Encoding_JSON,
    }
    
    resp, err := c.client.Get(ctx, req)
    if err != nil {
        return "", fmt.Errorf("failed to get config: %w", err)
    }
    
    return string(resp.Data.Data), nil
}

func (c *Client) CreateBGPPeer(ctx context.Context, peer *models.BGPPeer) error {
    config := fmt.Sprintf(`{
        "frr-routing:routing": {
            "control-plane-protocols": {
                "control-plane-protocol": [{
                    "type": "frr-bgp:bgp",
                    "name": "main",
                    "frr-bgp:bgp": {
                        "neighbors": {
                            "neighbor": [{
                                "remote-address": "%s",
                                "remote-as": %d,
                                "description": "%s"
                            }]
                        }
                    }
                }]
            }
        }
    }`, peer.RemoteAddr, peer.RemoteAS, peer.Description)
    
    req := &pb.CommitRequest{
        Type:     pb.CommitRequest_CANDIDATE,
        Encoding: pb.Encoding_JSON,
        Config:   []byte(config),
    }
    
    _, err := c.client.Commit(ctx, req)
    return err
}
```

Deliverables:
- [ ] gRPC client implementation
- [ ] Connection pooling
- [ ] Error handling and retries
- [ ] Unit tests for client
- [ ] Integration tests with real FRR

**Day 10: CI/CD Pipeline**

Set up GitHub Actions:

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  backend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Run tests
        working-directory: ./backend
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out

  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Run tests
        working-directory: ./frontend
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./frontend/coverage/coverage-final.json

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          working-directory: backend
      
      - name: ESLint
        working-directory: ./frontend
        run: |
          npm ci
          npm run lint
```

Deliverables:
- [ ] GitHub Actions workflow configured
- [ ] Automated testing on PR
- [ ] Code coverage reporting
- [ ] Linting checks

### Week 1-2 Deliverables

- [x] Complete development environment
- [x] Project structure established
- [x] Authentication system working
- [x] FRR gRPC client functional
- [x] CI/CD pipeline operational
- [x] Documentation updated

### Week 1-2 Testing

```bash
# Test authentication
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Test FRR connection
go test ./internal/frr/... -v

# Test frontend
cd frontend && npm test
```

---

## Week 3-4: BGP Peer Management

### Objectives

- Implement complete CRUD operations for BGP peers
- Integrate with FRR gRPC API
- Build frontend peer management UI
- Add form validation and error handling

### Backend Implementation

**Day 11-12: BGP Peer Service**

```go
// internal/service/bgp_service.go
package service

type BGPService struct {
    repo      repository.BGPPeerRepository
    frrClient *frr.Client
    audit     *audit.Logger
}

func (s *BGPService) CreatePeer(ctx context.Context, peer *models.BGPPeer) error {
    // Validate peer
    if err := s.validatePeer(peer); err != nil {
        return err
    }
    
    // Check for duplicates
    existing, _ := s.repo.FindByAddress(ctx, peer.RemoteAddr)
    if existing != nil {
        return ErrPeerAlreadyExists
    }
    
    // Create in FRR
    if err := s.frrClient.CreateBGPPeer(ctx, peer); err != nil {
        return fmt.Errorf("failed to create peer in FRR: %w", err)
    }
    
    // Save to database
    if err := s.repo.Create(ctx, peer); err != nil {
        // Rollback FRR changes
        s.frrClient.DeleteBGPPeer(ctx, peer.RemoteAddr)
        return err
    }
    
    // Audit log
    s.audit.Log(ctx, "bgp.peer.create", peer)
    
    return nil
}

func (s *BGPService) UpdatePeer(ctx context.Context, id uint, updates *models.BGPPeer) error {
    // Get existing peer
    peer, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Update in FRR
    if err := s.frrClient.UpdateBGPPeer(ctx, updates); err != nil {
        return err
    }
    
    // Update in database
    if err := s.repo.Update(ctx, id, updates); err != nil {
        // Rollback FRR changes
        s.frrClient.UpdateBGPPeer(ctx, peer)
        return err
    }
    
    s.audit.Log(ctx, "bgp.peer.update", updates)
    return nil
}

func (s *BGPService) DeletePeer(ctx context.Context, id uint) error {
    peer, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Delete from FRR
    if err := s.frrClient.DeleteBGPPeer(ctx, peer.RemoteAddr); err != nil {
        return err
    }
    
    // Delete from database
    if err := s.repo.Delete(ctx, id); err != nil {
        return err
    }
    
    s.audit.Log(ctx, "bgp.peer.delete", peer)
    return nil
}

func (s *BGPService) ListPeers(ctx context.Context) ([]*models.BGPPeer, error) {
    return s.repo.FindAll(ctx)
}

func (s *BGPService) GetPeer(ctx context.Context, id uint) (*models.BGPPeer, error) {
    return s.repo.FindByID(ctx, id)
}
```

**Day 13-14: API Endpoints**

```go
// internal/api/bgp_handler.go
package api

func (h *Handler) RegisterBGPRoutes(r *gin.RouterGroup) {
    bgp := r.Group("/bgp")
    bgp.Use(h.authMiddleware.RequireAuth())
    
    bgp.GET("/peers", h.ListPeers)
    bgp.GET("/peers/:id", h.GetPeer)
    bgp.POST("/peers", h.authMiddleware.RequirePermission("bgp:write"), h.CreatePeer)
    bgp.PUT("/peers/:id", h.authMiddleware.RequirePermission("bgp:write"), h.UpdatePeer)
    bgp.DELETE("/peers/:id", h.authMiddleware.RequirePermission("bgp:delete"), h.DeletePeer)
    bgp.POST("/peers/:id/enable", h.authMiddleware.RequirePermission("bgp:write"), h.EnablePeer)
    bgp.POST("/peers/:id/disable", h.authMiddleware.RequirePermission("bgp:write"), h.DisablePeer)
}

func (h *Handler) CreatePeer(c *gin.Context) {
    var peer models.BGPPeer
    if err := c.ShouldBindJSON(&peer); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    if err := h.bgpService.CreatePeer(c.Request.Context(), &peer); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, peer)
}
```

### Frontend Implementation

**Day 15-16: Peer List Component**

```typescript
// src/pages/BGPPeers.tsx
import { useQuery } from '@tanstack/react-query'
import { DataGrid, GridColDef } from '@mui/x-data-grid'
import { Button, Chip } from '@mui/material'
import { Add as AddIcon } from '@mui/icons-material'

const columns: GridColDef[] = [
  { field: 'id', headerName: 'ID', width: 70 },
  { field: 'remote_addr', headerName: 'Remote Address', width: 150 },
  { field: 'remote_as', headerName: 'AS Number', width: 120 },
  { field: 'description', headerName: 'Description', width: 200 },
  {
    field: 'state',
    headerName: 'State',
    width: 120,
    renderCell: (params) => (
      <Chip
        label={params.value}
        color={params.value === 'Established' ? 'success' : 'warning'}
        size="small"
      />
    ),
  },
  {
    field: 'enabled',
    headerName: 'Status',
    width: 100,
    renderCell: (params) => (
      <Chip
        label={params.value ? 'Enabled' : 'Disabled'}
        color={params.value ? 'success' : 'default'}
        size="small"
      />
    ),
  },
]

export function BGPPeers() {
  const { data: peers, isLoading } = useQuery({
    queryKey: ['bgp', 'peers'],
    queryFn: () => api.getBGPPeers(),
  })

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" mb={2}>
        <Typography variant="h4">BGP Peers</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => navigate('/bgp/peers/new')}
        >
          Add Peer
        </Button>
      </Box>

      <DataGrid
        rows={peers || []}
        columns={columns}
        loading={isLoading}
        autoHeight
        pageSizeOptions={[10, 25, 50]}
      />
    </Box>
  )
}
```

**Day 17-18: Peer Form Component**

```typescript
// src/components/BGPPeerForm.tsx
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const peerSchema = z.object({
  remote_as: z.number().min(1).max(4294967295),
  remote_addr: z.string().ip(),
  description: z.string().optional(),
  enabled: z.boolean().default(true),
})

type PeerFormData = z.infer<typeof peerSchema>

export function BGPPeerForm({ peer, onSubmit }: Props) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<PeerFormData>({
    resolver: zodResolver(peerSchema),
    defaultValues: peer,
  })

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <TextField
        label="Remote AS"
        type="number"
        {...register('remote_as', { valueAsNumber: true })}
        error={!!errors.remote_as}
        helperText={errors.remote_as?.message}
        fullWidth
        margin="normal"
      />

      <TextField
        label="Remote Address"
        {...register('remote_addr')}
        error={!!errors.remote_addr}
        helperText={errors.remote_addr?.message}
        fullWidth
        margin="normal"
      />

      <TextField
        label="Description"
        {...register('description')}
        error={!!errors.description}
        helperText={errors.description?.message}
        fullWidth
        margin="normal"
        multiline
        rows={3}
      />

      <FormControlLabel
        control={<Switch {...register('enabled')} />}
        label="Enabled"
      />

      <Box mt={2}>
        <Button
          type="submit"
          variant="contained"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Saving...' : 'Save'}
        </Button>
      </Box>
    </form>
  )
}
```

### Week 3-4 Deliverables

- [x] Complete CRUD API for BGP peers
- [x] FRR integration for peer management
- [x] Peer list view with filtering/sorting
- [x] Peer creation form with validation
- [x] Peer edit functionality
- [x] Peer deletion with confirmation
- [x] Enable/disable peer toggle
- [x] Unit tests for service layer
- [x] Integration tests with FRR
- [x] E2E tests for peer management

### Week 3-4 Testing

```bash
# Backend tests
cd backend
go test ./internal/service/... -v
go test ./internal/api/... -v

# Frontend tests
cd frontend
npm test -- BGPPeerForm
npm test -- BGPPeers

# E2E tests
npx playwright test bgp-peer-management
```

---

## Week 5-6: Real-time Session Monitoring

### Objectives

- Implement WebSocket server for real-time updates
- Poll BGP session state from FRR
- Display session statistics in frontend
- Show routes received/advertised counters
- Real-time session state changes

### Backend Implementation

**Day 19-20: WebSocket Server**

```go
// internal/websocket/hub.go
package websocket

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// internal/websocket/client.go
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    for {
        _, message
, _, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        // Handle incoming messages if needed
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(30 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.conn.WriteMessage(websocket.TextMessage, message)

        case <-ticker.C:
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**Day 21-22: Session Monitoring Service**

```go
// internal/service/monitor_service.go
package service

type MonitorService struct {
    frrClient *frr.Client
    hub       *websocket.Hub
    repo      repository.BGPPeerRepository
    interval  time.Duration
}

func (s *MonitorService) Start(ctx context.Context) {
    ticker := time.NewTicker(s.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.pollSessions(ctx)
        }
    }
}

func (s *MonitorService) pollSessions(ctx context.Context) {
    // Get all peers from database
    peers, err := s.repo.FindAll(ctx)
    if err != nil {
        log.Error("Failed to get peers", "error", err)
        return
    }

    // Get session status from FRR
    sessions, err := s.frrClient.GetBGPSummary(ctx)
    if err != nil {
        log.Error("Failed to get BGP summary", "error", err)
        return
    }

    // Update peer states
    updates := make([]SessionUpdate, 0)
    for _, peer := range peers {
        if session, ok := sessions[peer.RemoteAddr]; ok {
            peer.State = session.State
            peer.Uptime = session.Uptime
            peer.RoutesReceived = session.RoutesReceived
            peer.RoutesAdvertised = session.RoutesAdvertised
            peer.LastUpdate = time.Now()

            updates = append(updates, SessionUpdate{
                PeerID:           peer.ID,
                State:            peer.State,
                RoutesReceived:   peer.RoutesReceived,
                RoutesAdvertised: peer.RoutesAdvertised,
            })
        }
    }

    // Broadcast updates via WebSocket
    if len(updates) > 0 {
        data, _ := json.Marshal(updates)
        s.hub.Broadcast(data)
    }
}
```

### Frontend Implementation

**Day 23-24: WebSocket Client**

```typescript
// src/services/websocket.ts
import { useEffect, useRef } from 'react'
import { useDispatch } from 'react-redux'

export function useWebSocket(url: string) {
  const ws = useRef<WebSocket | null>(null)
  const dispatch = useDispatch()

  useEffect(() => {
    ws.current = new WebSocket(url)

    ws.current.onopen = () => {
      console.log('WebSocket connected')
    }

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data)
      dispatch(updateSessionStates(data))
    }

    ws.current.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    ws.current.onclose = () => {
      console.log('WebSocket disconnected')
      // Reconnect after 5 seconds
      setTimeout(() => {
        ws.current = new WebSocket(url)
      }, 5000)
    }

    return () => {
      ws.current?.close()
    }
  }, [url, dispatch])

  return ws.current
}

// src/pages/SessionMonitor.tsx
export function SessionMonitor() {
  const sessions = useSelector((state) => state.sessions.data)
  useWebSocket('ws://localhost:8080/ws')

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        BGP Session Monitor
      </Typography>

      <Grid container spacing={3}>
        {sessions.map((session) => (
          <Grid item xs={12} md={6} lg={4} key={session.id}>
            <Card>
              <CardContent>
                <Box display="flex" justifyContent="space-between" alignItems="center">
                  <Typography variant="h6">{session.remote_addr}</Typography>
                  <Chip
                    label={session.state}
                    color={session.state === 'Established' ? 'success' : 'error'}
                  />
                </Box>

                <Typography color="textSecondary" gutterBottom>
                  AS {session.remote_as}
                </Typography>

                <Divider sx={{ my: 2 }} />

                <Grid container spacing={2}>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="textSecondary">
                      Routes Received
                    </Typography>
                    <Typography variant="h6">
                      {session.routes_received.toLocaleString()}
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body2" color="textSecondary">
                      Routes Advertised
                    </Typography>
                    <Typography variant="h6">
                      {session.routes_advertised.toLocaleString()}
                    </Typography>
                  </Grid>
                </Grid>

                <Box mt={2}>
                  <Typography variant="body2" color="textSecondary">
                    Uptime: {formatUptime(session.uptime)}
                  </Typography>
                  <Typography variant="caption" color="textSecondary">
                    Last update: {formatTime(session.last_update)}
                  </Typography>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  )
}
```

### Week 5-6 Deliverables

- [x] WebSocket server implementation
- [x] Session state polling from FRR
- [x] Real-time updates to connected clients
- [x] Session statistics display
- [x] Routes received/advertised counters
- [x] Session uptime tracking
- [x] Automatic reconnection on disconnect
- [x] Performance testing with multiple clients

---

## Week 7-8: Configuration Management

### Objectives

- Implement configuration backup system
- Version control for configurations
- Restore functionality with rollback
- Configuration diff viewer
- Atomic transaction support

### Backend Implementation

**Day 25-26: Configuration Backup Service**

```go
// internal/service/config_service.go
package service

type ConfigService struct {
    frrClient *frr.Client
    repo      repository.ConfigRepository
    audit     *audit.Logger
}

func (s *ConfigService) CreateBackup(ctx context.Context, description string, userID uint) (*models.ConfigVersion, error) {
    // Get current FRR configuration
    config, err := s.frrClient.GetConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get FRR config: %w", err)
    }

    // Get next version number
    latestVersion, _ := s.repo.GetLatestVersion(ctx)
    nextVersion := 1
    if latestVersion != nil {
        nextVersion = latestVersion.Version + 1
    }

    // Create backup
    backup := &models.ConfigVersion{
        Version:     nextVersion,
        Config:      config,
        Description: description,
        CreatedBy:   userID,
    }

    if err := s.repo.Create(ctx, backup); err != nil {
        return nil, err
    }

    s.audit.Log(ctx, "config.backup.create", backup)
    return backup, nil
}

func (s *ConfigService) RestoreBackup(ctx context.Context, versionID uint, userID uint) error {
    // Get backup
    backup, err := s.repo.FindByID(ctx, versionID)
    if err != nil {
        return err
    }

    // Create backup of current config before restore
    _, err = s.CreateBackup(ctx, "Auto-backup before restore", userID)
    if err != nil {
        return fmt.Errorf("failed to create pre-restore backup: %w", err)
    }

    // Apply configuration to FRR
    if err := s.frrClient.SetConfig(ctx, backup.Config); err != nil {
        return fmt.Errorf("failed to restore config: %w", err)
    }

    s.audit.Log(ctx, "config.backup.restore", backup)
    return nil
}

func (s *ConfigService) GetDiff(ctx context.Context, version1ID, version2ID uint) (string, error) {
    v1, err := s.repo.FindByID(ctx, version1ID)
    if err != nil {
        return "", err
    }

    v2, err := s.repo.FindByID(ctx, version2ID)
    if err != nil {
        return "", err
    }

    // Generate diff
    diff := difflib.UnifiedDiff{
        A:        difflib.SplitLines(v1.Config),
        B:        difflib.SplitLines(v2.Config),
        FromFile: fmt.Sprintf("Version %d", v1.Version),
        ToFile:   fmt.Sprintf("Version %d", v2.Version),
        Context:  3,
    }

    return difflib.GetUnifiedDiffString(diff)
}

func (s *ConfigService) ListBackups(ctx context.Context) ([]*models.ConfigVersion, error) {
    return s.repo.FindAll(ctx)
}
```

**Day 27-28: Atomic Transactions**

```go
// internal/service/transaction.go
package service

type Transaction struct {
    id            string
    changes       []Change
    originalState string
    frrClient     *frr.Client
}

func (s *ConfigService) BeginTransaction(ctx context.Context) (*Transaction, error) {
    // Get current configuration
    config, err := s.frrClient.GetConfig(ctx)
    if err != nil {
        return nil, err
    }

    return &Transaction{
        id:            uuid.New().String(),
        changes:       make([]Change, 0),
        originalState: config,
        frrClient:     s.frrClient,
    }, nil
}

func (t *Transaction) AddChange(change Change) {
    t.changes = append(t.changes, change)
}

func (t *Transaction) Commit(ctx context.Context) error {
    // Apply all changes
    for _, change := range t.changes {
        if err := change.Apply(ctx, t.frrClient); err != nil {
            // Rollback on error
            t.Rollback(ctx)
            return fmt.Errorf("transaction failed, rolled back: %w", err)
        }
    }

    // Verify configuration is valid
    if err := t.frrClient.ValidateConfig(ctx); err != nil {
        t.Rollback(ctx)
        return fmt.Errorf("config validation failed, rolled back: %w", err)
    }

    return nil
}

func (t *Transaction) Rollback(ctx context.Context) error {
    return t.frrClient.SetConfig(ctx, t.originalState)
}
```

### Frontend Implementation

**Day 29-30: Configuration Management UI**

```typescript
// src/pages/ConfigurationManagement.tsx
export function ConfigurationManagement() {
  const { data: backups, isLoading } = useQuery({
    queryKey: ['config', 'backups'],
    queryFn: () => api.getConfigBackups(),
  })

  const createBackupMutation = useMutation({
    mutationFn: (description: string) => api.createConfigBackup(description),
    onSuccess: () => {
      queryClient.invalidateQueries(['config', 'backups'])
      toast.success('Backup created successfully')
    },
  })

  const restoreMutation = useMutation({
    mutationFn: (versionId: number) => api.restoreConfigBackup(versionId),
    onSuccess: () => {
      toast.success('Configuration restored successfully')
    },
  })

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" mb={3}>
        <Typography variant="h4">Configuration Management</Typography>
        <Button
          variant="contained"
          startIcon={<BackupIcon />}
          onClick={() => setBackupDialogOpen(true)}
        >
          Create Backup
        </Button>
      </Box>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Version</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Created By</TableCell>
              <TableCell>Created At</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {backups?.map((backup) => (
              <TableRow key={backup.id}>
                <TableCell>{backup.version}</TableCell>
                <TableCell>{backup.description}</TableCell>
                <TableCell>{backup.user.username}</TableCell>
                <TableCell>{formatDate(backup.created_at)}</TableCell>
                <TableCell>
                  <IconButton
                    onClick={() => handleRestore(backup.id)}
                    title="Restore"
                  >
                    <RestoreIcon />
                  </IconButton>
                  <IconButton
                    onClick={() => handleViewDiff(backup.id)}
                    title="View Diff"
                  >
                    <DiffIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <ConfigDiffDialog
        open={diffDialogOpen}
        onClose={() => setDiffDialogOpen(false)}
        version1={selectedVersion1}
        version2={selectedVersion2}
      />
    </Box>
  )
}
```

### Week 7-8 Deliverables

- [x] Configuration backup creation
- [x] Configuration restore functionality
- [x] Version history listing
- [x] Configuration diff viewer
- [x] Atomic transaction support
- [x] Automatic pre-restore backup
- [x] Configuration validation
- [x] Rollback on errors

---

## Week 9-10: Alerting System

### Objectives

- Implement alert model and storage
- Peer down detection
- Notification system (email/webhook)
- Alert history and acknowledgment
- Alert configuration UI

### Backend Implementation

**Day 31-32: Alert Service**

```go
// internal/service/alert_service.go
package service

type AlertService struct {
    repo         repository.AlertRepository
    notifier     *notification.Service
    monitorSvc   *MonitorService
}

func (s *AlertService) Start(ctx context.Context) {
    // Subscribe to session state changes
    s.monitorSvc.Subscribe(func(update SessionUpdate) {
        s.handleSessionUpdate(ctx, update)
    })
}

func (s *AlertService) handleSessionUpdate(ctx context.Context, update SessionUpdate) {
    peer, _ := s.peerRepo.FindByID(ctx, update.PeerID)
    
    // Check for peer down
    if update.State != "Established" && update.PreviousState == "Established" {
        alert := &models.Alert{
            Type:       "peer_down",
            Severity:   "critical",
            PeerID:     peer.ID,
            Message:    fmt.Sprintf("BGP peer %s (AS%d) is down", peer.RemoteAddr, peer.RemoteAS),
            Details:    fmt.Sprintf("State changed from %s to %s", update.PreviousState, update.State),
            Timestamp:  time.Now(),
        }
        
        s.repo.Create(ctx, alert)
        s.notifier.Send(ctx, alert)
    }
    
    // Check for peer up
    if update.State == "Established" && update.PreviousState != "Established" {
        alert := &models.Alert{
            Type:       "peer_up",
            Severity:   "info",
            PeerID:     peer.ID,
            Message:    fmt.Sprintf("BGP peer %s (AS%d) is up", peer.RemoteAddr, peer.RemoteAS),
            Timestamp:  time.Now(),
        }
        
        s.repo.Create(ctx, alert)
        s.notifier.Send(ctx, alert)
    }
}

func (s *AlertService) AcknowledgeAlert(ctx context.Context, alertID uint, userID uint) error {
    alert, err := s.repo.FindByID(ctx, alertID)
    if err != nil {
        return err
    }
    
    alert.Acknowledged = true
    alert.AcknowledgedBy = userID
    alert.AcknowledgedAt = time.Now()
    
    return s.repo.Update(ctx, alert)
}
```

**Day 33-34: Notification Service**

```go
// internal/notification/service.go
package notification

type Service struct {
    emailSender   *EmailSender
    webhookSender *WebhookSender
    config        *Config
}

func (s *Service) Send(ctx context.Context, alert *models.Alert) error {
    var wg sync.WaitGroup
    errors := make(chan error, 2)
    
    // Send email notification
    if s.config.EmailEnabled {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := s.emailSender.Send(ctx, alert); err != nil {
                errors <- fmt.Errorf("email failed: %w", err)
            }
        }()
    }
    
    // Send webhook notification
    if s.config.WebhookEnabled {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := s.webhookSender.Send(ctx, alert); err != nil {
                errors <- fmt.Errorf("webhook failed: %w", err)
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    // Collect errors
    var errs []error
    for err := range errors {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return fmt.Errorf("notification errors: %v", errs)
    }
    
    return nil
}
```

### Frontend Implementation

**Day 35-36: Alert Dashboard**

```typescript
// src/pages/Alerts.tsx
export function Alerts() {
  const { data: alerts } = useQuery({
    queryKey: ['alerts'],
    queryFn: () => api.getAlerts(),
    refetchInterval: 30000, // Refresh every 30 seconds
  })

  const acknowledgeMutation = useMutation({
    mutationFn: (alertId: number) => api.acknowledgeAlert(alertId),
    onSuccess: () => {
      queryClient.invalidateQueries(['alerts'])
    },
  })

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'error'
      case 'warning': return 'warning'
      case 'info': return 'info'
      default: return 'default'
    }
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Alerts
      </Typography>

      <List>
        {alerts?.map((alert) => (
          <ListItem
            key={alert.id}
            sx={{
              bgcolor: alert.acknowledged ? 'transparent' : 'action.hover',
              mb: 1,
              borderRadius: 1,
            }}
          >
            <ListItemIcon>
              <Chip
                label={alert.severity}
                color={getSeverityColor(alert.severity)}
                size="small"
              />
            </ListItemIcon>
            <ListItemText
              primary={alert.message}
              secondary={
                <>
                  <Typography variant="body2" component="span">
                    {alert.details}
                  </Typography>
                  <br />
                  <Typography variant="caption" color="textSecondary">
                    {formatDate(alert.timestamp)}
                  </Typography>
                </>
              }
            />
            {!alert.acknowledged && (
              <ListItemSecondaryAction>
                <Button
                  size="small"
                  onClick={() => acknowledgeMutation.mutate(alert.id)}
                >
                  Acknowledge
                </Button>
              </ListItemSecondaryAction>
            )}
          </ListItem>
        ))}
      </List>
    </Box>
  )
}
```

### Week 9-10 Deliverables

- [x] Alert model and database schema
- [x] Peer down/up detection
- [x] Email notification system
- [x] Webhook notification system
- [x] Alert history storage
- [x] Alert acknowledgment
- [x] Alert dashboard UI
- [x] Alert filtering and search

---

## Week 11-12: Topology Visualization & Polish

### Objectives

- Interactive BGP topology map
- Dashboard with statistics
- Bug fixes and refinements
- Documentation updates
- Testing and validation
- Performance optimization

### Implementation

**Day 37-38: Topology Visualization**

```typescript
// src/components/TopologyMap.tsx
import ReactFlow, { Node, Edge } from 'reactflow'
import 'reactflow/dist/style.css'

export function TopologyMap() {
  const { data: peers } = useQuery({
    queryKey: ['bgp', 'peers'],
    queryFn: () => api.getBGPPeers(),
  })

  const nodes: Node[] = useMemo(() => {
    const localNode: Node = {
      id: 'local',
      type: 'custom',
      position: { x: 400, y: 300 },
      data: {
        label: 'Local Router',
        as: 65000,
        type: 'local',
      },
    }

    const peerNodes: Node[] = peers?.map((peer, index) => ({
      id: `peer-${peer.id}`,
      type: 'custom',
      position: {
        x: 400 + 300 * Math.cos((2 * Math.PI * index) / peers.length),
        y: 300 + 300 * Math.sin((2 * Math.PI * index) / peers.length),
      },
      data: {
        label: peer.remote_addr,
        as: peer.remote_as,
        state: peer.state,
        type: 'peer',
      },
    })) || []

    return [localNode, ...peerNodes]
  }, [peers])

  const edges: Edge[] = useMemo(() => {
    return peers?.map((peer) => ({
      id: `edge-${peer.id}`,
      source: 'local',
      target: `peer-${peer.id}`,
      animated: peer.state === 'Established',
      style: {
        stroke: peer.state === 'Established' ? '#4caf50' : '#f44336',
      },
    })) || []
  }, [peers])

  return (
    <Box sx={{ height: '600px' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        nodeTypes={nodeTypes}
      />
    </Box>
  )
}
```

**Day 39-40: Dashboard**

```typescript
// src/pages/Dashboard.tsx
export function Dashboard() {
  const { data: stats } = useQuery({
    queryKey: ['dashboard', 'stats'],
    queryFn: () => api.getDashboardStats(),
  })

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      <Grid container spacing={3}>
        <Grid item xs={12} md={3}>
          <StatCard
            title="Total Peers"
            value={stats?.total_peers || 0}
            icon={<PeersIcon />}
            color="primary"
          />
        </Grid>
        <Grid item xs={12} md={3}>
          <StatCard
            title="Established"
            value={stats?.established_peers || 0}
            icon={<CheckIcon />}
            color="success"
          />
        </Grid>
        <Grid item xs={12} md={3}>
          <StatCard
            title="Down"
            value={stats?.down_peers || 0}
            icon={<ErrorIcon />}
            color="error"
          />
        </Grid>
        <Grid item xs={12} md={3}>
          <StatCard
            title="Total Routes"
            value={stats?.total_routes || 0}
            icon={<RoutesIcon />}
            color="info"
          />
        </Grid>

        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              BGP Topology
            </Typography>
            <TopologyMap />
          </Paper>
        </Grid>

        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Recent Alerts
            </Typography>
            <RecentAlerts limit={5} />
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Paper sx={{ p: 2 }}>
            <Typography variant="h6" gutterBottom>
              Session Status
            </Typography>
            <SessionStatusChart />
          </Paper>
        </Grid>
      </Grid>
    </Box>
  )
}
```

**Day 41-42: Bug Fixes & Refinements**

Focus areas:
- [ ] Fix any reported bugs
- [ ] Improve error handling
- [ ] Optimize database queries
- [ ] Improve UI responsiveness
- [ ] Add loading states
- [ ] Improve error messages
- [ ] Add input validation
- [ ] Fix edge cases

**Day 43-44: Documentation & Testing**

- [ ] Update API documentation
- [ ] Update user guide
- [ ] Update deployment guide
- [ ] Run full test suite
- [ ] Performance testing
- [ ] Security audit
- [ ] Accessibility testing

**Day 45-46: Final Testing & Release Prep**

- [ ] End-to-end testing
- [ ] User acceptance testing
- [ ] Performance benchmarks
- [ ] Security scan
- [ ] Create release notes
- [ ] Tag v0.1.0 release
- [ ] Deploy to staging
- [ ] Final review

### Week 11-12 Deliverables

- [x] Interactive topology visualization
- [x] Comprehensive dashboard
- [x] All bugs fixed
- [x] Documentation complete
- [x] All tests passing
- [x] Performance optimized
- [x] Security validated
- [x] v0.1.0 released

---

## Risk Management

### Identified Risks

| Risk | Probability | Impact | Mitigation Strategy |
|------|------------|--------|---------------------|
| FRR API instability | Medium | High | Pin FRR version, extensive testing, fallback mechanisms |
| Performance issues | Medium | Medium | Early performance testing, optimization sprints |
| Security vulnerabilities | Low | Critical | Security audit, penetration testing, code review |
| Scope creep | High | Medium | Strict feature freeze after Week 8, prioritization |
| Team availability | Medium | High | Buffer time in schedule, cross-training |
| Integration complexity | Medium | High | Early prototyping, incremental integration |
| Testing delays | Medium | Medium | Parallel testing, automated CI/CD |

### Contingency Plans

**If FRR integration blocked:**
- Use mock FRR responses for development
- Implement simulator mode
- Continue frontend development independently

**If performance inadequate:**
- Implement caching layer (Redis)
- Optimize database queries and indexes
- Add pagination and lazy loading
- Profile and optimize hot paths

**If security issues found:**
- Delay release for critical issues
- Hot-fix for non-critical issues
- Implement security advisory process
- Conduct additional security review

**If timeline slips:**
- Reduce scope (defer non-critical features to Phase 2)
- Add resources if available
- Extend timeline by maximum 2 weeks
- Prioritize MVP features only

---

## Success Criteria

### Functional Requirements

- [ ] Users can create, read, update, and delete BGP peers
- [ ] Real-time session monitoring with WebSocket updates
- [ ] Configuration backup and restore functionality
- [ ] Alert system for peer down/up events
- [ ] Interactive topology visualization
- [ ] Dashboard with key statistics

### Performance Requirements

- [ ] API response time < 200ms (p95)
- [ ] UI page load < 2s
- [ ] Support 100+ concurrent BGP peers
- [ ] WebSocket latency < 100ms
- [ ] Handle 1000+ concurrent WebSocket connections

### Quality Requirements

- [ ] Backend test coverage > 80%
- [ ] Frontend test coverage > 70%
- [ ] Zero critical security vulnerabilities
- [ ] All documentation complete and accurate
- [ ] Passes accessibility standards (WCAG 2.1 AA)

### Deployment Requirements

- [ ] Docker images published
- [ ] Installation documentation complete
- [ ] Deployment tested on Debian 12
- [ ] Systemd service configured
- [ ] Monitoring and logging operational

### User Acceptance

- [ ] 5+ beta testers provide positive feedback
- [ ] All critical bugs resolved
- [ ] User guide reviewed and approved
- [ ] Demo successfully presented to stakeholders

---

## Next Steps

After Phase 1 completion:

1. **Gather Feedback**: Collect user feedback and bug reports
2. **Plan Phase 2**: Define Phase 2 features (OSPF, static routes, advanced policies)
3. **Continuous Improvement**: Address technical debt and optimize performance
4. **Community Building**: Engage with open-source community
5. **Documentation**: Create video tutorials and advanced guides

---

**Last Updated**: 2024-01-15  
**Version**: 0.1.0-alpha  
**Status**: Planning Phase