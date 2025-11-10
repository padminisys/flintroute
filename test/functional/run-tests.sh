#!/bin/bash

# FlintRoute Functional Test Runner
# Main test execution script with full lifecycle management

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
PATTERN="./..."
CONFIG_FILE="config/test-config.yaml"
LOG_LEVEL="info"
NO_CLEANUP=false
VERBOSE=false
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_FAILED=false

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_help() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS]

FlintRoute Functional Test Runner

OPTIONS:
    --pattern PATTERN    Run tests matching pattern (default: ./...)
    --config FILE        Use specific config file (default: config/test-config.yaml)
    --log-level LEVEL    Set log level: debug|info|warn|error (default: info)
    --no-cleanup         Don't cleanup on success
    --verbose            Verbose output
    --help               Show this help message

EXAMPLES:
    # Run all tests
    ./run-tests.sh

    # Run specific test package
    ./run-tests.sh --pattern ./tests/01_authentication/...

    # Run with debug logging
    ./run-tests.sh --log-level debug --verbose

    # Run without cleanup (for debugging)
    ./run-tests.sh --no-cleanup

EXIT CODES:
    0 - All tests passed
    1 - Tests failed
    2 - Setup/environment error

EOF
}

# Parse command-line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --pattern)
                PATTERN="$2"
                shift 2
                ;;
            --config)
                CONFIG_FILE="$2"
                shift 2
                ;;
            --log-level)
                LOG_LEVEL="$2"
                shift 2
                ;;
            --no-cleanup)
                NO_CLEANUP=true
                shift
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# Function to check if config file exists
check_config() {
    if [[ ! -f "$SCRIPT_DIR/$CONFIG_FILE" ]]; then
        print_error "Config file not found: $CONFIG_FILE"
        return 1
    fi
    print_info "Using config file: $CONFIG_FILE"
}

# Function to setup environment
setup_environment() {
    print_info "Setting up test environment..."
    
    if [[ "$VERBOSE" == "true" ]]; then
        bash "$SCRIPT_DIR/scripts/setup-env.sh" --verbose
    else
        bash "$SCRIPT_DIR/scripts/setup-env.sh"
    fi
    
    if [[ $? -ne 0 ]]; then
        print_error "Environment setup failed"
        return 2
    fi
    
    print_success "Environment setup completed"
}

# Function to run tests
run_tests() {
    print_info "Running tests with pattern: $PATTERN"
    print_info "Log level: $LOG_LEVEL"
    
    # Set environment variables
    export TEST_CONFIG="$SCRIPT_DIR/$CONFIG_FILE"
    export TEST_LOG_LEVEL="$LOG_LEVEL"
    export TEST_RESULTS_DIR="$SCRIPT_DIR/results"
    
    # Create results directory if it doesn't exist
    mkdir -p "$TEST_RESULTS_DIR"
    
    # Prepare test command
    local test_cmd="go test"
    test_cmd="$test_cmd -v"
    test_cmd="$test_cmd -timeout 30m"
    test_cmd="$test_cmd -count=1"
    
    if [[ "$VERBOSE" == "true" ]]; then
        test_cmd="$test_cmd -v"
    fi
    
    # Add JSON output for report generation
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local json_output="$TEST_RESULTS_DIR/test-results-${timestamp}.json"
    test_cmd="$test_cmd -json"
    
    print_info "Executing: $test_cmd $PATTERN"
    print_info "Results will be saved to: $json_output"
    
    # Run tests and capture output
    cd "$SCRIPT_DIR"
    if eval "$test_cmd $PATTERN" | tee "$json_output"; then
        print_success "All tests passed"
        return 0
    else
        print_error "Some tests failed"
        TEST_FAILED=true
        return 1
    fi
}

# Function to generate reports
generate_reports() {
    print_info "Generating test reports..."
    
    local latest_result=$(ls -t "$SCRIPT_DIR/results"/test-results-*.json 2>/dev/null | head -1)
    
    if [[ -z "$latest_result" ]]; then
        print_warning "No test results found to generate reports"
        return 0
    fi
    
    # Generate HTML report (if tool available)
    if command -v go-test-report &> /dev/null; then
        local html_report="${latest_result%.json}.html"
        cat "$latest_result" | go-test-report > "$html_report" 2>/dev/null || true
        if [[ -f "$html_report" ]]; then
            print_success "HTML report generated: $html_report"
        fi
    fi
    
    # Generate summary
    local summary_file="$SCRIPT_DIR/results/test-summary.txt"
    {
        echo "FlintRoute Functional Test Summary"
        echo "=================================="
        echo "Date: $(date)"
        echo "Pattern: $PATTERN"
        echo "Config: $CONFIG_FILE"
        echo ""
        
        # Count pass/fail from JSON
        local total=$(grep -c '"Action":"pass"\|"Action":"fail"' "$latest_result" 2>/dev/null || echo "0")
        local passed=$(grep -c '"Action":"pass"' "$latest_result" 2>/dev/null || echo "0")
        local failed=$(grep -c '"Action":"fail"' "$latest_result" 2>/dev/null || echo "0")
        
        echo "Total Tests: $total"
        echo "Passed: $passed"
        echo "Failed: $failed"
        
        if [[ $failed -gt 0 ]]; then
            echo ""
            echo "Failed Tests:"
            grep '"Action":"fail"' "$latest_result" | grep -o '"Test":"[^"]*"' | cut -d'"' -f4 || true
        fi
    } > "$summary_file"
    
    print_success "Test summary generated: $summary_file"
    cat "$summary_file"
}

# Function to teardown environment
teardown_environment() {
    print_info "Tearing down test environment..."
    
    if [[ "$VERBOSE" == "true" ]]; then
        bash "$SCRIPT_DIR/scripts/teardown-env.sh" --verbose
    else
        bash "$SCRIPT_DIR/scripts/teardown-env.sh"
    fi
    
    if [[ $? -ne 0 ]]; then
        print_warning "Environment teardown had issues (non-critical)"
    else
        print_success "Environment teardown completed"
    fi
}

# Function to cleanup on success
cleanup_on_success() {
    if [[ "$NO_CLEANUP" == "true" ]]; then
        print_info "Skipping cleanup (--no-cleanup flag set)"
        return 0
    fi
    
    if [[ "$TEST_FAILED" == "true" ]]; then
        print_info "Skipping cleanup (tests failed, keeping artifacts for debugging)"
        return 0
    fi
    
    print_info "Running cleanup..."
    bash "$SCRIPT_DIR/scripts/cleanup-logs.sh" --keep-latest
    print_success "Cleanup completed"
}

# Main execution
main() {
    print_info "FlintRoute Functional Test Runner"
    print_info "=================================="
    
    # Parse arguments
    parse_args "$@"
    
    # Check configuration
    if ! check_config; then
        exit 2
    fi
    
    # Setup trap for cleanup on exit
    trap 'teardown_environment' EXIT
    
    # Setup environment
    if ! setup_environment; then
        print_error "Failed to setup environment"
        exit 2
    fi
    
    # Run tests
    local test_exit_code=0
    if ! run_tests; then
        test_exit_code=1
    fi
    
    # Generate reports
    generate_reports
    
    # Cleanup if successful
    cleanup_on_success
    
    # Final status
    echo ""
    if [[ $test_exit_code -eq 0 ]]; then
        print_success "Test run completed successfully"
    else
        print_error "Test run completed with failures"
    fi
    
    exit $test_exit_code
}

# Run main function
main "$@"