# Temporary Files

This directory contains temporary files generated during test execution.

## Contents

### Test Databases
SQLite database files created for each test run:
- `test.db` - Main test database
- `test.db-shm` - Shared memory file (SQLite)
- `test.db-wal` - Write-ahead log (SQLite)

Each test suite may create isolated databases:
- `test-auth-*.db` - Authentication tests
- `test-peer-*.db` - Peer management tests
- `test-session-*.db` - Session management tests

### Temporary Data Files
- `*.tmp` - Temporary data files
- `*.cache` - Cache files
- `*.lock` - Lock files for synchronization

### Test Artifacts
- Configuration snapshots
- State dumps for debugging
- Intermediate test data

## Database Lifecycle

### Creation
Databases are created at test startup:
```go
db, err := sql.Open("sqlite3", "./tmp/test.db")
```

### Isolation
Each test suite can use isolated databases:
```go
dbPath := fmt.Sprintf("./tmp/test-%s-%d.db", suiteName, time.Now().Unix())
```

### Cleanup
Databases are cleaned up based on test configuration:

**Automatic Cleanup** (default)
- Removed after successful test completion
- Retained on test failure for debugging

**Manual Cleanup**
```bash
# Remove all test databases
rm -f tmp/*.db*

# Remove specific test artifacts
rm -f tmp/test-auth-*.db*
```

## Configuration

Temporary file settings in [`../config/test-config.yaml`](../config/test-config.yaml):

```yaml
database:
  path: ./tmp/test.db

testing:
  cleanup_on_success: false  # Keep files for inspection
```

## Disk Space Management

### Monitoring
```bash
# Check directory size
du -sh tmp/

# List large files
du -h tmp/* | sort -rh | head -10

# Count database files
ls -1 tmp/*.db 2>/dev/null | wc -l
```

### Cleanup Strategies

**After Each Test Run**
```bash
make clean-test-artifacts
```

**Scheduled Cleanup** (keep last 5 runs)
```bash
ls -t tmp/test-*.db | tail -n +6 | xargs rm -f
```

**Emergency Cleanup** (remove all)
```bash
rm -rf tmp/*
```

## Debugging with Temporary Files

### Inspecting Test Databases
```bash
# Open database
sqlite3 tmp/test.db

# List tables
.tables

# Query data
SELECT * FROM peers;
SELECT * FROM users;
SELECT * FROM sessions;

# Export schema
.schema > tmp/schema.sql
```

### Analyzing Failed Tests
When tests fail, databases are retained:

```bash
# Find databases from failed runs
ls -lt tmp/test-*.db | head -5

# Compare with expected state
sqlite3 tmp/test-failed.db "SELECT * FROM peers;" > actual.txt
diff expected.txt actual.txt
```

### State Dumps
Tests may create state dumps for debugging:

```bash
# View state dump
cat tmp/state-dump-*.json | jq .

# Compare states
diff tmp/state-before.json tmp/state-after.json
```

## Performance Considerations

### Database Performance
- SQLite performs well for test workloads
- WAL mode enabled for better concurrency
- In-memory databases for faster tests (if configured)

### Disk I/O
- Temporary files use local disk (fast)
- Consider RAM disk for CI environments:
  ```bash
  # Mount RAM disk (Linux)
  sudo mount -t tmpfs -o size=512M tmpfs tmp/
  ```

### Cleanup Impact
- Automatic cleanup adds minimal overhead
- Manual cleanup recommended for CI to save space

## CI Integration

### Artifact Retention
Failed test databases uploaded as CI artifacts:

```yaml
# GitHub Actions
- name: Upload test databases
  if: failure()
  uses: actions/upload-artifact@v3
  with:
    name: test-databases
    path: test/functional/tmp/*.db
```

### Space Management
CI environments have limited disk space:

```yaml
# Cleanup before tests
- name: Clean temporary files
  run: rm -rf test/functional/tmp/*

# Cleanup after tests
- name: Cleanup
  if: always()
  run: make clean-test-artifacts
```

## Security Considerations

### Sensitive Data
Temporary files may contain:
- Test user credentials
- API tokens
- Configuration secrets

**Protection Measures:**
- Files are local to test environment
- Cleaned up after test completion
- Not committed to version control (see `.gitignore`)
- Restricted file permissions (0600)

### Data Isolation
Each test run uses isolated databases:
- Prevents test interference
- Ensures clean state
- Enables parallel execution

## Troubleshooting

### Database Locked Errors
```
database is locked
```

**Solutions:**
- Ensure previous tests completed
- Check for orphaned processes
- Remove lock files: `rm -f tmp/*.db-shm tmp/*.db-wal`

### Disk Space Issues
```
no space left on device
```

**Solutions:**
- Run cleanup: `make clean-test-artifacts`
- Check disk usage: `df -h`
- Increase disk space or use RAM disk

### Permission Errors
```
permission denied
```

**Solutions:**
- Check directory permissions: `chmod 755 tmp/`
- Verify ownership: `ls -la tmp/`
- Run with appropriate user permissions

## Best Practices

1. **Regular Cleanup**: Clean temporary files regularly
2. **Isolation**: Use unique names for parallel tests
3. **Retention**: Keep failed test artifacts for debugging
4. **Monitoring**: Monitor disk space usage
5. **Documentation**: Document custom temporary file usage