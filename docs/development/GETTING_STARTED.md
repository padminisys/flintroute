# FlintRoute Development Guide

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Running in Development Mode](#running-in-development-mode)
4. [Testing with FRR](#testing-with-frr)
5. [Adding New Features](#adding-new-features)
6. [Code Style and Standards](#code-style-and-standards)
7. [Testing](#testing)
8. [Debugging](#debugging)

## Development Environment Setup

### Prerequisites

- **Go 1.21+**: Backend development
- **Node.js 18+**: Frontend development
- **npm**: Package management
- **Git**: Version control
- **Docker & Docker Compose**: For FRR testing (optional)
- **SQLite**: Database (usually pre-installed)
- **Make**: Build automation (optional)

### Initial Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd flintroute
   ```

2. **Install Go dependencies:**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Install frontend dependencies:**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Create configuration files:**
   ```bash
   cp configs/config.example.yaml configs/config.yaml
   cp frontend/.env.example frontend/.env
   ```

5. **Build the project:**
   ```bash
   make build
   # Or manually:
   go build -o bin/flintroute ./cmd/flintroute
   cd frontend && npm run build && cd ..
   ```

### IDE Setup

#### VS Code (Recommended)

Install these extensions:
- Go (golang.go)
- ESLint (dbaeumer.vscode-eslint)
- Prettier (esbenp.prettier-vscode)
- TypeScript and JavaScript Language Features

Recommended settings (`.vscode/settings.json`):
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "[typescript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[typescriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

## Project Structure

```
flintroute/
├── cmd/
│   └── flintroute/          # Main application entry point
│       └── main.go
├── internal/                # Private application code
│   ├── api/                 # HTTP API handlers
│   │   ├── server.go        # Server setup and routing
│   │   ├── auth_handlers.go # Authentication endpoints
│   │   ├── bgp_handlers.go  # BGP management endpoints
│   │   └── config_handlers.go # Configuration endpoints
│   ├── auth/                # Authentication logic
│   │   ├── jwt.go           # JWT token management
│   │   └── middleware.go    # Auth middleware
│   ├── bgp/                 # BGP service layer
│   │   └── service.go       # BGP business logic
│   ├── config/              # Configuration management
│   │   └── config.go        # Config loading and validation
│   ├── database/            # Database layer
│   │   └── database.go      # SQLite operations
│   ├── frr/                 # FRR integration
│   │   └── client.go        # FRR VTY client
│   ├── models/              # Data models
│   │   └── models.go        # Shared data structures
│   └── websocket/           # WebSocket support
│       ├── hub.go           # WebSocket hub
│       └── handler.go       # WebSocket handlers
├── frontend/                # React frontend
│   ├── src/
│   │   ├── pages/           # Page components
│   │   ├── services/        # API and WebSocket clients
│   │   ├── store/           # Redux state management
│   │   └── main.tsx         # Application entry
│   └── package.json
├── configs/                 # Configuration files
│   ├── config.yaml          # Main config (gitignored)
│   └── config.example.yaml  # Example config
├── docs/                    # Documentation
├── data/                    # Runtime data (gitignored)
│   └── flintroute.db        # SQLite database
└── Makefile                 # Build automation
```

### Key Components

#### Backend Architecture

- **API Layer** (`internal/api/`): HTTP handlers and routing
- **Service Layer** (`internal/bgp/`): Business logic
- **Data Layer** (`internal/database/`): Database operations
- **Integration Layer** (`internal/frr/`): External system integration

#### Frontend Architecture

- **Pages**: Top-level route components
- **Services**: API communication and WebSocket
- **Store**: Redux state management with slices
- **Components**: Reusable UI components (to be added)

## Running in Development Mode

### Backend Development

**Option 1: Using Go run (with hot reload via air)**

Install air for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/flintroute"
  bin = "tmp/main"
  include_ext = ["go"]
  exclude_dir = ["frontend", "tmp", "vendor"]
  args_bin = ["--config", "configs/config.yaml"]
```

Run with hot reload:
```bash
air
```

**Option 2: Manual run**
```bash
go run ./cmd/flintroute --config configs/config.yaml
```

**Option 3: Build and run**
```bash
go build -o bin/flintroute ./cmd/flintroute
./bin/flintroute --config configs/config.yaml
```

### Frontend Development

```bash
cd frontend
npm run dev
```

This starts Vite dev server with:
- Hot Module Replacement (HMR)
- Fast refresh
- TypeScript checking
- Available at `http://localhost:5173`

### Full Stack Development

**Terminal 1 - Backend:**
```bash
air  # or go run ./cmd/flintroute --config configs/config.yaml
```

**Terminal 2 - Frontend:**
```bash
cd frontend && npm run dev
```

**Terminal 3 - Logs/Testing:**
```bash
# Watch logs
tail -f logs/flintroute.log

# Or test API
curl http://localhost:8080/api/health
```

## Testing with FRR

### Using Docker Compose

The project includes a Docker Compose setup for testing with FRR.

1. **Start FRR container:**
   ```bash
   docker-compose up -d frr
   ```

2. **Configure FRR:**
   ```bash
   # Access FRR shell
   docker-compose exec frr vtysh

   # Example BGP configuration
   configure terminal
   router bgp 65001
   bgp router-id 10.0.0.1
   neighbor 10.0.0.2 remote-as 65002
   exit
   exit
   write memory
   ```

3. **View FRR logs:**
   ```bash
   docker-compose logs -f frr
   ```

4. **Stop FRR:**
   ```bash
   docker-compose down
   ```

### Using Local FRR Installation

1. **Install FRR on Debian:**
   ```bash
   sudo apt-get update
   sudo apt-get install frr frr-pythontools
   ```

2. **Enable BGP daemon:**
   Edit `/etc/frr/daemons`:
   ```
   bgpd=yes
   zebra=yes
   ```

3. **Start FRR:**
   ```bash
   sudo systemctl start frr
   sudo systemctl enable frr
   ```

4. **Configure BGP:**
   ```bash
   sudo vtysh
   configure terminal
   router bgp 65001
   bgp router-id 10.0.0.1
   neighbor 10.0.0.2 remote-as 65002
   exit
   exit
   write memory
   ```

5. **Update FlintRoute config:**
   Edit `configs/config.yaml`:
   ```yaml
   frr:
     vty_address: "localhost:2605"
     vty_password: ""  # Set if configured
   ```

## Adding New Features

### Backend Feature Development

#### 1. Add Data Model

Edit `internal/models/models.go`:
```go
type NewFeature struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}
```

#### 2. Add Database Operations

Edit `internal/database/database.go`:
```go
func (db *Database) CreateNewFeature(feature *models.NewFeature) error {
    query := `INSERT INTO new_features (name, created_at) VALUES (?, ?)`
    result, err := db.db.Exec(query, feature.Name, time.Now())
    // ... handle error and get ID
    return nil
}
```

#### 3. Add Service Layer

Create `internal/newfeature/service.go`:
```go
package newfeature

type Service struct {
    db *database.Database
}

func NewService(db *database.Database) *Service {
    return &Service{db: db}
}

func (s *Service) CreateFeature(name string) (*models.NewFeature, error) {
    // Business logic here
    return s.db.CreateNewFeature(&models.NewFeature{Name: name})
}
```

#### 4. Add API Handlers

Create `internal/api/newfeature_handlers.go`:
```go
func (s *Server) handleCreateNewFeature(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call service
    // Return response
}

// Register in server.go setupRoutes():
r.Post("/api/newfeature", s.handleCreateNewFeature)
```

#### 5. Add Tests

Create `internal/newfeature/service_test.go`:
```go
func TestCreateFeature(t *testing.T) {
    // Setup test database
    // Test service methods
    // Assert results
}
```

### Frontend Feature Development

#### 1. Add Redux Slice

Create `frontend/src/store/newFeatureSlice.ts`:
```typescript
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

export const fetchNewFeatures = createAsyncThunk(
  'newFeature/fetch',
  async () => {
    const response = await api.get('/api/newfeature');
    return response.data;
  }
);

const newFeatureSlice = createSlice({
  name: 'newFeature',
  initialState: { items: [], loading: false },
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchNewFeatures.fulfilled, (state, action) => {
      state.items = action.payload;
    });
  },
});

export default newFeatureSlice.reducer;
```

#### 2. Add API Service

Edit `frontend/src/services/api.ts`:
```typescript
export const newFeatureAPI = {
  getAll: () => api.get('/api/newfeature'),
  create: (data: any) => api.post('/api/newfeature', data),
};
```

#### 3. Create Page Component

Create `frontend/src/pages/NewFeature.tsx`:
```typescript
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchNewFeatures } from '../store/newFeatureSlice';

export const NewFeaturePage: React.FC = () => {
  const dispatch = useDispatch();
  const features = useSelector((state: RootState) => state.newFeature.items);

  useEffect(() => {
    dispatch(fetchNewFeatures());
  }, [dispatch]);

  return <div>{/* Render features */}</div>;
};
```

#### 4. Add Route

Edit `frontend/src/App.tsx`:
```typescript
import { NewFeaturePage } from './pages/NewFeature';

// In Routes:
<Route path="/newfeature" element={<NewFeaturePage />} />
```

## Code Style and Standards

### Go Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write meaningful comments for exported functions
- Keep functions small and focused
- Use meaningful variable names

Example:
```go
// CreateBGPSession creates a new BGP session with the given configuration.
// It validates the configuration, stores it in the database, and applies
// it to FRR if connected.
func (s *Service) CreateBGPSession(config *models.BGPSessionConfig) error {
    if err := s.validateConfig(config); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    // ... implementation
}
```

### TypeScript/React Style

- Use functional components with hooks
- Use TypeScript strict mode
- Follow Airbnb style guide
- Use Prettier for formatting
- Use ESLint for linting
- Prefer const over let
- Use meaningful component and variable names

Example:
```typescript
interface Props {
  sessionId: string;
  onUpdate: (session: BGPSession) => void;
}

export const BGPSessionCard: React.FC<Props> = ({ sessionId, onUpdate }) => {
  const [loading, setLoading] = useState(false);
  
  // ... implementation
};
```

## Testing

### Backend Testing

Run all tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run specific package:
```bash
go test ./internal/bgp/...
```

### Frontend Testing

```bash
cd frontend
npm test
```

Run with coverage:
```bash
npm test -- --coverage
```

### Integration Testing

```bash
# Start services
docker-compose up -d

# Run integration tests
go test -tags=integration ./tests/integration/...

# Cleanup
docker-compose down
```

## Debugging

### Backend Debugging

#### Using Delve

Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Debug with Delve:
```bash
dlv debug ./cmd/flintroute -- --config configs/config.yaml
```

#### VS Code Debug Configuration

Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch FlintRoute",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/flintroute",
      "args": ["--config", "configs/config.yaml"]
    }
  ]
}
```

### Frontend Debugging

- Use browser DevTools (F12)
- Use React DevTools extension
- Use Redux DevTools extension
- Check console for errors
- Use `debugger` statements

### Common Issues

#### Port Already in Use
```bash
# Find process using port 8080
sudo lsof -i :8080
# Kill process
kill -9 <PID>
```

#### Database Locked
```bash
# Close all connections
rm data/flintroute.db
# Restart application
```

#### FRR Connection Failed
```bash
# Check FRR status
sudo systemctl status frr
# Check VTY port
sudo netstat -tlnp | grep 2605
```

## Development Workflow

1. **Create feature branch:**
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make changes and test:**
   ```bash
   # Backend
   go test ./...
   
   # Frontend
   cd frontend && npm test
   ```

3. **Format code:**
   ```bash
   # Backend
   gofmt -w .
   
   # Frontend
   cd frontend && npm run format
   ```

4. **Commit changes:**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```

5. **Push and create PR:**
   ```bash
   git push origin feature/new-feature
   ```

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [Redux Toolkit](https://redux-toolkit.js.org/)
- [Material-UI](https://mui.com/)
- [FRR Documentation](https://docs.frrouting.org/)
- [Project Architecture](../architecture/overview.md)
- [Testing Guide](./testing.md)

## Getting Help

- Check existing documentation in `docs/`
- Review code comments and examples
- Ask questions in team chat or issue tracker
- Refer to external documentation links above