#!/bin/bash

# FlintRoute Functional Test Runner - Retest
# Re-run tests without cleanup (for debugging)

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

FlintRoute Functional Test Runner - Retest Mode
Re-run tests without cleanup (useful for debugging)

This script will:
  1. Keep existing test environment and artifacts
  2. Re-run tests with --no-cleanup flag
  3. Preserve all logs and results for analysis

Use this when:
  - Debugging test failures
  - Iterating on test fixes
  - Analyzing test artifacts
  - Running tests multiple times quickly

OPTIONS:
    All options from run-tests.sh are supported:
    --pattern PATTERN    Run tests matching pattern
    --config FILE        Use specific config file
    --log-level LEVEL    Set log level (debug|info|warn|error)
    --verbose            Verbose output
    --help               Show this help message

EXAMPLES:
    # Retest all tests
    ./run-retest.sh

    # Retest specific test package
    ./run-retest.sh --pattern ./tests/01_authentication/...

    # Retest with debug logging
    ./run-retest.sh --log-level debug --verbose

EXIT CODES:
    0 - All tests passed
    1 - Tests failed
    2 - Setup/environment error

NOTE:
    This script automatically adds the --no-cleanup flag to preserve
    test artifacts. If you want a clean run, use run-clean.sh instead.

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
    print_info "FlintRoute Retest Mode"
    print_info "====================="
    print_warning "Running tests without cleanup (artifacts will be preserved)"
    echo ""
    
    # Run tests with --no-cleanup flag
    if bash "$SCRIPT_DIR/run-tests.sh" --no-cleanup "$@"; then
        print_success "Retest completed successfully"
        echo ""
        print_info "Test artifacts preserved in:"
        print_info "  - Logs: $SCRIPT_DIR/logs/"
        print_info "  - Results: $SCRIPT_DIR/results/"
        print_info "  - Database: $SCRIPT_DIR/tmp/"
        exit 0
    else
        print_error "Tests failed"
        echo ""
        print_info "Test artifacts preserved for debugging:"
        print_info "  - Logs: $SCRIPT_DIR/logs/"
        print_info "  - Results: $SCRIPT_DIR/results/"
        print_info "  - Database: $SCRIPT_DIR/tmp/"
        exit 1
    fi
}

# Run main function
main "$@"