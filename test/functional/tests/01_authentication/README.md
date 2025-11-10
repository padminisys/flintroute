# Authentication Tests

This directory contains functional tests for the FlintRoute authentication system.

## Test Files

### `01_login_test.go`

Comprehensive authentication test suite that verifies:

#### TestLogin
- **successful_login**: Validates login with correct credentials
  - Verifies access token generation
  - Verifies refresh token generation
  - Validates user information in response
  - Checks token expiration time

- **invalid_credentials**: Tests login failure with wrong credentials
  - Ensures proper error handling
  - Validates error messages

- **empty_username**: Tests validation for empty username
  - Ensures proper input validation
  - Validates error response

- **empty_password**: Tests validation for empty password
  - Ensures proper input validation
  - Validates error response

- **authenticated_request**: Tests authenticated API calls
  - Verifies token is stored correctly
  - Tests authentication state
  - Validates authenticated requests work

- **logout**: Tests logout functionality
  - Verifies logout succeeds
  - Ensures tokens are cleared
  - Validates authentication state after logout

#### TestHealthCheck
- **health_check_no_auth**: Tests health endpoint without authentication
  - Verifies public endpoint accessibility
  - Validates server is running

#### TestTokenRefresh
- **token_refresh**: Tests token refresh mechanism
  - Validates refresh token usage
  - Verifies new tokens are generated
  - Checks token expiration updates

- **invalid_refresh_token**: Tests refresh with invalid token
  - Ensures proper error handling
  - Validates security measures

## Running Tests

### Run All Authentication Tests
```bash
cd test/functional
go test ./tests/01_authentication/... -v
```

### Run Specific Test
```bash
cd test/functional
go test ./tests/01_authentication -run TestLogin -v
```

### Run with Detailed Logging
```bash
cd test/functional
go test ./tests/01_authentication/... -v 2>&1 | tee logs/auth-test.log
```

## Prerequisites

Before running these tests:

1. **FlintRoute Server Must Be Running**
   ```bash
   # From project root
   make run
   # Or
   go run cmd/flintroute/main.go
   ```
   Server should be accessible at `http://localhost:8080`

2. **Database Should Be Clean**
   ```bash
   cd test/functional
   ./scripts/cleanup-db.sh
   ```

3. **Admin User Fixture Must Exist**
   - File: `test/functional/fixtures/users/admin_user.yaml`
   - Should contain valid admin credentials

## Test Data

### Fixtures Used
- **admin_user.yaml**: Admin user credentials for authentication tests
  ```yaml
  username: admin
  password: admin123
  email: admin@flintroute.local
  role: admin
  active: true
  ```

## Expected Behavior

### Successful Test Run
When all tests pass, you should see:
```
=== RUN   TestLogin
=== RUN   TestLogin/successful_login
=== RUN   TestLogin/invalid_credentials
=== RUN   TestLogin/empty_username
=== RUN   TestLogin/empty_password
=== RUN   TestLogin/authenticated_request
=== RUN   TestLogin/logout
--- PASS: TestLogin (X.XXs)
=== RUN   TestHealthCheck
=== RUN   TestHealthCheck/health_check_no_auth
--- PASS: TestHealthCheck (X.XXs)
=== RUN   TestTokenRefresh
=== RUN   TestTokenRefresh/token_refresh
=== RUN   TestTokenRefresh/invalid_refresh_token
--- PASS: TestTokenRefresh (X.XXs)
PASS
```

### Test Logs
Logs are written to `test/functional/logs/test.log` and include:
- Test start/end timestamps
- Request/response details
- Authentication state changes
- Assertion results
- Error details (if any)

## Framework Verification

This test suite serves as verification that:

1. ✅ **Testing Framework Works**
   - Logger initialization succeeds
   - Test utilities are functional
   - Assertions work correctly

2. ✅ **API Client Works**
   - HTTP requests are sent correctly
   - Responses are parsed properly
   - Token management functions

3. ✅ **Fixture Loading Works**
   - YAML fixtures are loaded
   - Data is parsed correctly
   - Fixtures are accessible

4. ✅ **FlintRoute Server Integration**
   - Server endpoints are accessible
   - Authentication flow works
   - API responses match expected format

5. ✅ **Mock FRR Server Integration**
   - (Indirectly verified through server operation)

## Troubleshooting

### Test Compilation Errors
```bash
# Ensure dependencies are installed
cd test/functional
go mod download
go mod tidy
```

### Connection Refused Errors
- Verify FlintRoute server is running: `curl http://localhost:8080/health`
- Check server logs for errors
- Ensure port 8080 is not blocked

### Authentication Failures
- Verify admin user exists in database
- Check JWT secret configuration
- Review server authentication logs

### Fixture Loading Errors
- Verify fixture file exists: `test/functional/fixtures/users/admin_user.yaml`
- Check YAML syntax is valid
- Ensure file permissions are correct

## Integration with CI/CD

These tests can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Authentication Tests
  run: |
    # Start services
    make run &
    sleep 5
    
    # Run tests
    cd test/functional
    go test ./tests/01_authentication/... -v
    
    # Cleanup
    pkill flintroute
```

## Next Steps

After verifying authentication tests work:
1. Run peer management tests (`02_peer_management/`)
2. Run session management tests (`03_session_management/`)
3. Run configuration tests (`04_configuration/`)
4. Run alert tests (`05_alerts/`)
5. Run error handling tests (`06_error_handling/`)
6. Run workflow tests (`07_workflows/`)

## Contributing

When modifying authentication tests:
- Maintain test independence (each test should work standalone)
- Add both positive and negative test cases
- Update this README with new test descriptions
- Ensure tests clean up after themselves
- Follow existing naming conventions