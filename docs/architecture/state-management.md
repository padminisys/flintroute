# State Management Architecture

## Table of Contents
- [Overview](#overview)
- [State Categories](#state-categories)
- [Frontend State Management](#frontend-state-management)
- [Backend State Management](#backend-state-management)
- [FRR State Synchronization](#frr-state-synchronization)
- [Configuration Management](#configuration-management)
- [Real-time Updates](#real-time-updates)
- [Caching Strategy](#caching-strategy)
- [Conflict Resolution](#conflict-resolution)
- [State Persistence](#state-persistence)

---

## Overview

FlintRoute manages multiple types of state across different layers of the application. Proper state management is critical for:

- **Consistency**: Ensuring UI reflects actual FRR state
- **Performance**: Minimizing unnecessary API calls
- **Reliability**: Handling network failures gracefully
- **User Experience**: Providing real-time updates

### State Flow Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend State                        │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │ UI State   │  │ App State  │  │Server State│        │
│  │ (Local)    │  │ (Redux)    │  │(React Query)│       │
│  └────────────┘  └────────────┘  └──────┬─────┘        │
└────────────────────────────────────────┬─┴──────────────┘
                                         │
                                    HTTP/WebSocket
                                         │
┌────────────────────────────────────────▼────────────────┐
│                    Backend State                         │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │ Memory     │  │ Cache      │  │ Database   │        │
│  │ (Runtime)  │  │ (Redis)    │  │ (SQLite)   │        │
│  └────────────┘  └────────────┘  └────────────┘        │
└────────────────────────────────────────┬────────────────┘
                                         │
                                      gRPC
                                         │
┌────────────────────────────────────────▼────────────────┐
│                      FRR State                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐        │
│  │ Config DB  │  │ State DB   │  │ Runtime    │        │
│  │ (YANG)     │  │ (YANG)     │  │ (Daemons)  │        │
│  └────────────┘  └────────────┘  └────────────┘        │
└─────────────────────────────────────────────────────────┘
```

---

## State Categories

### 1. UI State (Frontend Only)

**Ephemeral state that doesn't need to persist:**

- Form input values
- Modal open/closed state
- Sidebar expanded/collapsed
- Selected tabs
- Sorting and filtering preferences
- Scroll positions

**Management**: React component state (useState, useReducer)

### 2. Application State (Frontend)

**Shared state across components:**

- User authentication status
- Current user information
- Active navigation
- Theme preferences
- Notification queue
- Global loading states

**Management**: Redux Toolkit

### 3. Server State (Frontend)

**Data fetched from backend:**

- BGP peer list
- Session status
- Configuration data
- Audit logs
- Alert history
- Metrics data

**Management**: React Query (TanStack Query)

### 4. Backend Application State

**Runtime state in backend:**

- Active WebSocket connections
- gRPC connection pool
- In-progress transactions
- Rate limit counters
- Cache entries

**Management**: In-memory data structures

### 5. Persistent State (Backend)

**Data that must survive restarts:**

- User accounts and roles
- Configuration history
- Audit logs
- Alert configurations
- Backup metadata

**Management**: SQLite/PostgreSQL database

### 6. FRR State

**Routing daemon state:**

- Running configuration
- BGP session state
- Routing tables
- Interface status
- Protocol statistics

**Management**: FRR internal databases (accessed via gRPC)

---

## Frontend State Management

### Redux Store Structure

```typescript
interface RootState {
  auth: {
    isAuthenticated: boolean;
    user: User | null;
    token: string | null;
    refreshToken: string | null;
    permissions: string[];
  };
  
  ui: {
    theme: 'light' | 'dark';
    sidebarOpen: boolean;
    notifications: Notification[];
    loading: {
      [key: string]: boolean;
    };
  };
  
  router: {
    currentRoute: string;
    previousRoute: string;
  };
  
  websocket: {
    connected: boolean;
    reconnecting: boolean;
    lastMessage: any;
  };
}
```

### Redux Slices

#### Auth Slice

```typescript
const authSlice = createSlice({
  name: 'auth',
  initialState: {
    isAuthenticated: false,
    user: null,
    token: null,
    refreshToken: null,
    permissions: [],
  },
  reducers: {
    loginSuccess: (state, action) => {
      state.isAuthenticated = true;
      state.user = action.payload.user;
      state.token = action.payload.token;
      state.refreshToken = action.payload.refreshToken;
      state.permissions = action.payload.permissions;
    },
    logout: (state) => {
      state.isAuthenticated = false;
      state.user = null;
      state.token = null;
      state.refreshToken = null;
      state.permissions = [];
    },
    updateToken: (state, action) => {
      state.token = action.payload.token;
    },
  },
});
```

#### UI Slice

```typescript
const uiSlice = createSlice({
  name: 'ui',
  initialState: {
    theme: 'light',
    sidebarOpen: true,
    notifications: [],
    loading: {},
  },
  reducers: {
    toggleSidebar: (state) => {
      state.sidebarOpen = !state.sidebarOpen;
    },
    setTheme: (state, action) => {
      state.theme = action.payload;
    },
    addNotification: (state, action) => {
      state.notifications.push(action.payload);
    },
    removeNotification: (state, action) => {
      state.notifications = state.notifications.filter(
        n => n.id !== action.payload
      );
    },
    setLoading: (state, action) => {
      state.loading[action.payload.key] = action.payload.value;
    },
  },
});
```

### React Query Configuration

```typescript
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // Stale time: how long data is considered fresh
      staleTime: 30 * 1000, // 30 seconds
      
      // Cache time: how long inactive data stays in cache
      cacheTime: 5 * 60 * 1000, // 5 minutes
      
      // Retry failed requests
      retry: 3,
      retryDelay: attemptIndex => Math.min(1000 * 2 ** attemptIndex, 30000),
      
      // Refetch on window focus
      refetchOnWindowFocus: true,
      
      // Refetch on reconnect
      refetchOnReconnect: true,
    },
    mutations: {
      retry: 1,
    },
  },
});
```

### Query Hooks

```typescript
// Fetch BGP peers
function useBGPPeers() {
  return useQuery({
    queryKey: ['bgp', 'peers'],
    queryFn: async () => {
      const response = await api.get('/api/v1/bgp/peers');
      return response.data;
    },
    staleTime: 30 * 1000, // 30 seconds
  });
}

// Fetch single peer with real-time updates
function useBGPPeer(peerId: string) {
  const queryClient = useQueryClient();
  
  // Subscribe to WebSocket updates
  useEffect(() => {
    const unsubscribe = websocket.subscribe(
      `bgp.peer.${peerId}`,
      (data) => {
        queryClient.setQueryData(['bgp', 'peer', peerId], data);
      }
    );
    return unsubscribe;
  }, [peerId, queryClient]);
  
  return useQuery({
    queryKey: ['bgp', 'peer', peerId],
    queryFn: async () => {
      const response = await api.get(`/api/v1/bgp/peers/${peerId}`);
      return response.data;
    },
  });
}

// Mutation for creating peer
function useCreateBGPPeer() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (peer: BGPPeerInput) => {
      const response = await api.post('/api/v1/bgp/peers', peer);
      return response.data;
    },
    onSuccess: () => {
      // Invalidate and refetch peer list
      queryClient.invalidateQueries({ queryKey: ['bgp', 'peers'] });
    },
  });
}
```

---

## Backend State Management

### In-Memory State

```go
type StateManager struct {
    // WebSocket connections
    wsConnections map[string]*websocket.Conn
    wsLock        sync.RWMutex
    
    // gRPC connection pool
    grpcPool      *GRPCPool
    
    // Active transactions
    transactions  map[string]*Transaction
    txLock        sync.RWMutex
    
    // Rate limiters
    rateLimiters  map[string]*rate.Limiter
    rlLock        sync.RWMutex
    
    // Cache
    cache         *Cache
}

// WebSocket connection management
func (sm *StateManager) AddWSConnection(userID string, conn *websocket.Conn) {
    sm.wsLock.Lock()
    defer sm.wsLock.Unlock()
    sm.wsConnections[userID] = conn
}

func (sm *StateManager) RemoveWSConnection(userID string) {
    sm.wsLock.Lock()
    defer sm.wsLock.Unlock()
    delete(sm.wsConnections, userID)
}

func (sm *StateManager) BroadcastToUser(userID string, message interface{}) error {
    sm.wsLock.RLock()
    conn, exists := sm.wsConnections[userID]
    sm.wsLock.RUnlock()
    
    if !exists {
        return errors.New("connection not found")
    }
    
    return conn.WriteJSON(message)
}
```

### Transaction Management

```go
type Transaction struct {
    ID            string
    UserID        string
    Type          string
    Status        TransactionStatus
    StartTime     time.Time
    EndTime       *time.Time
    BackupID      string
    Changes       []Change
    RollbackFunc  func() error
}

type TransactionStatus string

const (
    TxPending   TransactionStatus = "pending"
    TxRunning   TransactionStatus = "running"
    TxCommitted TransactionStatus = "committed"
    TxRolledBack TransactionStatus = "rolled_back"
    TxFailed    TransactionStatus = "failed"
)

func (sm *StateManager) BeginTransaction(userID, txType string) (*Transaction, error) {
    tx := &Transaction{
        ID:        generateID(),
        UserID:    userID,
        Type:      txType,
        Status:    TxPending,
        StartTime: time.Now(),
        Changes:   make([]Change, 0),
    }
    
    sm.txLock.Lock()
    sm.transactions[tx.ID] = tx
    sm.txLock.Unlock()
    
    return tx, nil
}

func (sm *StateManager) CommitTransaction(txID string) error {
    sm.txLock.Lock()
    tx, exists := sm.transactions[txID]
    sm.txLock.Unlock()
    
    if !exists {
        return errors.New("transaction not found")
    }
    
    tx.Status = TxCommitted
    now := time.Now()
    tx.EndTime = &now
    
    return nil
}

func (sm *StateManager) RollbackTransaction(txID string) error {
    sm.txLock.Lock()
    tx, exists := sm.transactions[txID]
    sm.txLock.Unlock()
    
    if !exists {
        return errors.New("transaction not found")
    }
    
    if tx.RollbackFunc != nil {
        if err := tx.RollbackFunc(); err != nil {
            return err
        }
    }
    
    tx.Status = TxRolledBack
    now := time.Now()
    tx.EndTime = &now
    
    return nil
}
```

---

## FRR State Synchronization

### Polling Strategy

```go
type FRRStateSyncer struct {
    grpcClient    *GRPCClient
    pollInterval  time.Duration
    lastSync      time.Time
    lastConfigHash string
}

func (s *FRRStateSyncer) Start(ctx context.Context) {
    ticker := time.NewTicker(s.pollInterval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := s.syncState(); err != nil {
                log.Error("Failed to sync FRR state", "error", err)
            }
        }
    }
}

func (s *FRRStateSyncer) syncState() error {
    // Fetch current configuration
    config, err := s.grpcClient.GetConfiguration()
    if err != nil {
        return err
    }
    
    // Calculate hash
    currentHash := calculateHash(config)
    
    // Check for drift
    if currentHash != s.lastConfigHash {
        log.Warn("Configuration drift detected")
        s.handleConfigDrift(config)
    }
    
    s.lastConfigHash = currentHash
    s.lastSync = time.Now()
    
    // Fetch BGP session state
    sessions, err := s.grpcClient.GetBGPSessions()
    if err != nil {
        return err
    }
    
    // Update cache and broadcast changes
    s.updateSessionState(sessions)
    
    return nil
}
```

### Event-Driven Updates

```go
type FRREventListener struct {
    grpcClient *GRPCClient
    eventChan  chan *FRREvent
}

func (l *FRREventListener) Start(ctx context.Context) {
    stream, err := l.grpcClient.SubscribeToEvents(ctx)
    if err != nil {
        log.Error("Failed to subscribe to FRR events", "error", err)
        return
    }
    
    for {
        event, err := stream.Recv()
        if err != nil {
            log.Error("Error receiving event", "error", err)
            return
        }
        
        l.eventChan <- event
    }
}

func (l *FRREventListener) HandleEvents(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case event := <-l.eventChan:
            l.processEvent(event)
        }
    }
}

func (l *FRREventListener) processEvent(event *FRREvent) {
    switch event.Type {
    case "bgp.peer.state_change":
        l.handleBGPStateChange(event)
    case "bgp.peer.established":
        l.handleBGPEstablished(event)
    case "bgp.peer.down":
        l.handleBGPDown(event)
    case "config.changed":
        l.handleConfigChange(event)
    }
}
```

---

## Configuration Management

### Configuration Versioning

```go
type ConfigVersion struct {
    ID          string
    Version     int
    Timestamp   time.Time
    UserID      string
    Description string
    ConfigData  []byte
    Hash        string
    ParentID    *string
}

type ConfigManager struct {
    db          *Database
    currentHash string
}

func (cm *ConfigManager) SaveConfiguration(userID, description string, config []byte) (*ConfigVersion, error) {
    hash := calculateHash(config)
    
    version := &ConfigVersion{
        ID:          generateID(),
        Version:     cm.getNextVersion(),
        Timestamp:   time.Now(),
        UserID:      userID,
        Description: description,
        ConfigData:  config,
        Hash:        hash,
        ParentID:    cm.getCurrentVersionID(),
    }
    
    if err := cm.db.SaveConfigVersion(version); err != nil {
        return nil, err
    }
    
    cm.currentHash = hash
    return version, nil
}

func (cm *ConfigManager) RestoreConfiguration(versionID string) error {
    version, err := cm.db.GetConfigVersion(versionID)
    if err != nil {
        return err
    }
    
    // Apply configuration to FRR
    if err := cm.applyConfiguration(version.ConfigData); err != nil {
        return err
    }
    
    // Create new version pointing to restored config
    return cm.SaveConfiguration(
        "system",
        fmt.Sprintf("Restored from version %d", version.Version),
        version.ConfigData,
    )
}
```

### Configuration Diff

```go
type ConfigDiff struct {
    Added    []ConfigLine
    Removed  []ConfigLine
    Modified []ConfigChange
}

type ConfigLine struct {
    LineNumber int
    Content    string
}

type ConfigChange struct {
    LineNumber int
    OldContent string
    NewContent string
}

func (cm *ConfigManager) DiffConfigurations(oldID, newID string) (*ConfigDiff, error) {
    oldVersion, err := cm.db.GetConfigVersion(oldID)
    if err != nil {
        return nil, err
    }
    
    newVersion, err := cm.db.GetConfigVersion(newID)
    if err != nil {
        return nil, err
    }
    
    return calculateDiff(oldVersion.ConfigData, newVersion.ConfigData), nil
}
```

---

## Real-time Updates

### WebSocket Message Protocol

```typescript
interface WebSocketMessage {
  type: 'subscribe' | 'unsubscribe' | 'event' | 'ping' | 'pong';
  channel?: string;
  data?: any;
  timestamp: number;
}

// Subscribe to BGP peer updates
const subscribeMessage: WebSocketMessage = {
  type: 'subscribe',
  channel: 'bgp.peers',
  timestamp: Date.now(),
};

// Event notification
const eventMessage: WebSocketMessage = {
  type: 'event',
  channel: 'bgp.peer.AS64512',
  data: {
    event: 'state_change',
    peer_id: 'peer_123',
    old_state: 'Active',
    new_state: 'Established',
    timestamp: '2024-01-15T10:30:45Z',
  },
  timestamp: Date.now(),
};
```

### Backend WebSocket Handler

```go
type WebSocketHandler struct {
    stateManager *StateManager
    pubsub       *PubSub
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Error("WebSocket upgrade failed", "error", err)
        return
    }
    defer conn.Close()
    
    userID := getUserIDFromContext(r.Context())
    h.stateManager.AddWSConnection(userID, conn)
    defer h.stateManager.RemoveWSConnection(userID)
    
    // Handle incoming messages
    for {
        var msg WebSocketMessage
        if err := conn.ReadJSON(&msg); err != nil {
            break
        }
        
        h.handleMessage(userID, &msg)
    }
}

func (h *WebSocketHandler) handleMessage(userID string, msg *WebSocketMessage) {
    switch msg.Type {
    case "subscribe":
        h.pubsub.Subscribe(userID, msg.Channel)
    case "unsubscribe":
        h.pubsub.Unsubscribe(userID, msg.Channel)
    case "ping":
        h.stateManager.BroadcastToUser(userID, &WebSocketMessage{
            Type:      "pong",
            Timestamp: time.Now().Unix(),
        })
    }
}
```

### Pub/Sub System

```go
type PubSub struct {
    subscribers map[string]map[string]bool // channel -> userID -> subscribed
    lock        sync.RWMutex
    stateManager *StateManager
}

func (ps *PubSub) Subscribe(userID, channel string) {
    ps.lock.Lock()
    defer ps.lock.Unlock()
    
    if ps.subscribers[channel] == nil {
        ps.subscribers[channel] = make(map[string]bool)
    }
    ps.subscribers[channel][userID] = true
}

func (ps *PubSub) Unsubscribe(userID, channel string) {
    ps.lock.Lock()
    defer ps.lock.Unlock()
    
    if ps.subscribers[channel] != nil {
        delete(ps.subscribers[channel], userID)
    }
}

func (ps *PubSub) Publish(channel string, message interface{}) {
    ps.lock.RLock()
    subscribers := ps.subscribers[channel]
    ps.lock.RUnlock()
    
    for userID := range subscribers {
        go ps.stateManager.BroadcastToUser(userID, &WebSocketMessage{
            Type:      "event",
            Channel:   channel,
            Data:      message,
            Timestamp: time.Now().Unix(),
        })
    }
}
```

---

## Caching Strategy

### Cache Layers

```
┌─────────────────────────────────────────────────────────┐
│                  Browser Cache                           │
│  - Static assets (CDN)                                  │
│  - API responses (React Query)                          │
│  - TTL: 30s - 5m                                        │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Backend Cache                           │
│  - Session data                                         │
│  - BGP state                                            │
│  - Configuration                                        │
│  - TTL: 1m - 10m                                        │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Database                                │
│  - Persistent data                                      │
│  - No TTL                                               │
└─────────────────────────────────────────────────────────┘
```

### Cache Implementation

```go
type Cache struct {
    data map[string]*CacheEntry
    lock sync.RWMutex
}

type CacheEntry struct {
    Value      interface{}
    Expiration time.Time
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.lock.Lock()
    defer c.lock.Unlock()
    
    c.data[key] = &CacheEntry{
        Value:      value,
        Expiration: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.lock.RLock()
    defer c.lock.RUnlock()
    
    entry, exists := c.data[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(entry.Expiration) {
        return nil, false
    }
    
    return entry.Value, true
}

func (c *Cache) Invalidate(key string) {
    c.lock.Lock()
    defer c.lock.Unlock()
    delete(c.data, key)
}

func (c *Cache) InvalidatePattern(pattern string) {
    c.lock.Lock()
    defer c.lock.Unlock()
    
    for key := range c.data {
        if matched, _ := filepath.Match(pattern, key); matched {
            delete(c.data, key)
        }
    }
}
```

### Cache Invalidation Strategy

```go
// Invalidate on configuration change
func (cm *ConfigManager) ApplyConfiguration(config []byte) error {
    if err := cm.applyToFRR(config); err != nil {
        return err
    }
    
    // Invalidate all BGP-related cache
    cm.cache.InvalidatePattern("bgp:*")
    
    // Broadcast configuration change
    cm.pubsub.Publish("config.changed", map[string]interface{}{
        "timestamp": time.Now(),
    })
    
    return nil
}

// Invalidate on BGP state change
func (l *FRREventListener) handleBGPStateChange(event *FRREvent) {
    peerID := event.Data["peer_id"].(string)
    
    // Invalidate specific peer cache
    l.cache.Invalidate(fmt.Sprintf("bgp:peer:%s", peerID))
    
    // Broadcast to subscribers
    l.pubsub.Publish(fmt.Sprintf("bgp.peer.%s", peerID), event.Data)
}
```

---

## Conflict Resolution

### Configuration Drift Handling

```go
type DriftHandler struct {
    configManager *ConfigManager
    alertManager  *AlertManager
}

func (dh *DriftHandler) HandleDrift(expected, actual []byte) error {
    // Calculate diff
    diff := calculateDiff(expected, actual)
    
    // Create alert
    alert := &Alert{
        Type:     "config_drift",
        Severity: "warning",
        Message:  "Configuration drift detected",
        Data: map[string]interface{}{
            "diff": diff,
        },
    }
    dh.alertManager.CreateAlert(alert)
    
    // Log to audit
    auditLog := &AuditLog{
        EventType: "config.drift_detected",
        Severity:  "warning",
        Data:      diff,
    }
    dh.configManager.LogAudit(auditLog)
    
    return nil
}

func (dh *DriftHandler) ResolveConflict(strategy string) error {
    switch strategy {
    case "revert_to_flintroute":
        return dh.revertToFlintRoute()
    case "accept_frr_config":
        return dh.acceptFRRConfig()
    case "manual":
        return nil // Wait for manual resolution
    default:
        return errors.New("unknown strategy")
    }
}
```

### Concurrent Modification Handling

```go
type OptimisticLock struct {
    Version int
    Hash    string
}

func (cm *ConfigManager) UpdateWithLock(configID string, lock *OptimisticLock, newConfig []byte) error {
    current, err := cm.db.GetConfigVersion(configID)
    if err != nil {
        return err
    }
    
    // Check version
    if current.Version != lock.Version {
        return errors.New("version mismatch: configuration was modified")
    }
    
    // Check hash
    if current.Hash != lock.Hash {
        return errors.New("hash mismatch: configuration was modified")
    }
    
    // Apply update
    return cm.SaveConfiguration("user", "Update", newConfig)
}
```

---

## State Persistence

### Database Schema

```sql
-- Configuration versions
CREATE TABLE config_versions (
    id TEXT PRIMARY KEY,
    version INTEGER NOT NULL,
    timestamp DATETIME NOT NULL,
    user_id TEXT NOT NULL,
    description TEXT,
    config_data BLOB NOT NULL,
    hash TEXT NOT NULL,
    parent_id TEXT,
    FOREIGN KEY (parent_id) REFERENCES config_versions(id)
);

-- BGP peer state cache
CREATE TABLE bgp_peer_state (
    peer_id TEXT PRIMARY KEY,
    state TEXT NOT NULL,
    uptime INTEGER,
    routes_received INTEGER,
    routes_advertised INTEGER,
    last_update DATETIME NOT NULL
);

-- Audit logs
CREATE TABLE audit_logs (
    id TEXT PRIMARY KEY,
    timestamp DATETIME NOT NULL,
    event_type TEXT NOT NULL,
    user_id TEXT,
    resource_type TEXT,
    resource_id TEXT,
    action TEXT NOT NULL,
    status TEXT NOT NULL,
    changes TEXT,
    metadata TEXT
);
```

### State Recovery

```go
func (sm *StateManager) RecoverState() error {
    // Recover configuration
    latestConfig, err := sm.db.GetLatestConfigVersion()
    if err != nil {
        return err
    }
    sm.currentConfigHash = latestConfig.Hash
    
    // Recover BGP state from FRR
    sessions, err := sm.grpcClient.GetBGPSessions()
    if err != nil {
        return err
    }
    
    for _, session := range sessions {
        sm.cache.Set(
            fmt.Sprintf("bgp:peer:%s", session.PeerID),
            session,
            5*time.Minute,
        )
    }
    
    return nil
}
```

---

## Best Practices

### 1. State Synchronization

- Poll FRR state every 30-60 seconds
- Use WebSocket for real-time updates
- Implement exponential backoff for retries
- Cache frequently accessed data

### 2. Consistency

- Use optimistic locking for concurrent updates
- Implement atomic transactions
- Validate state before applying changes
- Maintain audit trail

### 3. Performance

- Minimize database queries
- Use appropriate cache TTLs
- Batch updates when possible
- Lazy load data

### 4. Reliability

- Handle network failures gracefully
- Implement automatic reconnection
- Maintain state across restarts
- Regular state validation

---

## Next Steps

- [Architecture Diagrams](diagrams.md)
- [API Documentation](../api/grpc-services.md)
- [Development Guide](../development/setup.md)