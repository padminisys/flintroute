# Security Architecture

## Table of Contents
- [Overview](#overview)
- [Authentication](#authentication)
- [Authorization (RBAC)](#authorization-rbac)
- [Encryption & Transport Security](#encryption--transport-security)
- [Audit Logging](#audit-logging)
- [Input Validation](#input-validation)
- [Session Management](#session-management)
- [API Security](#api-security)
- [Configuration Security](#configuration-security)
- [Threat Model](#threat-model)
- [Security Best Practices](#security-best-practices)

---

## Overview

FlintRoute manages critical network infrastructure with root-level access to routing configurations. Security is paramount and built into every layer of the application.

### Security Principles

1. **Defense in Depth**: Multiple layers of security controls
2. **Least Privilege**: Minimal permissions by default
3. **Zero Trust**: Verify every request, trust nothing
4. **Audit Everything**: Complete audit trail of all actions
5. **Fail Secure**: Default to deny on errors
6. **Secure by Default**: Security enabled out of the box

### Security Layers

```
┌─────────────────────────────────────────────────────────┐
│                    Network Layer                         │
│  - TLS 1.3 Encryption                                   │
│  - Certificate Validation                               │
│  - Firewall Rules                                       │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                 Application Layer                        │
│  - JWT Authentication                                   │
│  - RBAC Authorization                                   │
│  - Rate Limiting                                        │
│  - Input Validation                                     │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   Data Layer                            │
│  - Encrypted Storage                                    │
│  - Audit Logging                                        │
│  - Configuration Backup                                 │
└─────────────────────────────────────────────────────────┘
```

---

## Authentication

### JWT-Based Authentication

FlintRoute uses JSON Web Tokens (JWT) for stateless authentication.

#### Token Structure

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "user_id",
    "username": "admin",
    "email": "admin@example.com",
    "roles": ["admin", "operator"],
    "permissions": ["bgp:read", "bgp:write", "config:write"],
    "iat": 1234567890,
    "exp": 1234571490,
    "jti": "unique-token-id"
  },
  "signature": "..."
}
```

#### Token Lifecycle

```
┌──────────────┐
│   User Login │
└──────┬───────┘
       │
       ▼
┌──────────────────────────────┐
│  Validate Credentials        │
│  - Username/Password         │
│  - LDAP/AD (optional)        │
│  - MFA (optional)            │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Generate JWT Token          │
│  - Access Token (15 min)     │
│  - Refresh Token (7 days)    │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Return Tokens to Client     │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Client Stores Tokens        │
│  - Access: Memory            │
│  - Refresh: HttpOnly Cookie  │
└──────────────────────────────┘
```

#### Token Validation

Every API request validates the JWT token:

1. **Signature Verification**: Verify token signature with public key
2. **Expiration Check**: Ensure token hasn't expired
3. **Revocation Check**: Check against revoked token list
4. **Claims Validation**: Verify required claims are present
5. **Permission Check**: Validate user has required permissions

#### Token Refresh Flow

```
┌──────────────────────────────┐
│  Access Token Expired        │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Client Sends Refresh Token  │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Validate Refresh Token      │
│  - Not expired               │
│  - Not revoked               │
│  - Valid signature           │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Generate New Access Token   │
│  (Optionally new refresh)    │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Return New Tokens           │
└──────────────────────────────┘
```

### Multi-Factor Authentication (MFA)

**Phase 2 Feature** - Support for:
- TOTP (Time-based One-Time Password)
- SMS-based codes
- Hardware tokens (YubiKey)
- Backup codes

### External Authentication

**Phase 2 Feature** - Integration with:
- LDAP/Active Directory
- SAML 2.0
- OAuth 2.0 / OpenID Connect
- RADIUS

---

## Authorization (RBAC)

### Role-Based Access Control

FlintRoute implements a flexible RBAC system with predefined roles and custom permissions.

#### Built-in Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| **Admin** | Full system access | All permissions |
| **Operator** | Day-to-day operations | Read all, write BGP/config |
| **Monitor** | Read-only access | Read all, no write |
| **Auditor** | Audit log access | Read audit logs only |

#### Permission Model

Permissions follow the format: `resource:action`

**Resource Types:**
- `bgp` - BGP configuration and peers
- `config` - System configuration
- `user` - User management
- `audit` - Audit logs
- `alert` - Alert configuration
- `backup` - Backup/restore operations

**Actions:**
- `read` - View resource
- `write` - Create/update resource
- `delete` - Delete resource
- `execute` - Execute operations

**Example Permissions:**
```
bgp:read          # View BGP configuration
bgp:write         # Modify BGP configuration
bgp:delete        # Delete BGP peers
config:read       # View system configuration
config:write      # Modify system configuration
user:write        # Create/modify users
audit:read        # View audit logs
backup:execute    # Execute backup/restore
```

#### Permission Hierarchy

```
┌─────────────────────────────────────────────────────────┐
│                      Admin Role                          │
│                    (All Permissions)                     │
└─────────────────────────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
┌───────▼────────┐ ┌─────▼──────┐ ┌───────▼────────┐
│   Operator     │ │  Monitor   │ │    Auditor     │
│                │ │            │ │                │
│ bgp:*          │ │ bgp:read   │ │ audit:read     │
│ config:read    │ │ config:read│ │                │
│ config:write   │ │ alert:read │ │                │
│ alert:*        │ │            │ │                │
│ backup:execute │ │            │ │                │
└────────────────┘ └────────────┘ └────────────────┘
```

#### Custom Roles

Administrators can create custom roles with specific permission sets:

```json
{
  "name": "bgp-operator",
  "description": "BGP-only operator",
  "permissions": [
    "bgp:read",
    "bgp:write",
    "bgp:delete",
    "config:read"
  ]
}
```

#### Authorization Flow

```
┌──────────────────────────────┐
│  User Makes Request          │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Extract JWT Token           │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Validate Token              │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Extract User Roles          │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Resolve Permissions         │
│  (from roles)                │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Check Required Permission   │
│  (resource:action)           │
└──────┬───────────────────────┘
       │
   ┌───┴───┐
   │       │
   ▼       ▼
┌─────┐ ┌──────┐
│Allow│ │ Deny │
└─────┘ └──────┘
```

---

## Encryption & Transport Security

### TLS Configuration

**Minimum TLS Version**: TLS 1.3
**Fallback**: TLS 1.2 (with strong ciphers only)

#### Recommended Cipher Suites (TLS 1.3)

```
TLS_AES_256_GCM_SHA384
TLS_CHACHA20_POLY1305_SHA256
TLS_AES_128_GCM_SHA256
```

#### Certificate Management

```
┌─────────────────────────────────────────────────────────┐
│              Certificate Authority (CA)                  │
│  - Internal CA or Let's Encrypt                         │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
┌───────▼────────┐ ┌──▼──────────┐ ┌▼────────────┐
│  Server Cert   │ │ Client Cert │ │  gRPC Cert  │
│  (Web UI)      │ │ (Optional)  │ │  (FRR)      │
└────────────────┘ └─────────────┘ └─────────────┘
```

**Certificate Requirements:**
- 2048-bit RSA or 256-bit ECDSA
- SHA-256 signature algorithm
- Valid for 90 days (auto-renewal)
- Subject Alternative Names (SAN) for all domains

### Data Encryption

#### At Rest

- **Database**: SQLite encryption extension or PostgreSQL pgcrypto
- **Backups**: AES-256 encryption
- **Secrets**: Encrypted with master key
- **Logs**: Sensitive data redacted

#### In Transit

- **Web UI ↔ Backend**: TLS 1.3
- **Backend ↔ FRR**: gRPC with TLS
- **WebSocket**: WSS (WebSocket Secure)

### Secret Management

```
┌─────────────────────────────────────────────────────────┐
│                   Master Key                             │
│  - Generated on first boot                              │
│  - Stored in secure location                            │
│  - Used to encrypt other secrets                        │
└─────────────────────┬───────────────────────────────────┘
                      │
        ┌─────────────┼─────────────┐
        │             │             │
┌───────▼────────┐ ┌──▼──────────┐ ┌▼────────────┐
│  JWT Secret    │ │ DB Password │ │  API Keys   │
│  (Encrypted)   │ │ (Encrypted) │ │ (Encrypted) │
└────────────────┘ └─────────────┘ └─────────────┘
```

**Secret Storage Options:**
1. **File-based**: Encrypted files on disk (default)
2. **Environment Variables**: For container deployments
3. **HashiCorp Vault**: For enterprise deployments (Phase 2)
4. **AWS Secrets Manager**: For cloud deployments (Phase 2)

---

## Audit Logging

### Audit Log Structure

Every action is logged with complete context:

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "event_id": "evt_abc123",
  "event_type": "bgp.peer.create",
  "severity": "info",
  "user": {
    "id": "user_123",
    "username": "admin",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0..."
  },
  "resource": {
    "type": "bgp_peer",
    "id": "peer_456",
    "name": "AS64512"
  },
  "action": {
    "operation": "create",
    "status": "success",
    "duration_ms": 234
  },
  "changes": {
    "before": null,
    "after": {
      "remote_as": 64512,
      "neighbor": "192.0.2.1",
      "description": "Upstream Provider"
    }
  },
  "metadata": {
    "request_id": "req_xyz789",
    "session_id": "sess_abc456"
  }
}
```

### Audit Event Types

| Category | Events |
|----------|--------|
| **Authentication** | login, logout, login_failed, token_refresh, password_change |
| **Authorization** | permission_denied, role_assigned, role_removed |
| **BGP** | peer_create, peer_update, peer_delete, peer_enable, peer_disable |
| **Configuration** | config_update, config_backup, config_restore, config_rollback |
| **User Management** | user_create, user_update, user_delete, user_disable |
| **System** | service_start, service_stop, service_restart, health_check |

### Audit Log Retention

- **Default Retention**: 90 days
- **Compliance Mode**: 1 year or more
- **Archive**: Compressed and encrypted
- **Immutable**: Cannot be modified or deleted by users

### Audit Log Access

- **View**: Auditor role or higher
- **Export**: Admin role only
- **Search**: Full-text search with filters
- **Alerts**: Configurable alerts on suspicious activity

---

## Input Validation

### Validation Layers

```
┌─────────────────────────────────────────────────────────┐
│                  Frontend Validation                     │
│  - Type checking (TypeScript)                           │
│  - Format validation (regex)                            │
│  - Range validation                                     │
│  - User feedback                                        │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Backend Validation                      │
│  - Schema validation (JSON Schema)                      │
│  - Business logic validation                            │
│  - SQL injection prevention                             │
│  - Command injection prevention                         │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   FRR Validation                         │
│  - Configuration syntax validation                      │
│  - Semantic validation                                  │
│  - Dry-run before apply                                 │
└─────────────────────────────────────────────────────────┘
```

### Validation Rules

#### IP Address Validation

```go
// IPv4 validation
func ValidateIPv4(ip string) error {
    if net.ParseIP(ip) == nil {
        return errors.New("invalid IPv4 address")
    }
    return nil
}

// IPv6 validation
func ValidateIPv6(ip string) error {
    if net.ParseIP(ip) == nil {
        return errors.New("invalid IPv6 address")
    }
    return nil
}

// CIDR validation
func ValidateCIDR(cidr string) error {
    _, _, err := net.ParseCIDR(cidr)
    return err
}
```

#### AS Number Validation

```go
func ValidateASN(asn uint32) error {
    // Valid range: 1-4294967295
    // Reserved: 0, 23456, 64496-64511, 65535, 4200000000-4294967295
    if asn == 0 || asn == 23456 || asn == 65535 {
        return errors.New("reserved AS number")
    }
    if asn >= 64496 && asn <= 64511 {
        return errors.New("reserved for documentation")
    }
    if asn >= 4200000000 {
        return errors.New("reserved for private use")
    }
    return nil
}
```

#### Configuration Validation

- **Syntax**: Validate against FRR configuration schema
- **Semantics**: Check for logical errors (e.g., duplicate peers)
- **Dependencies**: Verify required resources exist
- **Conflicts**: Detect conflicting configurations

### Sanitization

- **HTML**: Strip all HTML tags from user input
- **SQL**: Use parameterized queries (no string concatenation)
- **Shell**: Avoid shell execution, use native APIs
- **Path**: Validate and sanitize file paths

---

## Session Management

### Session Security

```
┌─────────────────────────────────────────────────────────┐
│                    Session Properties                    │
│                                                          │
│  - Session ID: Cryptographically random (256-bit)       │
│  - Timeout: 15 minutes of inactivity                    │
│  - Max Duration: 8 hours                                │
│  - Concurrent Sessions: 3 per user                      │
│  - IP Binding: Optional (detect IP changes)             │
└─────────────────────────────────────────────────────────┘
```

### Session Storage

- **Access Token**: Client-side (memory only, not localStorage)
- **Refresh Token**: HttpOnly cookie (secure, SameSite=Strict)
- **Session Data**: Server-side (Redis or database)

### Session Termination

- **Explicit Logout**: User-initiated
- **Timeout**: After inactivity period
- **Token Expiration**: After max duration
- **Admin Revocation**: Force logout by admin
- **Security Event**: Suspicious activity detected

---

## API Security

### Rate Limiting

```
┌─────────────────────────────────────────────────────────┐
│                    Rate Limit Tiers                      │
│                                                          │
│  Authentication Endpoints:                              │
│    - 5 requests per minute per IP                       │
│    - 20 requests per hour per IP                        │
│                                                          │
│  Read Operations:                                       │
│    - 100 requests per minute per user                   │
│    - 1000 requests per hour per user                    │
│                                                          │
│  Write Operations:                                      │
│    - 30 requests per minute per user                    │
│    - 300 requests per hour per user                     │
│                                                          │
│  Admin Operations:                                      │
│    - 10 requests per minute per user                    │
│    - 100 requests per hour per user                     │
└─────────────────────────────────────────────────────────┘
```

### CORS Configuration

```go
cors.Config{
    AllowOrigins:     []string{"https://flintroute.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}
```

### API Versioning

- **URL-based**: `/api/v1/bgp/peers`
- **Header-based**: `Accept: application/vnd.flintroute.v1+json`
- **Deprecation**: 6-month notice before removal

---

## Configuration Security

### Atomic Transactions

All configuration changes are atomic:

```
┌──────────────────────────────┐
│  Start Transaction           │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Backup Current Config       │
└──────┬───────────────────────┘
       │
       ▼
┌──────────────────────────────┐
│  Validate New Config         │
└──────┬───────────────────────┘
       │
   ┌───┴───┐
   │       │
   ▼       ▼
┌─────┐ ┌──────────────────────┐
│Fail │ │  Apply Config        │
│     │ └──────┬───────────────┘
│     │        │
│     │        ▼
│     │ ┌──────────────────────┐
│     │ │  Verify Success      │
│     │ └──────┬───────────────┘
│     │        │
│     │    ┌───┴───┐
│     │    │       │
│     │    ▼       ▼
│     │ ┌─────┐ ┌──────┐
│     │ │ OK  │ │ Fail │
│     │ └─────┘ └──┬───┘
│     │            │
│     ◄────────────┘
│     │
│     ▼
│ ┌──────────────────────┐
│ │  Rollback to Backup  │
│ └──────────────────────┘
│
▼
┌──────────────────────────────┐
│  Commit or Rollback          │
└──────────────────────────────┘
```

### Configuration Drift Detection

Monitor for out-of-band changes:

1. **Periodic Polling**: Check FRR config every 60 seconds
2. **Hash Comparison**: Compare config hash with expected
3. **Alert on Drift**: Notify admins of unexpected changes
4. **Sync Options**: 
   - Revert to FlintRoute config
   - Accept FRR config
   - Manual resolution

---

## Threat Model

### Identified Threats

| Threat | Mitigation |
|--------|-----------|
| **Unauthorized Access** | Strong authentication, MFA, rate limiting |
| **Privilege Escalation** | RBAC, least privilege, audit logging |
| **Man-in-the-Middle** | TLS 1.3, certificate pinning |
| **SQL Injection** | Parameterized queries, ORM |
| **XSS** | Content Security Policy, input sanitization |
| **CSRF** | SameSite cookies, CSRF tokens |
| **DoS** | Rate limiting, resource limits |
| **Configuration Tampering** | Atomic transactions, audit logging |
| **Credential Theft** | Secure storage, token rotation |
| **Session Hijacking** | Secure cookies, IP binding |

### Attack Scenarios

#### Scenario 1: Brute Force Login

**Attack**: Attacker attempts to guess passwords
**Mitigation**:
- Rate limiting (5 attempts per minute)
- Account lockout after 10 failed attempts
- CAPTCHA after 3 failed attempts
- Alert admins on suspicious activity

#### Scenario 2: Privilege Escalation

**Attack**: User attempts to access unauthorized resources
**Mitigation**:
- Strict RBAC enforcement
- Permission checks on every request
- Audit logging of all access attempts
- Alert on repeated permission denials

#### Scenario 3: Configuration Injection

**Attack**: Malicious configuration to disrupt routing
**Mitigation**:
- Input validation at all layers
- Configuration syntax validation
- Dry-run before apply
- Automatic rollback on errors
- Audit trail of all changes

---

## Security Best Practices

### Deployment Security

1. **Network Segmentation**
   - Isolate FlintRoute in management network
   - Firewall rules to restrict access
   - VPN for remote access

2. **System Hardening**
   - Minimal OS installation
   - Disable unnecessary services
   - Regular security updates
   - SELinux or AppArmor enabled

3. **Access Control**
   - SSH key-based authentication only
   - Disable root login
   - Use sudo for privileged operations
   - Regular access review

### Operational Security

1. **Regular Audits**
   - Review audit logs weekly
   - Check for suspicious activity
   - Verify user permissions
   - Review configuration changes

2. **Backup & Recovery**
   - Daily configuration backups
   - Encrypted backup storage
   - Test restore procedures
   - Off-site backup copies

3. **Monitoring**
   - Real-time security alerts
   - Failed login monitoring
   - Configuration change alerts
   - Resource usage monitoring

### Development Security

1. **Secure Coding**
   - Follow OWASP guidelines
   - Code review for all changes
   - Static analysis tools
   - Dependency scanning

2. **Testing**
   - Security testing in CI/CD
   - Penetration testing
   - Vulnerability scanning
   - Fuzzing critical components

3. **Secrets Management**
   - Never commit secrets to git
   - Use environment variables
   - Rotate secrets regularly
   - Use secret management tools

---

## Security Checklist

### Pre-Deployment

- [ ] TLS certificates configured
- [ ] Strong passwords set
- [ ] Default credentials changed
- [ ] Firewall rules configured
- [ ] Audit logging enabled
- [ ] Backup system configured
- [ ] Security updates applied

### Post-Deployment

- [ ] Regular security audits scheduled
- [ ] Monitoring alerts configured
- [ ] Incident response plan documented
- [ ] User training completed
- [ ] Backup restoration tested
- [ ] Access review process established

### Ongoing

- [ ] Weekly audit log review
- [ ] Monthly access review
- [ ] Quarterly security assessment
- [ ] Annual penetration testing
- [ ] Continuous security updates

---

## Compliance

### Standards & Frameworks

- **NIST Cybersecurity Framework**: Identify, Protect, Detect, Respond, Recover
- **CIS Controls**: Critical security controls implementation
- **ISO 27001**: Information security management
- **SOC 2**: Security, availability, and confidentiality

### Regulatory Compliance

- **GDPR**: Data protection and privacy (if applicable)
- **PCI DSS**: Payment card industry standards (if applicable)
- **HIPAA**: Healthcare data protection (if applicable)

---

## Next Steps

- [State Management](state-management.md)
- [Architecture Diagrams](diagrams.md)
- [API Security](../api/grpc-services.md)
- [Deployment Security](../deployment/installation.md)