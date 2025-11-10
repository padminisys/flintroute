#!/bin/bash

# FlintRoute Test Complete Cleanup
# Full cleanup of all test artifacts

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"

# Default: interactive mode
FORCE=false

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

Complete cleanup of all test artifacts

This script performs a full cleanup by:
  1. Stopping all services (teardown-env.sh)
  2. Cleaning database (cleanup-db.sh --full)
  3. Cleaning results (cleanup-results.sh)
  4. Cleaning logs (cleanup-logs.sh)
  5. Cleaning temporary files (tmp/*)
  6. Verifying cleanup completed

OPTIONS:
    --force    Skip confirmation prompt
    --help     Show this help message

EXAMPLES:
    # Interactive cleanup (with confirmation)
    ./cleanup-all.sh

    # Force cleanup (no confirmation)
    ./cleanup-all.sh --force

EXIT CODES:
    0 - Cleanup successful
    1 - Cleanup failed

WARNING:
    This will remove ALL test artifacts including:
    - Test database
    - Test results
    - Log files
    - Temporary files
    - Running services

EOF
}

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --force)
                FORCE=true
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

# Function to confirm cleanup
confirm_cleanup() {
    if [[ "$FORCE" == "true" ]]; then
        return 0
    fi
    
    echo ""
    print_warning "This will remove ALL test artifacts:"
    print_warning "  - Stop all running services"
    print_warning "  - Delete test database"
    print_warning "  - Delete test results"
    print_warning "  - Delete log files"
    print_warning "  - Delete temporary files"
    echo ""
    
    read -p "Are you sure you want to continue? (yes/no): " -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        print_info "Cleanup cancelled"
        exit 0
    fi
}

# Function to stop services
stop_services() {
    print_info "Step 1/6: Stopping all services..."
    
    if bash "$SCRIPT_DIR/teardown-env.sh"; then
        print_success "Services stopped"
        return 0
    else
        print_warning "Service teardown had issues (continuing anyway)"
        return 0
    fi
}

# Function to cleanup database
cleanup_database() {
    print_info "Step 2/6: Cleaning database..."
    
    if bash "$SCRIPT_DIR/cleanup-db.sh" --full; then
        print_success "Database cleaned"
        return 0
    else
        print_error "Database cleanup failed"
        return 1
    fi
}

# Function to cleanup results
cleanup_results() {
    print_info "Step 3/6: Cleaning test results..."
    
    if bash "$SCRIPT_DIR/cleanup-results.sh"; then
        print_success "Test results cleaned"
        return 0
    else
        print_error "Results cleanup failed"
        return 1
    fi
}

# Function to cleanup logs
cleanup_logs() {
    print_info "Step 4/6: Cleaning log files..."
    
    if bash "$SCRIPT_DIR/cleanup-logs.sh"; then
        print_success "Log files cleaned"
        return 0
    else
        print_error "Logs cleanup failed"
        return 1
    fi
}

# Function to cleanup temporary files
cleanup_temp() {
    print_info "Step 5/6: Cleaning temporary files..."
    
    local tmp_dir="$TEST_DIR/tmp"
    
    if [[ ! -d "$tmp_dir" ]]; then
        print_info "Temporary directory does not exist"
        return 0
    fi
    
    # Count files before cleanup
    local file_count=0
    if [[ -n "$(ls -A "$tmp_dir" 2>/dev/null)" ]]; then
        file_count=$(find "$tmp_dir" -type f ! -name "README.md" | wc -l)
    fi
    
    if [[ $file_count -eq 0 ]]; then
        print_info "No temporary files to remove"
        return 0
    fi
    
    print_info "Found $file_count temporary file(s)"
    
    # Remove all files except README.md
    find "$tmp_dir" -type f ! -name "README.md" -delete
    
    # Remove empty directories
    find "$tmp_dir" -type d -empty -delete 2>/dev/null || true
    
    # Recreate tmp directory if it was removed
    mkdir -p "$tmp_dir"
    
    # Verify README.md exists
    if [[ ! -f "$tmp_dir/README.md" ]]; then
        cat > "$tmp_dir/README.md" << 'EOF'
# Temporary Files Directory

This directory contains temporary files generated during test runs:

- `test.db` - Test database file
- `*.pid` - Process ID files for running services
- Other temporary test artifacts

Files in this directory are automatically cleaned up after test runs.
EOF
        print_info "Recreated README.md file"
    fi
    
    print_success "Temporary files cleaned"
}

# Function to verify cleanup
verify_cleanup() {
    print_info "Step 6/6: Verifying cleanup..."
    
    local issues=0
    
    # Check for running services
    if pgrep -f "mock-frr-server" > /dev/null 2>&1; then
        print_warning "Mock FRR server is still running"
        ((issues++))
    fi
    
    if lsof -ti:8080 > /dev/null 2>&1; then
        print_warning "Port 8080 is still in use"
        ((issues++))
    fi
    
    if lsof -ti:50051 > /dev/null 2>&1; then
        print_warning "Port 50051 is still in use"
        ((issues++))
    fi
    
    # Check for database file
    if [[ -f "$TEST_DIR/tmp/test.db" ]]; then
        print_warning "Database file still exists"
        ((issues++))
    fi
    
    # Check for PID files
    if [[ -f "$TEST_DIR/tmp/mock-frr.pid" ]] || [[ -f "$TEST_DIR/tmp/flintroute.pid" ]]; then
        print_warning "PID files still exist"
        ((issues++))
    fi
    
    # Check for test results
    if compgen -G "$TEST_DIR/results/*.json" > /dev/null 2>&1; then
        print_warning "Test result files still exist"
        ((issues++))
    fi
    
    # Check for log files
    if compgen -G "$TEST_DIR/logs/*.log" > /dev/null 2>&1; then
        print_warning "Log files still exist"
        ((issues++))
    fi
    
    if [[ $issues -eq 0 ]]; then
        print_success "Cleanup verification passed"
        return 0
    else
        print_warning "Cleanup verification found $issues issue(s)"
        return 1
    fi
}

# Main execution
main() {
    print_info "FlintRoute Complete Test Cleanup"
    print_info "================================="
    
    # Parse arguments
    parse_args "$@"
    
    # Confirm cleanup
    confirm_cleanup
    
    local exit_code=0
    
    # Step 1: Stop services
    if ! stop_services; then
        exit_code=1
    fi
    
    echo ""
    
    # Step 2: Cleanup database
    if ! cleanup_database; then
        exit_code=1
    fi
    
    echo ""
    
    # Step 3: Cleanup results
    if ! cleanup_results; then
        exit_code=1
    fi
    
    echo ""
    
    # Step 4: Cleanup logs
    if ! cleanup_logs; then
        exit_code=1
    fi
    
    echo ""
    
    # Step 5: Cleanup temporary files
    if ! cleanup_temp; then
        exit_code=1
    fi
    
    echo ""
    
    # Step 6: Verify cleanup
    if ! verify_cleanup; then
        exit_code=1
    fi
    
    echo ""
    echo "========================================"
    
    if [[ $exit_code -eq 0 ]]; then
        print_success "Complete cleanup finished successfully"
        print_info "Test environment is now clean and ready for new tests"
    else
        print_warning "Complete cleanup finished with warnings"
        print_info "Some cleanup steps had issues, but environment should be usable"
    fi
    
    exit $exit_code
}

# Run main function
main "$@"