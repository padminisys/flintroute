# System Architecture Overview

## Table of Contents
- [Introduction](#introduction)
- [High-Level Architecture](#high-level-architecture)
- [Component Architecture](#component-architecture)
- [Technology Stack](#technology-stack)
- [Design Principles](#design-principles)
- [Deployment Models](#deployment-models)
- [Scalability Considerations](#scalability-considerations)

---

## Introduction

FlintRoute is designed as a modern, cloud-native web application for managing FRR (Free Range Routing) instances. The architecture follows a three-tier model with clear separation of concerns:

1. **Presentation Layer**: React-based web UI
2. **Application Layer**: Go backend with business logic
3. **Data Layer**: FRR routing daemon with gRPC northbound API

### Design Goals

- **Security First**: Authentication, authorization, and audit logging built-in
- **Real-time**: WebSocket-based live updates for routing state
- **Atomic Operations**: All-or-nothing configuration changes
- **Scalability**: Support for managing multiple FRR instances
- **Reliability**: Configuration backup, rollback, and drift detection
- **Observability**: Built-in metrics, logging, and monitoring

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         End Users                                │
│                    (Network Engineers)                           │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS/WSS
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                      Load Balancer                               │
│                   (nginx/HAProxy)                                │
│                    TLS Termination                               │
└────────────────────────────┬────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
┌───────▼────────┐  ┌───────▼────────┐  ┌───────▼────────┐
│  FlintRoute    │  │  FlintRoute    │  │  FlintRoute    │
│   Instance 1   │  │   Instance 2   │  │   Instance N   │
│                │  │                │  │                │
│ ┌────────────┐ │  │ ┌────────────┐ │  │ ┌────────────┐ │
│ │  Frontend  │ │  │ │  Frontend  │ │  │ │  Frontend  │ │
│ │   (React)  │ │  │ │   (React)  │ │  │ │   (React)  │ │
│ └─────┬──────┘ │  │ └─────┬──────┘ │  │ └─────┬──────┘ │
│       │        │  │       │        │  │       │        │
│ ┌─────▼──────┐ │  │ ┌─────▼──────┐ │  │ ┌─────▼──────┐ │
│ │  Backend   │ │  │ │  Backend   │ │  │ │  Backend   │ │
│ │    (Go)    │ │  │ │    (Go)    │ │  │ │    (Go)    │ │
│ └─────┬──────┘ │  │ └─────┬──────┘ │  │ └─────┬──────┘ │
└───────┼────────┘  └───────┼────────┘  └───────┼────────┘
        │                   │                   │
        │ gRPC              │ gRPC              │ gRPC
        │                   │                   │
┌───────▼────────┐  ┌───────▼────────┐  ┌───────▼────────┐
│  FRR Instance  │  │  FRR Instance  │  │  FRR Instance  │
│   (Router 1)   │  │   (Router 2)   │  │   (Router N)   │
│                │  │                │  │                │
│ ┌────────────┐ │  │ ┌────────────┐ │  │ ┌────────────┐ │
│ │   BGPd     │ │  │ │   BGPd     │ │  │ │   BGPd     │ │
│ │   Zebra    │ │  │ │   Zebra    │ │  │ │   Zebra    │ │
│ │   Static   │ │  │ │   Static   │ │  │ │   Static   │ │
│ │   OSPF     │ │  │ │   OSPF     │ │  │ │   OSPF     │ │
│ └────────────┘ │  │ └────────────┘ │  │ └────────────┘ │
└────────────────┘  └────────────────┘  └────────────────┘
```

---

## Component Architecture

### Frontend (React Application)

```
┌─────────────────────────────────────────────────────────┐
│                    React Application                     │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │              UI Components                      │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐       │    │
│  │  │Dashboard │ │BGP Peers │ │Sessions  │       │    │
│  │  └──────────┘ └──────────┘ └──────────┘       │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐       │    │
│  │  │Config    │ │Alerts    │ │Topology  │       │    │
│  │  └──────────┘ └──────────┘ └──────────┘       │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │           State Management (Redux)              │    │
│  │  - Authentication State                         │    │
│  │  - BGP Session State                            │    │
│  │  - Configuration State                          │    │
│  │  - Alert State                                  │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │              API Client Layer                   │    │
│  │  - REST API Client                              │    │
│  │  - WebSocket Client                             │    │
│  │  - Authentication Interceptor                   │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

**Key Features:**
- **Component-Based**: Reusable UI components with Material-UI or Ant Design
- **State Management**: Redux for global state, React Query for server state
- **Real-time Updates**: WebSocket connection for live BGP session updates
- **Responsive Design**: Mobile-friendly interface
- **Type Safety**: TypeScript for compile-time type checking

### Backend (Go Application)

```
┌─────────────────────────────────────────────────────────┐
│                    Go Backend Service                    │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │              HTTP/WebSocket Server              │    │
│  │  - REST API Endpoints                           │    │
│  │  - WebSocket Handler                            │    │
│  │  - Static File Server                           │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │            Middleware Layer                     │    │
│  │  - Authentication (JWT)                         │    │
│  │  - Authorization (RBAC)                         │    │
│  │  - Logging & Audit                              │    │
│  │  - Rate Limiting                                │    │
│  │  - CORS                                         │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │            Business Logic Layer                 │    │
│  │  ┌──────────────┐  ┌──────────────┐            │    │
│  │  │ BGP Service  │  │Config Service│            │    │
│  │  └──────────────┘  └──────────────┘            │    │
│  │  ┌──────────────┐  ┌──────────────┐            │    │
│  │  │Monitor Svc   │  │ Alert Service│            │    │
│  │  └──────────────┘  └──────────────┘            │    │
│  │  ┌──────────────┐  ┌──────────────┐            │    │
│  │  │ Auth Service │  │Backup Service│            │    │
│  │  └──────────────┘  └──────────────┘            │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │            FRR Integration Layer                │    │
│  │  - gRPC Client Pool                             │    │
│  │  - Connection Management                        │    │
│  │  - Request/Response Translation                 │    │
│  │  - Error Handling & Retry                       │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │            Data Access Layer                    │    │
│  │  - Configuration Repository                     │    │
│  │  - User Repository                              │    │
│  │  - Audit Log Repository                         │    │
│  │  - Metrics Repository                           │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │              Storage Layer                      │    │
│  │  - SQLite (default)                             │    │
│  │  - PostgreSQL (optional)                        │    │
│  │  - File System (backups)                        │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

**Key Features:**
- **Modular Design**: Clear separation of concerns with service-oriented architecture
- **gRPC Client**: Efficient communication with FRR northbound API
- **Transaction Support**: Atomic configuration changes with rollback capability
- **Event-Driven**: WebSocket-based real-time updates to frontend
- **Observability**: Structured logging, metrics collection, and tracing

### FRR Integration

```
┌─────────────────────────────────────────────────────────┐
│                    FRR (Free Range Routing)              │
│                                                          │
│  ┌────────────────────────────────────────────────┐    │
│  │          gRPC Northbound API Server             │    │
│  │  - Configuration Management                     │    │
│  │  - State Retrieval                              │    │
│  │  - Operational Commands                         │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │              Management Layer                   │    │
│  │  - YANG Models                                  │    │
│  │  - Configuration Database                       │    │
│  │  - State Database                               │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │              Protocol Daemons                   │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐       │    │
│  │  │  BGPd    │ │  OSPFd   │ │  Staticd │       │    │
│  │  └──────────┘ └──────────┘ └──────────┘       │    │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐       │    │
│  │  │  IS-ISd  │ │  RIPd    │ │  BFDd    │       │    │
│  │  └──────────┘ └──────────┘ └──────────┘       │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │                  Zebra                          │    │
│  │  - Routing Table Management                     │    │
│  │  - Kernel Interface                             │    │
│  │  - Protocol Coordination                        │    │
│  └────────────────────────────────────────────────┘    │
│                         │                               │
│  ┌────────────────────────────────────────────────┐    │
│  │              Linux Kernel                       │    │
│  │  - Routing Tables (FIB)                         │    │
│  │  - Network Interfaces                           │    │
│  │  - Netfilter/iptables                           │    │
│  └────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────┘
```

---

## Technology Stack

### Frontend Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Framework | React | 18+ | UI framework |
| Language | TypeScript | 5+ | Type-safe JavaScript |
| State Management | Redux Toolkit | 2+ | Global state management |
| Server State | React Query | 5+ | Server state caching |
| UI Library | Material-UI / Ant Design | Latest | Component library |
| Routing | React Router | 6+ | Client-side routing |
| WebSocket | Socket.io-client | 4+ | Real-time communication |
| Charts | Recharts / D3.js | Latest | Data visualization |
| Forms | React Hook Form | 7+ | Form management |
| Build Tool | Vite | 5+ | Fast build tool |

### Backend Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.21+ | Backend language |
| HTTP Framework | Gin / Echo | Latest | Web framework |
| gRPC | grpc-go | 1.60+ | FRR communication |
| WebSocket | gorilla/websocket | Latest | Real-time updates |
| Database | SQLite / PostgreSQL | Latest | Data persistence |
| ORM | GORM | 1.25+ | Database abstraction |
| Authentication | golang-jwt | 5+ | JWT token handling |
| Logging | zap / logrus | Latest | Structured logging |
| Configuration | viper | Latest | Config management |
| Testing | testify | Latest | Testing framework |

### Infrastructure Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| OS | Debian 12 (Bookworm) | Base operating system |
| Routing | FRR stable | Routing daemon |
| Service Manager | systemd | Service orchestration |
| Reverse Proxy | nginx / Caddy | TLS termination, load balancing |
| Container | Docker (optional) | Development environment |
| Orchestration | Kubernetes (future) | Multi-node deployment |

---

## Design Principles

### 1. Security by Design

- **Authentication**: JWT tokens with secure storage
- **Authorization**: Role-Based Access Control (RBAC)
- **Encryption**: TLS 1.3 for all communications
- **Audit Logging**: Complete audit trail of all actions
- **Input Validation**: Comprehensive validation at all layers
- **Least Privilege**: Minimal permissions by default

### 2. Reliability

- **Atomic Transactions**: All-or-nothing configuration changes
- **Configuration Backup**: Automatic backup before changes
- **Rollback Support**: Quick rollback to previous configurations
- **Drift Detection**: Monitor for out-of-band changes
- **Health Checks**: Continuous health monitoring
- **Graceful Degradation**: Fail safely when components are unavailable

### 3. Performance

- **Efficient gRPC**: Binary protocol for FRR communication
- **Connection Pooling**: Reuse gRPC connections
- **Caching**: Cache frequently accessed data
- **Lazy Loading**: Load data on demand
- **WebSocket**: Efficient real-time updates
- **Database Indexing**: Optimize query performance

### 4. Scalability

- **Stateless Backend**: Horizontal scaling capability
- **Multi-Instance**: Manage multiple FRR instances
- **Load Balancing**: Distribute load across instances
- **Database Scaling**: Support for PostgreSQL clustering
- **Caching Layer**: Redis for distributed caching (future)

### 5. Observability

- **Structured Logging**: JSON-formatted logs
- **Metrics Collection**: Built-in metrics endpoint
- **Distributed Tracing**: Request tracing across components
- **Health Endpoints**: Liveness and readiness probes
- **Audit Trail**: Complete action history

### 6. Maintainability

- **Clean Architecture**: Clear separation of concerns
- **Dependency Injection**: Testable components
- **Documentation**: Comprehensive inline documentation
- **Testing**: Unit, integration, and e2e tests
- **Code Quality**: Linting and static analysis
- **Version Control**: Git-based workflow

---

## Deployment Models

### Single-Node Deployment

```
┌─────────────────────────────────────┐
│         Single Server               │
│                                     │
│  ┌───────────────────────────────┐ │
│  │     FlintRoute Service        │ │
│  │  (Frontend + Backend)         │ │
│  └───────────────┬───────────────┘ │
│                  │                  │
│  ┌───────────────▼───────────────┐ │
│  │         FRR Instance          │ │
│  │    (BGPd, Zebra, etc.)        │ │
│  └───────────────────────────────┘ │
│                                     │
│  ┌───────────────────────────────┐ │
│  │      SQLite Database          │ │
│  └───────────────────────────────┘ │
└─────────────────────────────────────┘
```

**Use Case**: Small deployments, development, testing
**Pros**: Simple setup, minimal resources
**Cons**: Single point of failure, limited scalability

### Multi-Node Deployment

```
┌──────────────────┐     ┌──────────────────┐
│  FlintRoute UI   │     │  FlintRoute UI   │
│    Instance 1    │     │    Instance 2    │
└────────┬─────────┘     └────────┬─────────┘
         │                        │
         └────────┬───────────────┘
                  │
         ┌────────▼─────────┐
         │  Load Balancer   │
         └────────┬─────────┘
                  │
    ┌─────────────┼─────────────┐
    │             │             │
┌───▼────┐   ┌───▼────┐   ┌───▼────┐
│ FRR 1  │   │ FRR 2  │   │ FRR N  │
│Router 1│   │Router 2│   │Router N│
└────────┘   └────────┘   └────────┘
```

**Use Case**: Production deployments, high availability
**Pros**: Redundancy, scalability, load distribution
**Cons**: More complex setup, requires coordination

### Centralized Management

```
┌─────────────────────────────────────────┐
│      FlintRoute Central Manager         │
│                                         │
│  ┌───────────────────────────────────┐ │
│  │     Management UI & Backend       │ │
│  └───────────────┬───────────────────┘ │
│                  │                      │
│  ┌───────────────▼───────────────────┐ │
│  │      PostgreSQL Database          │ │
│  └───────────────────────────────────┘ │
└──────────────────┬──────────────────────┘
                   │
        ┌──────────┼──────────┐
        │          │          │
    ┌───▼───┐  ┌──▼────┐  ┌──▼────┐
    │ FRR   │  │ FRR   │  │ FRR   │
    │Site 1 │  │Site 2 │  │Site N │
    └───────┘  └───────┘  └───────┘
```

**Use Case**: Multi-site deployments, enterprise
**Pros**: Centralized control, unified view
**Cons**: Network dependency, complexity

---

## Scalability Considerations

### Horizontal Scaling

**Frontend Scaling:**
- Static assets served via CDN
- Multiple UI instances behind load balancer
- Session affinity not required (stateless)

**Backend Scaling:**
- Stateless design enables horizontal scaling
- Load balancer distributes requests
- Shared database for state persistence
- WebSocket connections distributed across instances

**Database Scaling:**
- SQLite for single-node deployments
- PostgreSQL with replication for multi-node
- Read replicas for query scaling
- Connection pooling for efficiency

### Vertical Scaling

**Resource Requirements:**
- **Minimum**: 2 CPU cores, 4GB RAM
- **Recommended**: 4 CPU cores, 8GB RAM
- **Large Deployment**: 8+ CPU cores, 16GB+ RAM

**FRR Resource Usage:**
- Scales with number of routes and peers
- BGP full table: ~2GB RAM
- Multiple protocol daemons: additional overhead

### Performance Optimization

1. **Caching Strategy**
   - Cache BGP session state (30s TTL)
   - Cache configuration data (5m TTL)
   - Invalidate on configuration changes

2. **Database Optimization**
   - Index frequently queried fields
   - Partition large tables
   - Archive old audit logs

3. **Network Optimization**
   - gRPC connection pooling
   - HTTP/2 for frontend
   - WebSocket for real-time updates

4. **Monitoring**
   - Track response times
   - Monitor resource usage
   - Alert on performance degradation

---

## Next Steps

- [Security Architecture](security.md)
- [State Management](state-management.md)
- [Architecture Diagrams](diagrams.md)
- [API Documentation](../api/grpc-services.md)