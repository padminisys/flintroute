t # Testing Strategy Guide

## Table of Contents
- [Overview](#overview)
- [Testing Philosophy](#testing-philosophy)
- [Unit Testing](#unit-testing)
- [Integration Testing](#integration-testing)
- [End-to-End Testing](#end-to-end-testing)
- [Testing with Containerlab](#testing-with-containerlab)
- [CI/CD Testing Pipeline](#cicd-testing-pipeline)
- [Performance Testing](#performance-testing)
- [Security Testing](#security-testing)
- [Test Data Management](#test-data-management)
- [Best Practices](#best-practices)

---

## Overview

FlintRoute employs a comprehensive testing strategy covering all layers of the application, from unit tests to full end-to-end scenarios with real FRR instances.

### Testing Pyramid

```
           ┌─────────────┐
           │     E2E     │  ← 10% (Slow, Expensive)
           └─────────────┘
         ┌─────────────────┐
         │  Integration    │  ← 30% (Medium Speed)
         └─────────────────┘
       ┌─────────────────────┐
       │    Unit Tests       │  ← 60% (Fast, Cheap)
       └─────────────────────┘
```

### Test Coverage Goals

| Component | Target Coverage | Current |
|-----------|----------------|---------|
| Backend (Go) | 80%+ | TBD |
| Frontend (React) | 70%+ | TBD |
| Integration | 100% critical paths | TBD |
| E2E | Key user flows | TBD |

---

## Testing Philosophy

### Core Principles

1. **Test Behavior, Not Implementation**: Focus on what the code does, not how
2. **Fast Feedback**: Unit tests should run in seconds
3. **Reliable**: Tests should be deterministic and not flaky
4. **Maintainable**: Tests should be easy to understand and update
5. **Realistic**: Integration tests should use real FRR instances
6. **Automated**: All tests run in CI/CD pipeline

### Test Naming Convention

```go
// Go: TestFunctionName_Scenario_ExpectedBehavior
func TestCreateBGPPeer_ValidInput_ReturnsSuccess(t *testing.T) {}
func TestCreateBGPPeer_InvalidASN_ReturnsError(t *testing.T) {}
```

```typescript
// TypeScript: describe what, it should do what
describe('BGPPeerForm', () => {
  it('should validate AS number format', () => {})
  it('should submit form with valid data', () => {})
})
```

---

## Unit Testing

### Backend Unit Testing (Go)

#### Setup Testing Framework

```bash
cd backend

# Install testing dependencies
go get -u github.com/stretchr/testify/assert
go get -u github.com/stretchr/testify/mock
go get -u github.com/stretchr/testify/suite
```

#### Example: Service Layer Test

**File**: `internal/service/bgp_service_test.go`

```go
package service

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/padminisys/flintroute/internal/models"
)

// Mock FRR client
type MockFRRClient struct {
    mock.Mock
}

func (m *MockFRRClient) CreatePeer(ctx context.Context, peer *models.BGPPeer) error {
    args := m.Called(ctx, peer)
    return args.Error(0)
}

func TestBGPService_CreatePeer_ValidInput_Success(t *testing.T) {
    // Arrange
    mockFRR := new(MockFRRClient)
    mockRepo := new(MockRepository)
    service := NewBGPService(mockFRR, mockRepo)
    
    peer := &models.BGPPeer{
        RemoteAS:   65001,
        RemoteAddr: "192.0.2.1",
        Description: "Test Peer",
    }
    
    mockFRR.On("CreatePeer", mock.Anything, peer).Return(nil)
    mockRepo.On("Save", mock.Anything, peer).Return(nil)
    
    // Act
    err := service.CreatePeer(context.Background(), peer)
    
    // Assert
    assert.NoError(t, err)
    mockFRR.AssertExpectations(t)
    mockRepo.AssertExpectations(t)
}

func TestBGPService_CreatePeer_InvalidASN_Error(t *testing.T) {
    // Arrange
    service := NewBGPService(nil, nil)
    
    peer := &models.BGPPeer{
        RemoteAS:   0, // Invalid
        RemoteAddr: "192.0.2.1",
    }
    
    // Act
    err := service.CreatePeer(context.Background(), peer)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid AS number")
}

func TestBGPService_CreatePeer_DuplicatePeer_Error(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := NewBGPService(nil, mockRepo)
    
    peer := &models.BGPPeer{
        RemoteAS:   65001,
        RemoteAddr: "192.0.2.1",
    }
    
    mockRepo.On("FindByAddress", mock.Anything, "192.0.2.1").
        Return(&models.BGPPeer{}, nil)
    
    // Act
    err := service.CreatePeer(context.Background(), peer)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "peer already exists")
}
```

#### Example: Repository Test

**File**: `internal/repository/bgp_repository_test.go`

```go
package repository

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/padminisys/flintroute/internal/models"
)

type BGPRepositoryTestSuite struct {
    suite.Suite
    db   *gorm.DB
    repo *BGPRepository
}

func (suite *BGPRepositoryTestSuite) SetupTest() {
    // Create in-memory SQLite database
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    suite.NoError(err)
    
    // Auto-migrate schema
    err = db.AutoMigrate(&models.BGPPeer{})
    suite.NoError(err)
    
    suite.db = db
    suite.repo = NewBGPRepository(db)
}

func (suite *BGPRepositoryTestSuite) TearDownTest() {
    sqlDB, _ := suite.db.DB()
    sqlDB.Close()
}

func (suite *BGPRepositoryTestSuite) TestCreate_ValidPeer_Success() {
    // Arrange
    peer := &models.BGPPeer{
        RemoteAS:    65001,
        RemoteAddr:  "192.0.2.1",
        Description: "Test Peer",
    }
    
    // Act
    err := suite.repo.Create(context.Background(), peer)
    
    // Assert
    suite.NoError(err)
    suite.NotZero(peer.ID)
}

func (suite *BGPRepositoryTestSuite) TestFindByID_ExistingPeer_ReturnsPeer() {
    // Arrange
    peer := &models.BGPPeer{
        RemoteAS:   65001,
        RemoteAddr: "192.0.2.1",
    }
    suite.repo.Create(context.Background(), peer)
    
    // Act
    found, err := suite.repo.FindByID(context.Background(), peer.ID)
    
    // Assert
    suite.NoError(err)
    suite.Equal(peer.RemoteAS, found.RemoteAS)
    suite.Equal(peer.RemoteAddr, found.RemoteAddr)
}

func (suite *BGPRepositoryTestSuite) TestFindByID_NonExistent_ReturnsError() {
    // Act
    _, err := suite.repo.FindByID(context.Background(), 999)
    
    // Assert
    suite.Error(err)
}

func TestBGPRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(BGPRepositoryTestSuite))
}
```

#### Run Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Run specific package
go test ./internal/service/...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestBGPService_CreatePeer ./internal/service/
```

### Frontend Unit Testing (React)

#### Setup Testing Framework

```bash
cd frontend

# Install testing dependencies
npm install --save-dev \
    @testing-library/react \
    @testing-library/jest-dom \
    @testing-library/user-event \
    vitest \
    @vitest/ui \
    jsdom
```

#### Configure Vitest

**File**: `frontend/vitest.config.ts`

```typescript
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
      ],
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
})
```

#### Example: Component Test

**File**: `src/components/BGPPeerForm.test.tsx`

```typescript
import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { BGPPeerForm } from './BGPPeerForm'

describe('BGPPeerForm', () => {
  it('should render all form fields', () => {
    render(<BGPPeerForm onSubmit={vi.fn()} />)
    
    expect(screen.getByLabelText(/remote as/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/remote address/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /submit/i })).toBeInTheDocument()
  })

  it('should validate AS number format', async () => {
    render(<BGPPeerForm onSubmit={vi.fn()} />)
    
    const asnInput = screen.getByLabelText(/remote as/i)
    await userEvent.type(asnInput, 'invalid')
    
    fireEvent.blur(asnInput)
    
    await waitFor(() => {
      expect(screen.getByText(/invalid as number/i)).toBeInTheDocument()
    })
  })

  it('should validate IP address format', async () => {
    render(<BGPPeerForm onSubmit={vi.fn()} />)
    
    const ipInput = screen.getByLabelText(/remote address/i)
    await userEvent.type(ipInput, '999.999.999.999')
    
    fireEvent.blur(ipInput)
    
    await waitFor(() => {
      expect(screen.getByText(/invalid ip address/i)).toBeInTheDocument()
    })
  })

  it('should submit form with valid data', async () => {
    const onSubmit = vi.fn()
    render(<BGPPeerForm onSubmit={onSubmit} />)
    
    await userEvent.type(screen.getByLabelText(/remote as/i), '65001')
    await userEvent.type(screen.getByLabelText(/remote address/i), '192.0.2.1')
    await userEvent.type(screen.getByLabelText(/description/i), 'Test Peer')
    
    fireEvent.click(screen.getByRole('button', { name: /submit/i }))
    
    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        remoteAS: 65001,
        remoteAddr: '192.0.2.1',
        description: 'Test Peer',
      })
    })
  })

  it('should disable submit button while submitting', async () => {
    const onSubmit = vi.fn(() => new Promise(resolve => setTimeout(resolve, 100)))
    render(<BGPPeerForm onSubmit={onSubmit} />)
    
    await userEvent.type(screen.getByLabelText(/remote as/i), '65001')
    await userEvent.type(screen.getByLabelText(/remote address/i), '192.0.2.1')
    
    const submitButton = screen.getByRole('button', { name: /submit/i })
    fireEvent.click(submitButton)
    
    expect(submitButton).toBeDisabled()
    
    await waitFor(() => {
      expect(submitButton).not.toBeDisabled()
    })
  })
})
```

#### Example: Hook Test

**File**: `src/hooks/useBGPPeers.test.ts`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useBGPPeers } from './useBGPPeers'
import * as api from '@/services/api'

vi.mock('@/services/api')

describe('useBGPPeers', () => {
  let queryClient: QueryClient

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
      },
    })
  })

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  )

  it('should fetch BGP peers successfully', async () => {
    const mockPeers = [
      { id: 1, remoteAS: 65001, remoteAddr: '192.0.2.1' },
      { id: 2, remoteAS: 65002, remoteAddr: '192.0.2.2' },
    ]

    vi.mocked(api.getBGPPeers).mockResolvedValue(mockPeers)

    const { result } = renderHook(() => useBGPPeers(), { wrapper })

    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true)
    })

    expect(result.current.data).toEqual(mockPeers)
  })

  it('should handle fetch error', async () => {
    vi.mocked(api.getBGPPeers).mockRejectedValue(new Error('Network error'))

    const { result } = renderHook(() => useBGPPeers(), { wrapper })

    await waitFor(() => {
      expect(result.current.isError).toBe(true)
    })

    expect(result.current.error).toBeDefined()
  })
})
```

#### Run Frontend Tests

```bash
cd frontend

# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch

# Run with UI
npm run test:ui

# Run specific test file
npm test BGPPeerForm.test.tsx
```

---

## Integration Testing

### Backend Integration Tests

Integration tests verify that components work together correctly with real dependencies.

#### Setup Integration Test Environment

**File**: `backend/internal/integration/setup_test.go`

```go
package integration

import (
    "context"
    "testing"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type TestEnvironment struct {
    FRRContainer testcontainers.Container
    DBContainer  testcontainers.Container
}

func SetupTestEnvironment(t *testing.T) *TestEnvironment {
    ctx := context.Background()
    
    // Start FRR container
    frrReq := testcontainers.ContainerRequest{
        Image:        "frrouting/frr:v8.4.0",
        ExposedPorts: []string{"50051/tcp", "2605/tcp"},
        WaitingFor:   wait.ForListeningPort("50051/tcp"),
    }
    
    frrContainer, err := testcontainers.GenericContainer(ctx, 
        testcontainers.GenericContainerRequest{
            ContainerRequest: frrReq,
            Started:          true,
        })
    if err != nil {
        t.Fatalf("Failed to start FRR container: %v", err)
    }
    
    return &TestEnvironment{
        FRRContainer: frrContainer,
    }
}

func (env *TestEnvironment) Cleanup(t *testing.T) {
    ctx := context.Background()
    if env.FRRContainer != nil {
        env.FRRContainer.Terminate(ctx)
    }
}
```

#### Example: FRR Integration Test

**File**: `backend/internal/integration/frr_client_test.go`

```go
package integration

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/padminisys/flintroute/internal/frr"
    "github.com/padminisys/flintroute/internal/models"
)

func TestFRRClient_CreatePeer_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup
    env := SetupTestEnvironment(t)
    defer env.Cleanup(t)
    
    host, _ := env.FRRContainer.Host(context.Background())
    port, _ := env.FRRContainer.MappedPort(context.Background(), "50051")
    
    client, err := frr.NewClient(host, port.Port())
    require.NoError(t, err)
    defer client.Close()
    
    // Test
    peer := &models.BGPPeer{
        RemoteAS:    65001,
        RemoteAddr:  "192.0.2.1",
        Description: "Integration Test Peer",
    }
    
    err = client.CreatePeer(context.Background(), peer)
    assert.NoError(t, err)
    
    // Verify
    peers, err := client.ListPeers(context.Background())
    assert.NoError(t, err)
    assert.Len(t, peers, 1)
    assert.Equal(t, uint32(65001), peers[0].RemoteAS)
}

func TestFRRClient_SessionMonitoring_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    env := SetupTestEnvironment(t)
    defer env.Cleanup(t)
    
    host, _ := env.FRRContainer.Host(context.Background())
    port, _ := env.FRRContainer.MappedPort(context.Background(), "50051")
    
    client, err := frr.NewClient(host, port.Port())
    require.NoError(t, err)
    defer client.Close()
    
    // Get session status
    status, err := client.GetSessionStatus(context.Background())
    assert.NoError(t, err)
    assert.NotNil(t, status)
}
```

#### Run Integration Tests

```bash
cd backend

# Run integration tests
go test -tags=integration ./internal/integration/...

# Run all tests including integration
go test ./...

# Skip integration tests (for fast feedback)
go test -short ./...
```

### Frontend Integration Tests

**File**: `frontend/src/integration/BGPManagement.test.tsx`

```typescript
import { describe, it, expect, beforeAll, afterAll } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { setupServer } from 'msw/node'
import { rest } from 'msw'
import { App } from '@/App'

const server = setupServer(
  rest.get('/api/v1/bgp/peers', (req, res, ctx) => {
    return res(ctx.json([
      { id: 1, remoteAS: 65001, remoteAddr: '192.0.2.1', status: 'Established' },
    ]))
  }),
  
  rest.post('/api/v1/bgp/peers', (req, res, ctx) => {
    return res(ctx.json({ id: 2, ...req.body }))
  })
)

beforeAll(() => server.listen())
afterAll(() => server.close())

describe('BGP Management Integration', () => {
  it('should display list of peers and create new peer', async () => {
    render(<App />)
    
    // Wait for peers to load
    await waitFor(() => {
      expect(screen.getByText('192.0.2.1')).toBeInTheDocument()
    })
    
    // Click add peer button
    await userEvent.click(screen.getByRole('button', { name: /add peer/i }))
    
    // Fill form
    await userEvent.type(screen.getByLabelText(/remote as/i), '65002')
    await userEvent.type(screen.getByLabelText(/remote address/i), '192.0.2.2')
    
    // Submit
    await userEvent.click(screen.getByRole('button', { name: /submit/i }))
    
    // Verify new peer appears
    await waitFor(() => {
      expect(screen.getByText('192.0.2.2')).toBeInTheDocument()
    })
  })
})
```

---

## End-to-End Testing

### Setup Playwright

```bash
cd frontend

# Install Playwright
npm install --save-dev @playwright/test

# Install browsers
npx playwright install
```

#### Configure Playwright

**File**: `frontend/playwright.config.ts`

```typescript
import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
  ],
  webServer: {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
  },
})
```

#### Example: E2E Test

**File**: `frontend/e2e/bgp-management.spec.ts`

```typescript
import { test, expect } from '@playwright/test'

test.describe('BGP Peer Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: /login/i }).click()
    await page.getByLabel(/username/i).fill('admin')
    await page.getByLabel(/password/i).fill('admin')
    await page.getByRole('button', { name: /submit/i }).click()
    await expect(page.getByText(/dashboard/i)).toBeVisible()
  })

  test('should create new BGP peer', async ({ page }) => {
    // Navigate to BGP peers page
    await page.getByRole('link', { name: /bgp peers/i }).click()
    
    // Click add peer button
    await page.getByRole('button', { name: /add peer/i }).click()
    
    // Fill form
    await page.getByLabel(/remote as/i).fill('65001')
    await page.getByLabel(/remote address/i).fill('192.0.2.1')
    await page.getByLabel(/description/i).fill('Test Peer')
    
    // Submit
    await page.getByRole('button', { name: /create/i }).click()
    
    // Verify success message
    await expect(page.getByText(/peer created successfully/i)).toBeVisible()
    
    // Verify peer appears in list
    await expect(page.getByText('192.0.2.1')).toBeVisible()
    await expect(page.getByText('AS 65001')).toBeVisible()
  })

  test('should monitor BGP session status', async ({ page }) => {
    await page.getByRole('link', { name: /sessions/i }).click()
    
    // Wait for session data to load
    await expect(page.getByText(/session status/i)).toBeVisible()
    
    // Verify real-time updates (WebSocket)
    const statusBefore = await page.getByTestId('session-status').textContent()
    
    // Wait for potential status change
    await page.waitForTimeout(5000)
    
    const statusAfter = await page.getByTestId('session-status').textContent()
    
    // Status should be present (may or may not change)
    expect(statusAfter).toBeTruthy()
  })

  test('should backup and restore configuration', async ({ page }) => {
    await page.getByRole('link', { name: /configuration/i }).click()
    
    // Create backup
    await page.getByRole('button', { name: /backup/i }).click()
    await expect(page.getByText(/backup created/i)).toBeVisible()
    
    // Verify backup appears in list
    const backupRow = page.getByRole('row').filter({ hasText: /just now/i })
    await expect(backupRow).toBeVisible()
    
    // Restore backup
    await backupRow.getByRole('button', { name: /restore/i }).click()
    await page.getByRole('button', { name: /confirm/i }).click()
    
    await expect(page.getByText(/configuration restored/i)).toBeVisible()
  })
})
```

#### Run E2E Tests

```bash
cd frontend

# Run all E2E tests
npx playwright test

# Run in headed mode
npx playwright test --headed

# Run specific test
npx playwright test bgp-management

# Debug mode
npx playwright test --debug

# Generate report
npx playwright show-report
```

---

## Testing with Containerlab

### Automated Testing with Containerlab

**File**: `tests/containerlab/test_bgp_peering.sh`

```bash
#!/bin/bash
set -e

# Deploy test topology
echo "Deploying Containerlab topology..."
sudo containerlab deploy -t topology.yml

# Wait for FRR to start
echo "Waiting for FRR to initialize..."
sleep 30

# Run tests
echo "Running BGP peering tests..."

# Test 1: Verify BGP sessions are established
echo "Test 1: Checking BGP session status..."
sudo docker exec clab-test-rs1 vtysh -c "show bgp summary" | grep "Established" || {
    echo "FAIL: BGP sessions not established"
    exit 1
}
echo "PASS: BGP sessions established"

# Test 2: Verify routes are received
echo "Test 2: Checking received routes..."
ROUTES=$(sudo docker exec clab-test-rs1 vtysh -c "show bgp ipv4 unicast" | grep -c "192.168")
if [ "$ROUTES" -gt 0 ]; then
    echo "PASS: Received $ROUTES routes"
else
    echo "FAIL: No routes received"
    exit 1
fi

# Test 3: Test peer down detection
echo "Test 3: Testing peer down detection..."
sudo docker stop clab-test-peer1
sleep 10
sudo docker exec clab-test-rs1 vtysh -c "show bgp summary" | grep "peer1" | grep -v "Established" || {
    echo "FAIL: Peer down not detected"
    exit 1
}
echo "PASS: Peer down detected"

# Restart peer
sudo docker start clab-test-peer1
sleep 30

# Cleanup
echo "Cleaning up..."
sudo containerlab destroy -t topology.yml

echo "All tests passed!"
```

Make executable and run:

```bash
chmod +x tests/containerlab/test_bgp_peering.sh
./tests/containerlab/test_bgp_peering.sh
```

---

## CI/CD Testing Pipeline

### GitHub Actions Workflow

**File**: `.github/workflows/test.yml`

```yaml
name: Test

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
      
      - name: Install dependencies
        working-directory: ./backend
        run: go mod download
      
      - name: Run unit tests
        working-directory: ./backend
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out
          flags: backend

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
      
      - name: Run unit tests
        working-directory: ./frontend
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./frontend/coverage/coverage-final.json
          flags: frontend

  integration-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install Containerlab
        run: |
          bash -c "$(curl -sL https://get.containerlab.dev)"
      
      - name: Run integration tests
        working-directory: ./backend
        run: go test -tags=integration ./internal/integration/...

  e2e-test:
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
      
      - name: Install Playwright browsers
        working-directory: ./frontend
        run: npx playwright install --with-deps
      
      - name: Run E2E tests
        working-directory: ./frontend
        run: npx playwright test
      
      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

---

## Performance Testing

### Load Testing Backend API

**File**: `tests/performance/load_test.go`

```go
package performance

import (
    "context"
    "fmt"
    "sync"
    "testing"
    "time"
)

func BenchmarkCreatePeer(b *testing.B) {
    client := setupTestClient()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        peer := &models.BGPPeer{
            RemoteAS:   uint32(65000 + i),
            RemoteAddr: fmt.Sprintf("192.0.2.%d", i%255),
        }
        client.CreatePeer(context.Background(), peer)
    }
}

func TestConcurrentPeerCreation(t *testing.T) {
    client := setupTestClient()
    concurrency := 100
    
    var wg sync.WaitGroup
    start := time.Now()
    
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            peer := &models.BGPPeer{
                RemoteAS:   uint32(65000 + id),
                RemoteAddr: fmt.Sprintf("192.0.2.%d", id%255),
            }
            client.CreatePeer(context.Background(), peer)
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    t.Logf("Created %d peers in %v", concurrency, duration)
    t.Logf("Average: %v per peer", duration/time.Duration(concurrency))
}
```

### Load Testing with k6

**File**: `tests/performance/load_test.js`

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 20 },  // Ramp up
    { duration: '1m', target: 20 },   // Stay at 20 users
    { duration: '30s', target: 0 },   // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests under 500ms
    http_req_failed: ['rate<0.01'],   // Less than 1% failures
  },
};

export default function () {
  // Login
  const loginRes = http.post('http://localhost:8080/api/v1/auth/login', {
    username: 'admin',
    password: 'admin',
  });
  
  check(loginRes, {
    'login successful': (r) => r.status === 200,
  });
  
  const token = loginRes.json('token');
  const headers = { Authorization: `Bearer ${token}` };
  
  // List peers
  const listRes = http.get('http://localhost:8080/api/v1/bgp/peers', { headers });
  check(listRes, {
    'list peers successful': (r) => r.status === 200,
  });
  
  sleep(1);
}
```

Run k6 tests:

```bash
# Install k6
sudo apt-get install k6

# Run load test
k6 run tests/performance/load_test.js
```

---

## Security Testing

### Static Analysis

```bash
# Go security scanning
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# Dependency vulnerability scanning
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Frontend security scanning
npm audit
npm audit fix
```

### OWASP ZAP Scanning

```bash
# Run ZAP in Docker
docker run -t owasp/zap2docker-stable zap-baseline.py \
    -t http://localhost:8080 \
    -r zap-report.html
```

### Penetration Testing Checklist

- [ ] SQL Injection testing
- [ ] XSS testing
- [ ] CSRF testing
- [ ] Authentication bypass attempts
- [ ] Authorization bypass attempts
- [ ] Rate limiting verification
- [ ] Input validation testing
- [ ] Session management testing

---

## Test Data Management

### Test Fixtures

**File**: `backend/internal/testdata/fixtures.go`

```go
package testdata

import "github.com/padminisys/flintroute/internal/models"

func GetTestPeers() []*models.BGPPeer {
    return []*models.BGPPeer{
        {
            ID:          1,
            RemoteAS:    65001,
            RemoteAddr:  "192.0.2.1",
            Description: "Test Peer 1",
        },
        {
            ID:          2,
            RemoteAS:    65002,
            RemoteAddr:  "192.0.2.2",
            Description: "Test Peer 2",
        },
    }
}

func GetTestUser() *models.User {
    return &models.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Role:     "admin",
    }
}
```

### Database Seeding

**File**: `backend/internal/testdata/seed.go`

```go
package testdata

import (
    "gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
    // Seed users
    users := GetTestUsers()
    for _, user := range users {
        if err := db.Create(user).Error; err != nil {
            return err
        }
    }
    
    // Seed BGP peers
    peers := GetTestPeers()
    for _, peer := range peers {
        if err := db.Create(peer).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## Best Practices

### Test Organization

```
backend/
├── internal/
│   ├── service/
│   │   ├── bgp_service.go
│   │   └── bgp_service_test.go      # Unit tests
│   ├── integration/
│   │   └── bgp_integration_test.go  # Integration tests
│   └── testdata/
│       └── fixtures.go               # Test data
└── tests/
    ├── e2e/                          # E2E tests
    └── performance/                  # Performance tests

frontend/
├── src/
│   ├── components/
│   │   ├── BGPPeerForm.tsx
│   │   └── BGPPeerForm.test.tsx     # Unit tests
│   └── integration/
│       └── BGPManagement.test.tsx   # Integration tests
└── e2e/
    └── bgp-management.spec.ts       # E2E tests
```

### Test Naming

```go
// Good
func TestCreateBGPPeer_ValidInput_ReturnsSuccess(t *testing.T) {}
func TestCreateBGPPeer_InvalidASN_ReturnsError(t *testing.T) {}

// Bad
func TestCreatePeer(t *testing.T) {}
func TestPeer(t *testing.T) {}
```

### Test Independence

```go
// Good - Each test is independent
func TestA(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
    // Test logic
}

func TestB(t *testing.T) {
    db := setupTestDB()
    defer db.Close()
    // Test logic
}

// Bad - Tests depend on each other
var sharedDB *gorm.DB

func TestA(t *testing.T) {
    sharedDB = setupTestDB()
    // Test logic
}

func TestB(t *testing.T) {
    // Uses sharedDB from TestA
}
```

### Mocking Guidelines

1. **Mock External Dependencies**: Always mock external services (FRR, databases in unit tests)
2. **Use Interfaces**: Design code with interfaces for easy mocking
3. **Verify Interactions**: Use mock assertions to verify behavior
4. **Keep Mocks Simple**: Don't over-complicate mock implementations

### Test Coverage Goals

- **Critical Paths**: 100% coverage
- **Business Logic**: 90%+ coverage
- **Utilities**: 80%+ coverage
- **UI Components**: 70%+ coverage

### Continuous Improvement

1. **Review Test Failures**: Investigate and fix flaky tests immediately
2. **Refactor Tests**: Keep tests maintainable and readable
3. **Update Tests**: Update tests when requirements change
4. **Monitor Coverage**: Track coverage trends over time
5. **Performance**: Keep test suite fast (< 5 minutes for full suite)

---

## Quick Reference

### Run All Tests

```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && npm test

# Integration
go test -tags=integration ./...

# E2E
cd frontend && npx playwright test

# Performance
k6 run tests/performance/load_test.js
```

### Coverage Reports

```bash
# Backend coverage
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Frontend coverage
cd frontend
npm run test:coverage
open coverage/index.html
```

### Debug Tests

```bash
# Go debug
dlv test ./internal/service -- -test.run TestCreatePeer

# Frontend debug
npm test -- --inspect-brk

# Playwright debug
npx playwright test --debug
```

---

## Next Steps

1. **Implement Tests**: Start with unit tests for core functionality
2. **Set Up CI/CD**: Configure automated testing in GitHub Actions
3. **Integration Testing**: Set up Containerlab for integration tests
4. **E2E Testing**: Implement critical user flows with Playwright
5. **Monitor Coverage**: Track and improve test coverage over time

---

**Last Updated**: 2024-01-15  
**Version**: 0.1.0-alpha