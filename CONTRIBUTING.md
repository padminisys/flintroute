# Contributing to FlintRoute

Thank you for your interest in contributing to FlintRoute! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Pull Request Process](#pull-request-process)
- [Community](#community)

---

## Code of Conduct

This project adheres to a Code of Conduct that all contributors are expected to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

**In short:**
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions

---

## Getting Started

### Prerequisites

- **Go**: 1.21 or higher
- **Node.js**: 18 or higher
- **Docker**: For running FRR test instances
- **Git**: For version control
- **Make**: For build automation

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   # Fork the repository on GitHub
   git clone https://github.com/YOUR_USERNAME/flintroute.git
   cd flintroute
   
   # Add upstream remote
   git remote add upstream https://github.com/padminisys/flintroute.git
   ```

2. **Install Dependencies**
   ```bash
   # Backend dependencies
   go mod download
   
   # Frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Set Up FRR Test Environment**
   ```bash
   # Start FRR in Docker
   docker-compose up -d frr
   
   # Verify FRR is running
   docker exec flintroute-frr vtysh -c "show version"
   ```

4. **Run Development Servers**
   ```bash
   # Terminal 1: Backend
   make dev-backend
   
   # Terminal 2: Frontend
   make dev-frontend
   ```

5. **Verify Setup**
   - Backend: http://localhost:8080/health
   - Frontend: http://localhost:3000
   - API Docs: http://localhost:8080/api/docs

See [Development Setup Guide](docs/development/setup.md) for detailed instructions.

---

## Development Workflow

### Branch Strategy

We use a simplified Git Flow:

- `main` - Production-ready code
- `develop` - Integration branch for features
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Urgent production fixes
- `docs/*` - Documentation updates

### Creating a Feature Branch

```bash
# Update your local repository
git checkout develop
git pull upstream develop

# Create feature branch
git checkout -b feature/your-feature-name

# Make your changes
# ...

# Commit your changes
git add .
git commit -m "feat: add your feature description"

# Push to your fork
git push origin feature/your-feature-name
```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements

**Examples:**
```bash
feat(bgp): add peer group support
fix(auth): resolve token refresh issue
docs(api): update gRPC service documentation
test(bgp): add integration tests for peer management
```

### Keeping Your Branch Updated

```bash
# Fetch latest changes
git fetch upstream

# Rebase your branch
git rebase upstream/develop

# If conflicts occur, resolve them and continue
git add .
git rebase --continue

# Force push to your fork (if already pushed)
git push origin feature/your-feature-name --force-with-lease
```

---

## Coding Standards

### Go Code Style

We follow standard Go conventions:

1. **Formatting**
   ```bash
   # Format code
   gofmt -w .
   
   # Or use goimports
   goimports -w .
   ```

2. **Linting**
   ```bash
   # Run golangci-lint
   golangci-lint run
   ```

3. **Naming Conventions**
   - Use camelCase for variables and functions
   - Use PascalCase for exported identifiers
   - Use descriptive names (avoid abbreviations)

4. **Error Handling**
   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to create peer: %w", err)
   }
   
   // Bad
   if err != nil {
       return err
   }
   ```

5. **Comments**
   ```go
   // Package-level comment
   package bgp
   
   // CreatePeer creates a new BGP peer with the given configuration.
   // It validates the configuration and applies it to FRR.
   // Returns an error if validation fails or FRR rejects the configuration.
   func CreatePeer(config *PeerConfig) error {
       // Implementation
   }
   ```

### TypeScript/React Code Style

1. **Formatting**
   ```bash
   # Format code
   npm run format
   
   # Check formatting
   npm run format:check
   ```

2. **Linting**
   ```bash
   # Run ESLint
   npm run lint
   
   # Fix auto-fixable issues
   npm run lint:fix
   ```

3. **Component Structure**
   ```typescript
   // Good: Functional component with TypeScript
   interface PeerListProps {
     peers: BGPPeer[];
     onPeerSelect: (peer: BGPPeer) => void;
   }
   
   export const PeerList: React.FC<PeerListProps> = ({ peers, onPeerSelect }) => {
     return (
       <div>
         {peers.map(peer => (
           <PeerCard key={peer.id} peer={peer} onClick={() => onPeerSelect(peer)} />
         ))}
       </div>
     );
   };
   ```

4. **Hooks**
   ```typescript
   // Custom hooks should start with 'use'
   function useBGPPeers() {
     return useQuery({
       queryKey: ['bgp', 'peers'],
       queryFn: fetchBGPPeers,
     });
   }
   ```

### Code Organization

```
flintroute/
├── cmd/                    # Application entry points
│   └── flintroute/
│       └── main.go
├── internal/               # Private application code
│   ├── api/               # API handlers
│   ├── auth/              # Authentication
│   ├── bgp/               # BGP management
│   ├── config/            # Configuration
│   └── frr/               # FRR integration
├── pkg/                    # Public libraries
│   └── models/            # Shared data models
├── frontend/              # React application
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── hooks/        # Custom hooks
│   │   ├── pages/        # Page components
│   │   ├── services/     # API services
│   │   └── store/        # Redux store
│   └── public/
└── docs/                  # Documentation
```

---

## Testing Guidelines

### Backend Testing

1. **Unit Tests**
   ```go
   func TestCreatePeer(t *testing.T) {
       tests := []struct {
           name    string
           config  *PeerConfig
           wantErr bool
       }{
           {
               name: "valid peer",
               config: &PeerConfig{
                   RemoteAS: 64512,
                   Neighbor: "192.0.2.1",
               },
               wantErr: false,
           },
           {
               name: "invalid AS",
               config: &PeerConfig{
                   RemoteAS: 0,
                   Neighbor: "192.0.2.1",
               },
               wantErr: true,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               err := CreatePeer(tt.config)
               if (err != nil) != tt.wantErr {
                   t.Errorf("CreatePeer() error = %v, wantErr %v", err, tt.wantErr)
               }
           })
       }
   }
   ```

2. **Run Tests**
   ```bash
   # Run all tests
   go test ./...
   
   # Run with coverage
   go test -cover ./...
   
   # Generate coverage report
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   ```

### Frontend Testing

1. **Component Tests**
   ```typescript
   import { render, screen } from '@testing-library/react';
   import { PeerList } from './PeerList';
   
   describe('PeerList', () => {
     it('renders peer list', () => {
       const peers = [
         { id: '1', name: 'Peer 1', remoteAS: 64512 },
         { id: '2', name: 'Peer 2', remoteAS: 64513 },
       ];
       
       render(<PeerList peers={peers} onPeerSelect={() => {}} />);
       
       expect(screen.getByText('Peer 1')).toBeInTheDocument();
       expect(screen.getByText('Peer 2')).toBeInTheDocument();
     });
   });
   ```

2. **Run Tests**
   ```bash
   # Run all tests
   npm test
   
   # Run with coverage
   npm test -- --coverage
   
   # Run E2E tests
   npm run test:e2e
   ```

### Test Coverage Requirements

- Backend: Minimum 80% coverage
- Frontend: Minimum 70% coverage
- Critical paths: 100% coverage

---

## Documentation

### Code Documentation

1. **Go Documentation**
   ```go
   // Package bgp provides BGP peer management functionality.
   //
   // It integrates with FRR's gRPC API to configure and monitor
   // BGP sessions.
   package bgp
   
   // PeerConfig represents the configuration for a BGP peer.
   type PeerConfig struct {
       // RemoteAS is the autonomous system number of the peer.
       RemoteAS uint32 `json:"remote_as"`
       
       // Neighbor is the IP address of the peer.
       Neighbor string `json:"neighbor"`
   }
   ```

2. **TypeScript Documentation**
   ```typescript
   /**
    * Fetches the list of BGP peers from the API.
    * 
    * @returns Promise resolving to array of BGP peers
    * @throws {APIError} If the request fails
    */
   export async function fetchBGPPeers(): Promise<BGPPeer[]> {
     // Implementation
   }
   ```

### User Documentation

- Update relevant documentation in `docs/`
- Include examples and use cases
- Add screenshots for UI changes
- Update API documentation for API changes

---

## Pull Request Process

### Before Submitting

1. **Ensure Tests Pass**
   ```bash
   make test
   ```

2. **Check Code Quality**
   ```bash
   make lint
   ```

3. **Update Documentation**
   - Update relevant docs in `docs/`
   - Update CHANGELOG.md
   - Update API documentation if needed

4. **Rebase on Latest**
   ```bash
   git fetch upstream
   git rebase upstream/develop
   ```

### Submitting Pull Request

1. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request**
   - Go to GitHub and create a PR
   - Use the PR template
   - Link related issues
   - Add screenshots for UI changes

3. **PR Title Format**
   ```
   feat(bgp): add peer group support
   ```

4. **PR Description Template**
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests added/updated
   - [ ] Integration tests added/updated
   - [ ] Manual testing performed
   
   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] Tests pass locally
   - [ ] No new warnings
   
   ## Related Issues
   Closes #123
   ```

### Review Process

1. **Automated Checks**
   - CI/CD pipeline must pass
   - Code coverage must meet requirements
   - Linting must pass

2. **Code Review**
   - At least one approval required
   - Address all review comments
   - Update PR based on feedback

3. **Merge**
   - Squash and merge (default)
   - Delete branch after merge

---

## Community

### Getting Help

- **GitHub Discussions**: Ask questions and discuss ideas
- **GitHub Issues**: Report bugs and request features
- **Discord/Slack**: Real-time chat (coming soon)
- **Email**: support@flintroute.com

### Reporting Bugs

Use the bug report template:

```markdown
**Describe the bug**
A clear description of the bug

**To Reproduce**
Steps to reproduce:
1. Go to '...'
2. Click on '...'
3. See error

**Expected behavior**
What you expected to happen

**Screenshots**
If applicable

**Environment**
- OS: [e.g., Debian 12]
- FlintRoute version: [e.g., 0.1.0]
- FRR version: [e.g., 8.4]

**Additional context**
Any other relevant information
```

### Feature Requests

Use the feature request template:

```markdown
**Is your feature request related to a problem?**
Description of the problem

**Describe the solution you'd like**
Clear description of desired functionality

**Describe alternatives you've considered**
Alternative solutions or features

**Additional context**
Any other relevant information
```

---

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project website (coming soon)

Thank you for contributing to FlintRoute! 🎉

---

## License

By contributing to FlintRoute, you agree that your contributions will be licensed under the Apache License 2.0.