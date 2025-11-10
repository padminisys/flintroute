#!/bin/bash

# FlintRoute Functional Test Runner - Clean Run
# Full cleanup followed by test execution

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_help() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS]

FlintRoute Functional Test Runner - Clean Run
Performs full cleanup before running tests

This script will:
  1. Stop all running services
  2. Clean all test artifacts (database, logs, results)
  3. Run the full test suite

OPTIONS:
    All options from run-tests.sh are supported:
    --pattern PATTERN    Run tests matching pattern
    --config FILE        Use specific config file
    --log-level LEVEL    Set log level (debug|info|warn|error)
    --verbose            Verbose output
    --help               Show this help message

EXAMPLES:
    # Clean run with all tests
    ./run-clean.sh

    # Clean run with specific pattern
    ./run-clean.sh --pattern ./tests/01_authentication/...

    # Clean run with verbose output
    ./run-clean.sh --verbose

EXIT CODES:
    0 - All tests passed
    1 - Tests failed
    2 - Setup/cleanup error

EOF
}

# Check for help flag
for arg in "$@"; do
    if [[ "$arg" == "--help" ]]; then
        show_help
        exit 0
    fi
done

# Main execution
main() {
    print_info "FlintRoute Clean Test Run"
    print_info "========================="
    echo ""
    
    # Step 1: Full cleanup
    print_info "Step 1: Running full cleanup..."
    if bash "$SCRIPT_DIR/scripts/cleanup-all.sh"; then
        print_success "Cleanup completed"
    else
        print_error "Cleanup failed"
        exit 2
    fi
    
    echo ""
    
    # Step 2: Run tests
    print_info "Step 2: Running tests..."
    if bash "$SCRIPT_DIR/run-tests.sh" "$@"; then
        print_success "Clean test run completed successfully"
        exit 0
    else
        print_error "Tests failed"
        exit 1
    fi
}

# Run main function
main "$@"