# FlintRoute Documentation

Welcome to the FlintRoute documentation! This guide will help you understand, deploy, and contribute to FlintRoute - an open-source FRR management UI.

## 📚 Documentation Structure

### Getting Started

- **[README](../README.md)** - Project overview and quick start
- **[Installation Guide](deployment/installation.md)** - How to install FlintRoute
- **[Configuration Guide](deployment/configuration.md)** - Configuration options
- **[Quick Start Tutorial](getting-started.md)** - Step-by-step tutorial (coming soon)

### Architecture

- **[System Architecture](architecture/overview.md)** - High-level architecture and components
- **[Security Model](architecture/security.md)** - Authentication, authorization, and security
- **[State Management](architecture/state-management.md)** - How state is managed across the system
- **[Architecture Diagrams](architecture/diagrams.md)** - Visual architecture diagrams

### API Documentation

- **[gRPC Services](api/grpc-services.md)** - FRR gRPC integration
- **[Data Models](api/data-models.md)** - Data structures and schemas
- **[REST Endpoints](api/rest-endpoints.md)** - HTTP API reference

### Development

- **[Development Setup](development/setup.md)** - Set up your development environment
- **[FRR Installation](development/frr-installation.md)** - Install and configure FRR
- **[Testing Guide](development/testing.md)** - Testing with simulators
- **[Contributing Guidelines](../CONTRIBUTING.md)** - How to contribute

### Deployment

- **[Installation](deployment/installation.md)** - Production installation
- **[Configuration](deployment/configuration.md)** - Configuration management
- **[Systemd Service](deployment/systemd.md)** - Running as a system service
- **[Monitoring](deployment/monitoring.md)** - Monitoring and observability (coming soon)

### Phase 1 MVP

- **[Features](phase1/features.md)** - Phase 1 feature specifications
- **[Roadmap](phase1/roadmap.md)** - Implementation timeline
- **[Success Criteria](phase1/success-criteria.md)** - Acceptance criteria
- **[Testing Checklist](phase1/testing-checklist.md)** - Testing requirements

---

## 🎯 Quick Navigation

### For Users

**I want to...**
- [Install FlintRoute](deployment/installation.md)
- [Configure FlintRoute](deployment/configuration.md)
- [Manage BGP peers](phase1/features.md#bgp-peer-management)
- [Monitor sessions](phase1/features.md#session-monitoring)
- [Set up alerts](phase1/features.md#alerting-system)
- [Backup configurations](phase1/features.md#configuration-management)

### For Developers

**I want to...**
- [Set up development environment](development/setup.md)
- [Understand the architecture](architecture/overview.md)
- [Learn about security](architecture/security.md)
- [Contribute code](../CONTRIBUTING.md)
- [Run tests](development/testing.md)
- [Understand the API](api/grpc-services.md)

### For Operators

**I want to...**
- [Deploy to production](deployment/installation.md)
- [Configure systemd service](deployment/systemd.md)
- [Monitor the system](deployment/monitoring.md)
- [Troubleshoot issues](troubleshooting.md) (coming soon)
- [Upgrade FlintRoute](deployment/upgrade.md) (coming soon)

---

## 📖 Documentation by Topic

### BGP Management

- [BGP Peer Management](phase1/features.md#bgp-peer-management)
- [Session Monitoring](phase1/features.md#session-monitoring)
- [Route Policies](phase1/features.md#route-policies) (Phase 2)
- [Peer Groups](phase1/features.md#peer-groups) (Phase 2)

### Security

- [Authentication](architecture/security.md#authentication)
- [Authorization (RBAC)](architecture/security.md#authorization-rbac)
- [Encryption](architecture/security.md#encryption--transport-security)
- [Audit Logging](architecture/security.md#audit-logging)
- [Security Best Practices](architecture/security.md#security-best-practices)

### Configuration

- [Configuration Management](phase1/features.md#configuration-management)
- [Backup & Restore](phase1/features.md#backup--restore)
- [Configuration Drift](architecture/state-management.md#conflict-resolution)
- [Version Control](architecture/state-management.md#configuration-versioning)

### Monitoring & Alerting

- [Real-time Monitoring](phase1/features.md#session-monitoring)
- [Alert Configuration](phase1/features.md#alerting-system)
- [Metrics Collection](architecture/overview.md#observability)
- [Dashboard](phase1/features.md#visualization)

---

## 🔧 Technical Reference

### Architecture Components

| Component | Description | Documentation |
|-----------|-------------|---------------|
| Frontend | React web application | [Architecture](architecture/overview.md#frontend-react-application) |
| Backend | Go API server | [Architecture](architecture/overview.md#backend-go-application) |
| FRR Integration | gRPC client | [gRPC Services](api/grpc-services.md) |
| Database | SQLite/PostgreSQL | [Data Models](api/data-models.md) |
| Authentication | JWT + RBAC | [Security](architecture/security.md#authentication) |

### Technology Stack

| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| Frontend | React | 18+ | UI framework |
| Frontend | TypeScript | 5+ | Type safety |
| Frontend | Redux Toolkit | 2+ | State management |
| Backend | Go | 1.21+ | API server |
| Backend | gRPC | 1.60+ | FRR communication |
| Database | SQLite | 3+ | Data storage |
| Routing | FRR | 8.4+ | Routing daemon |
| OS | Debian | 12 | Base system |

---

## 📋 Common Tasks

### Installation & Setup

1. [Install FRR on Debian 12](development/frr-installation.md)
2. [Install FlintRoute](deployment/installation.md)
3. [Configure FlintRoute](deployment/configuration.md)
4. [Set up systemd service](deployment/systemd.md)
5. [Access the web UI](deployment/installation.md#accessing-the-ui)

### Development Workflow

1. [Set up development environment](development/setup.md)
2. [Run development servers](development/setup.md#running-development-servers)
3. [Make code changes](../CONTRIBUTING.md#development-workflow)
4. [Run tests](development/testing.md)
5. [Submit pull request](../CONTRIBUTING.md#pull-request-process)

### Operations

1. [Monitor system health](deployment/monitoring.md)
2. [Review audit logs](architecture/security.md#audit-logging)
3. [Backup configurations](phase1/features.md#configuration-management)
4. [Restore from backup](phase1/features.md#backup--restore)
5. [Troubleshoot issues](troubleshooting.md) (coming soon)

---

## 🎓 Learning Path

### Beginner

1. Read the [README](../README.md)
2. Understand the [Architecture Overview](architecture/overview.md)
3. Follow the [Installation Guide](deployment/installation.md)
4. Try the [Quick Start Tutorial](getting-started.md) (coming soon)

### Intermediate

1. Learn about [Security](architecture/security.md)
2. Understand [State Management](architecture/state-management.md)
3. Explore the [API Documentation](api/grpc-services.md)
4. Review [Phase 1 Features](phase1/features.md)

### Advanced

1. Study the [Architecture Diagrams](architecture/diagrams.md)
2. Review [Data Models](api/data-models.md)
3. Understand [Testing Strategy](development/testing.md)
4. Contribute to [Development](../CONTRIBUTING.md)

---

## 🔍 Search & Find

### By Feature

- **BGP**: [Features](phase1/features.md#bgp-peer-management), [API](api/grpc-services.md#bgp-service)
- **Authentication**: [Security](architecture/security.md#authentication), [API](api/rest-endpoints.md#authentication)
- **Monitoring**: [Features](phase1/features.md#session-monitoring), [Architecture](architecture/overview.md#observability)
- **Configuration**: [Features](phase1/features.md#configuration-management), [State Management](architecture/state-management.md#configuration-management)
- **Alerts**: [Features](phase1/features.md#alerting-system), [API](api/rest-endpoints.md#alerts)

### By Role

**Network Engineer**
- [BGP Management](phase1/features.md#bgp-peer-management)
- [Session Monitoring](phase1/features.md#session-monitoring)
- [Configuration Backup](phase1/features.md#configuration-management)

**System Administrator**
- [Installation](deployment/installation.md)
- [Configuration](deployment/configuration.md)
- [Systemd Service](deployment/systemd.md)
- [Monitoring](deployment/monitoring.md)

**Developer**
- [Development Setup](development/setup.md)
- [Architecture](architecture/overview.md)
- [API Documentation](api/grpc-services.md)
- [Contributing](../CONTRIBUTING.md)

**Security Engineer**
- [Security Architecture](architecture/security.md)
- [Authentication](architecture/security.md#authentication)
- [Authorization](architecture/security.md#authorization-rbac)
- [Audit Logging](architecture/security.md#audit-logging)

---

## 📝 Documentation Status

### Complete ✅

- [x] README
- [x] Architecture Overview
- [x] Security Documentation
- [x] State Management
- [x] Phase 1 Features
- [x] Phase 1 Roadmap
- [x] Contributing Guidelines

### In Progress 🚧

- [ ] Architecture Diagrams
- [ ] gRPC Services Documentation
- [ ] Data Models Documentation
- [ ] Development Setup Guide
- [ ] FRR Installation Guide
- [ ] Testing Guide

### Planned 📋

- [ ] Quick Start Tutorial
- [ ] User Guide
- [ ] Administrator Guide
- [ ] Troubleshooting Guide
- [ ] FAQ
- [ ] API Reference (OpenAPI)
- [ ] Deployment Best Practices
- [ ] Performance Tuning Guide
- [ ] Upgrade Guide

---

## 🤝 Contributing to Documentation

Documentation improvements are always welcome! See [Contributing Guidelines](../CONTRIBUTING.md#documentation) for details.

**Ways to contribute:**
- Fix typos and grammar
- Improve clarity and examples
- Add missing documentation
- Update outdated information
- Translate documentation
- Add diagrams and screenshots

---

## 📞 Getting Help

- **GitHub Issues**: [Report bugs or request features](https://github.com/padminisys/flintroute/issues)
- **GitHub Discussions**: [Ask questions and discuss ideas](https://github.com/padminisys/flintroute/discussions)
- **Email**: support@flintroute.com
- **Website**: [flintroute.com](https://flintroute.com)

---

## 📄 License

FlintRoute is licensed under the Apache License 2.0. See [LICENSE](../LICENSE) for details.

---

**Last Updated**: 2024-01-15
**Version**: 0.1.0-alpha
**Status**: Phase 1 Development