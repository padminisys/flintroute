#!/bin/bash

# FlintRoute Test Logs Cleanup
# Clean log files with option to keep latest

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
LOGS_DIR="$TEST_DIR/logs"

# Default: remove all logs
KEEP_LATEST=false

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

Clean log files from test runs

This script removes log files from the logs directory.
By default, all log files are removed. Use --keep-latest
to preserve the most recent log file of each type.

OPTIONS:
    --keep-latest    Keep the most recent log file of each type
    --help           Show this help message

LOG FILES:
    - mock-frr-server.log
    - flintroute-server.log
    - test-runner.log
    - Any other *.log files

EXAMPLES:
    # Remove all log files
    ./cleanup-logs.sh

    # Keep the latest log file of each type
    ./cleanup-logs.sh --keep-latest

EXIT CODES:
    0 - Cleanup successful
    1 - Cleanup failed

EOF
}

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --keep-latest)
                KEEP_LATEST=true
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

# Function to get latest file by modification time
get_latest_file() {
    local pattern=$1
    
    if compgen -G "$LOGS_DIR/$pattern" > /dev/null 2>&1; then
        ls -t "$LOGS_DIR"/$pattern 2>/dev/null | head -1
    fi
}

# Function to cleanup logs
cleanup_logs() {
    print_info "Cleaning logs directory: $LOGS_DIR"
    
    if [[ ! -d "$LOGS_DIR" ]]; then
        print_warning "Logs directory does not exist"
        return 0
    fi
    
    # Count log files
    local log_count=0
    if compgen -G "$LOGS_DIR/*.log" > /dev/null 2>&1; then
        log_count=$(ls -1 "$LOGS_DIR"/*.log 2>/dev/null | wc -l)
    fi
    
    if [[ $log_count -eq 0 ]]; then
        print_info "No log files to remove"
        return 0
    fi
    
    print_info "Found $log_count log file(s)"
    
    if [[ "$KEEP_LATEST" == "true" ]]; then
        print_info "Keeping latest log files..."
        
        # Get unique log file base names
        local log_types=()
        for log_file in "$LOGS_DIR"/*.log; do
            if [[ -f "$log_file" ]]; then
                local basename=$(basename "$log_file")
                # Extract base name without timestamp if present
                local base_type=$(echo "$basename" | sed 's/-[0-9]\{8\}_[0-9]\{6\}\.log$/\.log/' | sed 's/\.log$//')
                if [[ ! " ${log_types[@]} " =~ " ${base_type} " ]]; then
                    log_types+=("$base_type")
                fi
            fi
        done
        
        # For each log type, keep the latest and remove others
        local removed=0
        local kept=0
        
        for log_type in "${log_types[@]}"; do
            local latest=$(get_latest_file "${log_type}*.log")
            
            if [[ -n "$latest" ]]; then
                print_info "Keeping latest: $(basename "$latest")"
                ((kept++))
                
                # Remove all other files of this type
                for log_file in "$LOGS_DIR"/${log_type}*.log; do
                    if [[ -f "$log_file" ]] && [[ "$log_file" != "$latest" ]]; then
                        rm -f "$log_file"
                        print_info "Removed: $(basename "$log_file")"
                        ((removed++))
                    fi
                done
            fi
        done
        
        print_success "Kept $kept latest log file(s), removed $removed old log file(s)"
    else
        # Remove all log files
        local removed=0
        for log_file in "$LOGS_DIR"/*.log; do
            if [[ -f "$log_file" ]]; then
                rm -f "$log_file"
                print_info "Removed: $(basename "$log_file")"
                ((removed++))
            fi
        done
        
        print_success "Removed $removed log file(s)"
    fi
    
    # Verify README.md is still there
    if [[ ! -f "$LOGS_DIR/README.md" ]]; then
        cat > "$LOGS_DIR/README.md" << 'EOF'
# Test Logs Directory

This directory contains log files from test runs:

- `mock-frr-server.log` - Mock FRR server logs
- `flintroute-server.log` - FlintRoute server logs
- `test-runner.log` - Test execution logs

Logs are automatically cleaned up after successful test runs.
Use `--no-cleanup` flag to preserve logs for debugging.
EOF
        print_info "Recreated README.md file"
    fi
}

# Main execution
main() {
    print_info "FlintRoute Test Logs Cleanup"
    print_info "============================="
    
    # Parse arguments
    parse_args "$@"
    
    if [[ "$KEEP_LATEST" == "true" ]]; then
        print_info "Mode: Keep latest log files"
    else
        print_info "Mode: Remove all log files"
    fi
    
    echo ""
    
    if cleanup_logs; then
        echo ""
        print_success "Logs cleanup completed successfully"
        exit 0
    else
        echo ""
        print_error "Logs cleanup failed"
        exit 1
    fi
}

# Run main function
main "$@"