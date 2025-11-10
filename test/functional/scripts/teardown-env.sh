#!/bin/bash

# FlintRoute Test Environment Teardown
# Stop all services gracefully

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
VERBOSE=false

# PID files
MOCK_FRR_PID_FILE="$TEST_DIR/tmp/mock-frr.pid"
FLINTROUTE_PID_FILE="$TEST_DIR/tmp/flintroute.pid"

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

# Parse arguments
for arg in "$@"; do
    if [[ "$arg" == "--verbose" ]]; then
        VERBOSE=true
    fi
done

# Function to stop a process gracefully
stop_process() {
    local pid=$1
    local name=$2
    local timeout=${3:-10}
    
    if ! ps -p "$pid" > /dev/null 2>&1; then
        print_warning "$name (PID: $pid) is not running"
        return 0
    fi
    
    print_info "Stopping $name (PID: $pid)..."
    
    # Send SIGTERM
    kill -TERM "$pid" 2>/dev/null || true
    
    # Wait for process to exit
    local count=0
    while ps -p "$pid" > /dev/null 2>&1; do
        if [[ $count -ge $timeout ]]; then
            print_warning "$name did not stop gracefully, sending SIGKILL..."
            kill -KILL "$pid" 2>/dev/null || true
            sleep 1
            break
        fi
        sleep 1
        ((count++))
    done
    
    if ps -p "$pid" > /dev/null 2>&1; then
        print_error "Failed to stop $name"
        return 1
    else
        print_success "$name stopped"
        return 0
    fi
}

# Function to stop FlintRoute server
stop_flintroute() {
    if [[ ! -f "$FLINTROUTE_PID_FILE" ]]; then
        print_warning "FlintRoute PID file not found"
        return 0
    fi
    
    local pid=$(cat "$FLINTROUTE_PID_FILE")
    stop_process "$pid" "FlintRoute server" 15
    
    # Remove PID file
    rm -f "$FLINTROUTE_PID_FILE"
}

# Function to stop mock FRR server
stop_mock_frr() {
    if [[ ! -f "$MOCK_FRR_PID_FILE" ]]; then
        print_warning "Mock FRR PID file not found"
        return 0
    fi
    
    local pid=$(cat "$MOCK_FRR_PID_FILE")
    stop_process "$pid" "Mock FRR server" 10
    
    # Remove PID file
    rm -f "$MOCK_FRR_PID_FILE"
}

# Function to kill orphaned processes
kill_orphaned_processes() {
    print_info "Checking for orphaned processes..."
    
    local killed=false
    
    # Kill any mock-frr-server processes
    if pgrep -f "mock-frr-server" > /dev/null 2>&1; then
        print_warning "Found orphaned mock-frr-server processes"
        pkill -TERM -f "mock-frr-server" 2>/dev/null || true
        sleep 2
        pkill -KILL -f "mock-frr-server" 2>/dev/null || true
        killed=true
    fi
    
    # Kill any FlintRoute test processes on port 8080
    if lsof -ti:8080 > /dev/null 2>&1; then
        print_warning "Found process using port 8080"
        local port_pid=$(lsof -ti:8080)
        kill -TERM "$port_pid" 2>/dev/null || true
        sleep 2
        kill -KILL "$port_pid" 2>/dev/null || true
        killed=true
    fi
    
    # Kill any processes on gRPC port 50051
    if lsof -ti:50051 > /dev/null 2>&1; then
        print_warning "Found process using port 50051"
        local port_pid=$(lsof -ti:50051)
        kill -TERM "$port_pid" 2>/dev/null || true
        sleep 2
        kill -KILL "$port_pid" 2>/dev/null || true
        killed=true
    fi
    
    if [[ "$killed" == "true" ]]; then
        print_success "Orphaned processes cleaned up"
    else
        print_info "No orphaned processes found"
    fi
}

# Function to close database connections
close_database_connections() {
    print_info "Closing database connections..."
    
    # SQLite doesn't require explicit connection closing
    # Just ensure no processes are holding the database file
    local db_file="$TEST_DIR/tmp/test.db"
    
    if [[ -f "$db_file" ]]; then
        # Check if any process has the database file open
        if lsof "$db_file" > /dev/null 2>&1; then
            print_warning "Database file is still open by some process"
            lsof "$db_file" || true
        else
            print_success "No active database connections"
        fi
    fi
}

# Function to verify teardown
verify_teardown() {
    print_info "Verifying teardown..."
    
    local all_ok=true
    
    # Check if mock FRR is stopped
    if pgrep -f "mock-frr-server" > /dev/null 2>&1; then
        print_error "Mock FRR server is still running"
        all_ok=false
    fi
    
    # Check if FlintRoute is stopped
    if lsof -ti:8080 > /dev/null 2>&1; then
        print_error "Port 8080 is still in use"
        all_ok=false
    fi
    
    # Check if gRPC port is free
    if lsof -ti:50051 > /dev/null 2>&1; then
        print_error "Port 50051 is still in use"
        all_ok=false
    fi
    
    # Check PID files are removed
    if [[ -f "$MOCK_FRR_PID_FILE" ]] || [[ -f "$FLINTROUTE_PID_FILE" ]]; then
        print_warning "PID files still exist"
        rm -f "$MOCK_FRR_PID_FILE" "$FLINTROUTE_PID_FILE"
    fi
    
    if [[ "$all_ok" == "true" ]]; then
        print_success "Teardown verified"
        return 0
    else
        print_warning "Teardown verification found issues"
        return 1
    fi
}

# Main execution
main() {
    print_info "FlintRoute Test Environment Teardown"
    print_info "====================================="
    
    local exit_code=0
    
    # Stop FlintRoute server
    if ! stop_flintroute; then
        exit_code=1
    fi
    
    # Stop mock FRR server
    if ! stop_mock_frr; then
        exit_code=1
    fi
    
    # Close database connections
    close_database_connections
    
    # Kill any orphaned processes
    kill_orphaned_processes
    
    # Verify teardown
    if ! verify_teardown; then
        exit_code=1
    fi
    
    echo ""
    if [[ $exit_code -eq 0 ]]; then
        print_success "Test environment teardown completed successfully"
    else
        print_warning "Test environment teardown completed with warnings"
    fi
    
    exit $exit_code
}

# Run main function
main "$@"