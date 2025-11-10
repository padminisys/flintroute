# Contributing to FlintRoute Testing Framework

Thank you for your interest in contributing to the FlintRoute Functional Testing Framework! This document provides guidelines and instructions for contributing.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
4. [Development Workflow](#development-workflow)
5. [Coding Standards](#coding-standards)
6. [Testing Guidelines](#testing-guidelines)
7. [Documentation](#documentation)
8. [Pull Request Process](#pull-request-process)
9. [Issue Guidelines](#issue-guidelines)

---

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors. We expect all participants to:

- Be respectful and considerate
- Accept constructive criticism gracefully
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment or discriminatory language
- Personal attacks or trolling
- Publishing others' private information
- Other conduct that could reasonably be considered inappropriate

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- Go 1.21 or later installed
- Git for version control
- Basic understanding of Go testing
- Familiarity with BGP concepts (helpful but not required)

### Setting Up Development Environment

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/flintroute.git
   cd flintroute/test/functional
   ```

2. **Check prerequisites**
   ```bash
   ./scripts/check-prerequisites.sh
   ```

3. **Run existing tests**
   ```bash
   ./run-clean.sh
   ```

4. **Verify everything works**
   ```bash
   # All tests should pass
   echo $?  # Should output 0
   ```

---

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

#### 1. Bug Reports
- Found a bug? Report it!
- Include steps to reproduce
- Provide expected vs actual behavior
- Include relevant logs

#### 2. Feature Requests
- Suggest new test scenarios
- Propose improvements to existing tests
- Request new documentation

#### 3. Code Contributions
- Add new test cases
- Improve existing tests
- Enhance mock server functionality
- Optimize test execution

#### 4. Documentation
- Fix typos or unclear explanations
- Add examples
- Improve guides
- Translate documentation

#### 5. Test Fixtures
- Add new test data
- Improve existing fixtures
- Add edge case scenarios

---

## Development Workflow

### 1. Create a Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

### Branch Naming Conventions

- `feature/` - New features or enhancements
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `test/` - Test additions or improvements
- `refactor/` - Code refactoring

### 2. Make Changes

Follow the coding standards and guidelines in this document.

### 3. Test Your Changes

```bash
# Run all tests
./run-clean.sh

# Run specific tests
./run-tests.sh --pattern ./tests/YOUR_TEST/...

# Run with verbose output
./run-tests.sh --verbose --log-level debug
```

### 4. Commit Changes

```bash
# Stage changes
git add .

# Commit with descriptive message
git commit -m "Add: Brief description of changes"
```

### Commit Message Format

```
Type: Brief description (50 chars or less)

Detailed explanation of what changed and why.
Include any relevant issue numbers.

Fixes #123
```

**Types:**
- `Add:` - New feature or test
- `Fix:` - Bug fix
- `Update:` - Update existing functionality
- `Refactor:` - Code refactoring
- `Docs:` - Documentation changes
- `Test:` - Test additions or changes

### 5. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name

# Create pull request on GitHub
```

---

## Coding Standards

### Go Code Style

#### General Guidelines

1. **Follow Go conventions**
   - Use `gofmt` for formatting
   - Follow effective Go guidelines
   - Use meaningful variable names

2. **Code organization**
   ```go
   // Package declaration
   package mytest_test
   
   // Imports (grouped: stdlib, external, internal)
   import (
       "testing"
       "time"
       
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       
       "github.com/yourusername/flintroute/test/functional/pkg/testutil"
   )
   
   // Constants
   const (
       defaultTimeout = 30 * time.Second
   )
   
   // Test functions
   func TestFeature(t *testing.T) {
       // Test implementation
   }
   ```

3. **Naming conventions**
   - Test files: `feature_test.go`
   - Test functions: `TestFeatureName`
   - Subtests: `Feature_Scenario`
   - Variables: `camelCase`
   - Constants: `camelCase` or `UPPER_CASE`

#### Test Structure

```go
func TestFeature(t *testing.T) {
    // Setup
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    // Arrange - Prepare test data
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    
    // Act - Execute operation
    err := ctx.Client.CreatePeer(peer)
    
    // Assert - Verify results
    require.NoError(t, err)
    assert.NotNil(t, peer)
}
```

#### Error Handling

```go
// ✅ Good - Check and handle errors
result, err := ctx.Client.DoSomething()
require.NoError(t, err, "DoSomething should not fail")
assert.NotNil(t, result)

// ❌ Bad - Ignore errors
result, _ := ctx.Client.DoSomething()
```

#### Comments

```go
// ✅ Good - Explain why, not what
// Wait for session to establish because BGP negotiation takes time
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)

// ❌ Bad - State the obvious
// Call WaitForSessionState
err := testutil.WaitForSessionState(t, ctx, "peer1", "Established", 10*time.Second)
```

### Shell Script Style

1. **Use shellcheck**
   ```bash
   shellcheck script.sh
   ```

2. **Error handling**
   ```bash
   set -euo pipefail  # Exit on error, undefined vars, pipe failures
   ```

3. **Functions**
   ```bash
   # Function with clear purpose
   check_prerequisites() {
       local required_cmd=$1
       
       if ! command -v "$required_cmd" &> /dev/null; then
           print_error "$required_cmd not found"
           return 1
       fi
       
       return 0
   }
   ```

4. **Variables**
   ```bash
   # Use uppercase for constants
   readonly REQUIRED_GO_VERSION="1.21"
   
   # Use lowercase for local variables
   local test_result=0
   ```

### YAML Style

1. **Indentation**
   - Use 2 spaces (no tabs)
   - Consistent indentation

2. **Structure**
   ```yaml
   # Clear, hierarchical structure
   name: peer1
   description: Test BGP peer
   remote_ip: 192.168.1.1
   remote_asn: 65001
   enabled: true
   ```

3. **Comments**
   ```yaml
   # Explain purpose of fixture
   name: peer-with-authentication
   description: Peer with MD5 authentication
   password: test-password-123  # Test password only
   ```

---

## Testing Guidelines

### Writing Tests

#### 1. Test Independence

Each test must be completely independent:

```go
// ✅ Good - Independent test
func TestCreatePeer(t *testing.T) {
    ctx := testutil.NewTestContext(t)
    defer ctx.Cleanup()
    
    peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")
    err := ctx.Client.CreatePeer(peer)
    assert.NoError(t, err)
}

// ❌ Bad - Depends on other tests
func TestUpdatePeer(t *testing.T) {
    // Assumes peer from TestCreatePeer exists
    peer := &models.BGPPeer{Name: "peer1"}
    err := ctx.Client.UpdatePeer(peer)
    assert.NoError(t, err)
}
```

#### 2. Test Coverage

Test both success and failure cases:

```go
func TestPeerCreation(t *testing.T) {
    t.Run("Valid_Peer", func(t *testing.T) {
        // Test success case
    })
    
    t.Run("Invalid_IP", func(t *testing.T) {
        // Test failure case
    })
    
    t.Run("Duplicate_Name", func(t *testing.T) {
        // Test conflict case
    })
}
```

#### 3. Descriptive Names

Use clear, descriptive test names:

```go
// ✅ Good
func TestLoginWithValidCredentials(t *testing.T)
func TestCreatePeerWithInvalidIP(t *testing.T)
func TestSessionEstablishmentTimeout(t *testing.T)

// ❌ Bad
func TestLogin(t *testing.T)
func TestPeer(t *testing.T)
func TestSession(t *testing.T)
```

#### 4. Use Fixtures

Always use fixtures for test data:

```go
// ✅ Good - Use fixtures
peer := testutil.LoadPeerFixture(t, "fixtures/peers/valid/basic_peer.yaml")

// ❌ Bad - Hardcode data
peer := &models.BGPPeer{
    Name:      "test-peer",
    RemoteIP:  "192.168.1.1",
    RemoteASN: 65001,
}
```

### Adding New Test Suites

1. **Create directory**
   ```bash
   mkdir -p tests/08_new_feature
   ```

2. **Add test files**
   ```bash
   touch tests/08_new_feature/feature_test.go
   ```

3. **Follow naming convention**
   - Use numbered prefix (08_, 09_, etc.)
   - Use descriptive name
   - Add README.md explaining the suite

4. **Update documentation**
   - Add to TESTING_GUIDE.md
   - Update QUICK_REFERENCE.md
   - Update this CONTRIBUTING.md

### Adding Test Fixtures

1. **Choose appropriate directory**
   ```bash
   # For valid fixtures
   fixtures/peers/valid/
   
   # For invalid fixtures
   fixtures/peers/invalid/
   ```

2. **Create YAML file**
   ```yaml
   # fixtures/peers/valid/my_new_peer.yaml
   name: my-new-peer
   description: Description of this fixture
   remote_ip: 192.168.1.1
   remote_asn: 65001
   enabled: true
   ```

3. **Document the fixture**
   - Add to fixtures/README.md
   - Explain its purpose
   - Provide usage example

---

## Documentation

### Documentation Standards

1. **Clarity**
   - Write for beginners
   - Explain concepts clearly
   - Use examples

2. **Completeness**
   - Cover all features
   - Include edge cases
   - Provide troubleshooting

3. **Accuracy**
   - Keep documentation up-to-date
   - Test all examples
   - Verify commands work

### Documentation Structure

```markdown
# Title

Brief description

## Section

Detailed explanation

### Subsection

Specific details

#### Code Example

```language
code here
```

**Explanation**: What the code does
```

### Adding Documentation

1. **Update existing docs**
   - Fix typos
   - Clarify explanations
   - Add missing information

2. **Add new sections**
   - Follow existing structure
   - Use consistent formatting
   - Include examples

3. **Update references**
   - Update table of contents
   - Update cross-references
   - Update quick reference

---

## Pull Request Process

### Before Submitting

1. **Run all tests**
   ```bash
   ./run-clean.sh
   ```

2. **Check code style**
   ```bash
   gofmt -w .
   go vet ./...
   ```

3. **Update documentation**
   - Update relevant docs
   - Add examples if needed
   - Update CHANGELOG.md

4. **Write clear commit messages**
   - Follow commit message format
   - Reference related issues

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Test improvement

## Testing
- [ ] All tests pass
- [ ] Added new tests
- [ ] Updated existing tests

## Documentation
- [ ] Updated relevant documentation
- [ ] Added code examples
- [ ] Updated CHANGELOG.md

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] No new warnings generated

## Related Issues
Fixes #123
Related to #456
```

### Review Process

1. **Automated checks**
   - CI/CD pipeline runs
   - Tests must pass
   - Code style checks

2. **Code review**
   - At least one approval required
   - Address review comments
   - Update as needed

3. **Merge**
   - Squash commits if needed
   - Update branch if needed
   - Merge when approved

---

## Issue Guidelines

### Reporting Bugs

Use this template:

```markdown
## Bug Description
Clear description of the bug

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.21.5]
- FlintRoute version: [e.g., 1.0.0]

## Logs
```
Relevant log output
```

## Additional Context
Any other relevant information
```

### Requesting Features

Use this template:

```markdown
## Feature Description
Clear description of the feature

## Use Case
Why is this feature needed?

## Proposed Solution
How should it work?

## Alternatives Considered
Other approaches considered

## Additional Context
Any other relevant information
```

### Asking Questions

Use this template:

```markdown
## Question
Clear question

## Context
What are you trying to do?

## What You've Tried
Steps already attempted

## Environment
Relevant environment details
```

---

## Recognition

### Contributors

All contributors will be:
- Listed in project contributors
- Mentioned in release notes
- Credited in documentation

### Types of Recognition

- **Code Contributors**: Listed in GitHub contributors
- **Documentation Contributors**: Credited in docs
- **Bug Reporters**: Mentioned in CHANGELOG
- **Feature Requesters**: Credited in release notes

---

## Getting Help

### Resources

- **Documentation**: [TESTING_GUIDE.md](TESTING_GUIDE.md)
- **API Reference**: [API_REFERENCE.md](API_REFERENCE.md)
- **Quick Reference**: [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **FAQ**: [FAQ.md](FAQ.md)

### Communication

- **GitHub Issues**: For bugs and features
- **GitHub Discussions**: For questions and ideas
- **Pull Requests**: For code contributions

---

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

## Thank You!

Thank you for contributing to the FlintRoute Testing Framework! Your contributions help make this project better for everyone.

---

**Document Version**: 1.0.0  
**Last Updated**: November 10, 2025  
**Maintained By**: FlintRoute Development Team