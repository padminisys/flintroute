# Test Results

This directory contains test execution results in multiple formats for different use cases.

## Result Formats

### JSON Format
Machine-readable results for programmatic analysis:

**test-results.json**
```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "duration_seconds": 125.5,
  "total_tests": 45,
  "passed": 43,
  "failed": 2,
  "skipped": 0,
  "suites": [
    {
      "name": "01_authentication",
      "tests": 8,
      "passed": 8,
      "failed": 0,
      "duration_seconds": 12.3
    }
  ],
  "failures": [
    {
      "suite": "02_peer_management",
      "test": "TestCreatePeerWithInvalidIP",
      "error": "expected status 400, got 500",
      "stack_trace": "..."
    }
  ]
}
```

### XML Format (JUnit)
Standard format for CI/CD integration:

**test-results.xml**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites tests="45" failures="2" time="125.5">
  <testsuite name="01_authentication" tests="8" failures="0" time="12.3">
    <testcase name="TestUserLogin" time="1.2"/>
    <testcase name="TestTokenRefresh" time="0.8"/>
  </testsuite>
  <testsuite name="02_peer_management" tests="10" failures="1" time="15.7">
    <testcase name="TestCreatePeer" time="1.5"/>
    <testcase name="TestCreatePeerWithInvalidIP" time="0.9">
      <failure message="expected status 400, got 500">
        Stack trace...
      </failure>
    </testcase>
  </testsuite>
</testsuites>
```

### HTML Format
Human-readable reports with visual indicators:

**test-report.html**
- Summary dashboard with pass/fail statistics
- Detailed test results with expandable sections
- Error messages and stack traces
- Execution timeline
- Coverage information (if available)

## File Naming Convention

Results are timestamped for historical tracking:

```
test-results-2024-01-15T10-30-45Z.json
test-results-2024-01-15T10-30-45Z.xml
test-report-2024-01-15T10-30-45Z.html
```

Latest results are also available without timestamp:
```
test-results.json
test-results.xml
test-report.html
```

## Result Structure

### Test Suite Results
Each test suite generates detailed results:

```json
{
  "suite": "02_peer_management",
  "tests": [
    {
      "name": "TestCreatePeer",
      "status": "passed",
      "duration_ms": 1500,
      "assertions": 5
    },
    {
      "name": "TestUpdatePeer",
      "status": "failed",
      "duration_ms": 900,
      "error": "assertion failed: expected 200, got 500",
      "stack_trace": "...",
      "logs": ["..."]
    }
  ]
}
```

### Coverage Results
Code coverage data (if enabled):

**coverage.json**
```json
{
  "total_coverage": 78.5,
  "packages": [
    {
      "name": "internal/api",
      "coverage": 85.2,
      "files": [
        {
          "name": "auth_handlers.go",
          "coverage": 92.3,
          "lines_covered": 120,
          "lines_total": 130
        }
      ]
    }
  ]
}
```

## Viewing Results

### Command Line
```bash
# View summary
jq '.total_tests, .passed, .failed' test-results.json

# List failures
jq '.failures[] | {suite, test, error}' test-results.json

# View specific suite
jq '.suites[] | select(.name=="02_peer_management")' test-results.json
```

### Web Browser
```bash
# Open HTML report
open test-report.html  # macOS
xdg-open test-report.html  # Linux
start test-report.html  # Windows
```

### CI Dashboard
Results are automatically uploaded to CI systems:
- GitHub Actions: Test summary in workflow run
- GitLab CI: Test report tab
- Jenkins: Test results trend

## Result Analysis

### Trend Analysis
Compare results over time:

```bash
# Compare pass rates
for file in test-results-*.json; do
  echo "$file: $(jq -r '.passed / .total_tests * 100' $file)%"
done

# Find flaky tests (intermittent failures)
grep -h '"status":"failed"' test-results-*.json | \
  jq -r '.test' | sort | uniq -c | sort -rn
```

### Performance Tracking
Monitor test execution time:

```bash
# Suite duration trends
jq '.suites[] | {name, duration: .duration_seconds}' test-results.json

# Slowest tests
jq '.suites[].tests[] | {name, duration: .duration_ms}' test-results.json | \
  jq -s 'sort_by(.duration) | reverse | .[0:10]'
```

## CI Integration

### GitHub Actions
```yaml
- name: Run Tests
  run: make test-functional

- name: Publish Test Results
  uses: EnricoMi/publish-unit-test-result-action@v2
  if: always()
  with:
    files: test/functional/results/test-results.xml

- name: Upload Results
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: test-results
    path: test/functional/results/
```

### GitLab CI
```yaml
test:
  script:
    - make test-functional
  artifacts:
    when: always
    reports:
      junit: test/functional/results/test-results.xml
    paths:
      - test/functional/results/
```

## Retention Policy

Results are retained according to:
- **Local**: Last 10 runs (configurable)
- **CI**: According to CI retention policy (typically 30-90 days)
- **Archive**: Long-term storage for release testing

## Cleanup

Remove old results:

```bash
# Keep only latest 10 results
ls -t test-results-*.json | tail -n +11 | xargs rm -f
ls -t test-results-*.xml | tail -n +11 | xargs rm -f
ls -t test-report-*.html | tail -n +11 | xargs rm -f

# Remove all results
rm -f test-results*.* test-report*.* coverage*.*
```

## Troubleshooting

### Missing Results
- Verify tests completed successfully
- Check file permissions (directory must be writable)
- Ensure result generation is enabled in test configuration

### Incomplete Results
- Check for test timeouts
- Review logs for errors during result generation
- Verify all test suites completed

### CI Upload Failures
- Verify artifact path is correct
- Check file size limits (some CI systems have limits)
- Ensure proper permissions for CI user