#!/bin/bash

# FlintRoute Test Environment Setup
# Start all required services for testing

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_ROOT="$(dirname "$(dirname "$TEST_DIR")")"
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

# Function to check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    if ! bash "$SCRIPT_DIR/check-prerequisites.sh"; then
        print_error "Prerequisites check failed"
        return 1
    fi
    
    print_success "Prerequisites check passed"
}

# Function to create necessary directories
create_directories() {
    print_info "Creating necessary directories..."
    
    mkdir -p "$TEST_DIR/logs"
    mkdir -p "$TEST_DIR/results"
    mkdir -p "$TEST_DIR/tmp"
    
    print_success "Directories created"
}

# Function to build mock FRR server
build_mock_frr() {
    print_info "Building mock FRR server..."
    
    if [[ "$VERBOSE" == "true" ]]; then
        bash "$SCRIPT_DIR/build-mock-server.sh"
    else
        bash "$SCRIPT_DIR/build-mock-server.sh" > /dev/null 2>&1
    fi
    
    if [[ $? -ne 0 ]]; then
        print_error "Failed to build mock FRR server"
        return 1
    fi
    
    print_success "Mock FRR server built"
}

# Function to start mock FRR server
start_mock_frr() {
    print_info "Starting mock FRR server..."
    
    local mock_binary="$TEST_DIR/cmd/mock-frr-server/mock-frr-server"
    local mock_config="$TEST_DIR/config/mock-frr-config.yaml"
    local mock_log="$TEST_DIR/logs/mock-frr-server.log"
    
    if [[ ! -f "$mock_binary" ]]; then
        print_error "Mock FRR server binary not found: $mock_binary"
        return 1
    fi
    
    # Start mock FRR server in background
    if [[ "$VERBOSE" == "true" ]]; then
        "$mock_binary" --config "$mock_config" > "$mock_log" 2>&1 &
    else
        "$mock_binary" --config "$mock_config" > "$mock_log" 2>&1 &
    fi
    
    local mock_pid=$!
    echo "$mock_pid" > "$MOCK_FRR_PID_FILE"
    
    print_info "Mock FRR server started (PID: $mock_pid)"
    
    # Wait for mock FRR server to be ready
    print_info "Waiting for mock FRR server to be ready..."
    if ! bash "$SCRIPT_DIR/wait-for-service.sh" grpc localhost 50051 30; then
        print_error "Mock FRR server failed to start"
        cat "$mock_log"
        return 1
    fi
    
    print_success "Mock FRR server is ready"
}

# Function to initialize test database
init_database() {
    print_info "Initializing test database..."
    
    local db_file="$TEST_DIR/tmp/test.db"
    
    # Remove old database if exists
    if [[ -f "$db_file" ]]; then
        rm -f "$db_file"
        print_info "Removed old database"
    fi
    
    # Database will be initialized by FlintRoute on startup
    print_success "Database initialization prepared"
}

# Function to start FlintRoute server
start_flintroute() {
    print_info "Starting FlintRoute server..."
    
    local flintroute_binary="$PROJECT_ROOT/bin/flintroute"
    local flintroute_config="$TEST_DIR/config/test-config.yaml"
    local flintroute_log="$TEST_DIR/logs/flintroute-server.log"
    
    # Build FlintRoute if not exists
    if [[ ! -f "$flintroute_binary" ]]; then
        print_info "Building FlintRoute server..."
        cd "$PROJECT_ROOT"
        if [[ "$VERBOSE" == "true" ]]; then
            make build
        else
            make build > /dev/null 2>&1
        fi
        
        if [[ $? -ne 0 ]]; then
            print_error "Failed to build FlintRoute server"
            cd "$TEST_DIR"
            return 1
        fi
        cd "$TEST_DIR"
    fi
    
    if [[ ! -f "$flintroute_binary" ]]; then
        print_error "FlintRoute binary not found: $flintroute_binary"
        return 1
    fi
    
    # Set environment variables for test mode
    export FLINTROUTE_CONFIG="$flintroute_config"
    export FLINTROUTE_ENV="test"
    export FLINTROUTE_DB_PATH="$TEST_DIR/tmp/test.db"
    export FLINTROUTE_FRR_HOST="localhost"
    export FLINTROUTE_FRR_PORT="50051"
    
    # Start FlintRoute server in background
    if [[ "$VERBOSE" == "true" ]]; then
        "$flintroute_binary" > "$flintroute_log" 2>&1 &
    else
        "$flintroute_binary" > "$flintroute_log" 2>&1 &
    fi
    
    local flintroute_pid=$!
    echo "$flintroute_pid" > "$FLINTROUTE_PID_FILE"
    
    print_info "FlintRoute server started (PID: $flintroute_pid)"
    
    # Wait for FlintRoute server to be ready
    print_info "Waiting for FlintRoute server to be ready..."
    if ! bash "$SCRIPT_DIR/wait-for-service.sh" http localhost 8080 60; then
        print_error "FlintRoute server failed to start"
        cat "$flintroute_log"
        return 1
    fi
    
    print_success "FlintRoute server is ready"
}

# Function to verify all services
verify_services() {
    print_info "Verifying all services..."
    
    local all_ok=true
    
    # Check mock FRR server
    if [[ -f "$MOCK_FRR_PID_FILE" ]]; then
        local mock_pid=$(cat "$MOCK_FRR_PID_FILE")
        if ps -p "$mock_pid" > /dev/null 2>&1; then
            print_success "Mock FRR server is running (PID: $mock_pid)"
        else
            print_error "Mock FRR server is not running"
            all_ok=false
        fi
    else
        print_error "Mock FRR server PID file not found"
        all_ok=false
    fi
    
    # Check FlintRoute server
    if [[ -f "$FLINTROUTE_PID_FILE" ]]; then
        local flintroute_pid=$(cat "$FLINTROUTE_PID_FILE")
        if ps -p "$flintroute_pid" > /dev/null 2>&1; then
            print_success "FlintRoute server is running (PID: $flintroute_pid)"
        else
            print_error "FlintRoute server is not running"
            all_ok=false
        fi
    else
        print_error "FlintRoute server PID file not found"
        all_ok=false
    fi
    
    if [[ "$all_ok" == "true" ]]; then
        print_success "All services verified"
        return 0
    else
        print_error "Service verification failed"
        return 1
    fi
}

# Main execution
main() {
    print_info "FlintRoute Test Environment Setup"
    print_info "=================================="
    
    # Check prerequisites
    if ! check_prerequisites; then
        exit 1
    fi
    
    # Create directories
    create_directories
    
    # Build mock FRR server
    if ! build_mock_frr; then
        exit 1
    fi
    
    # Start mock FRR server
    if ! start_mock_frr; then
        exit 1
    fi
    
    # Initialize database
    if ! init_database; then
        exit 1
    fi
    
    # Start FlintRoute server
    if ! start_flintroute; then
        # Cleanup on failure
        bash "$SCRIPT_DIR/teardown-env.sh" 2>/dev/null || true
        exit 1
    fi
    
    # Verify all services
    if ! verify_services; then
        # Cleanup on failure
        bash "$SCRIPT_DIR/teardown-env.sh" 2>/dev/null || true
        exit 1
    fi
    
    echo ""
    print_success "Test environment setup completed successfully"
    print_info "Services:"
    print_info "  - Mock FRR: localhost:50051 (gRPC)"
    print_info "  - FlintRoute: localhost:8080 (HTTP)"
    print_info "Logs:"
    print_info "  - Mock FRR: $TEST_DIR/logs/mock-frr-server.log"
    print_info "  - FlintRoute: $TEST_DIR/logs/flintroute-server.log"
    
    exit 0
}

# Run main function
main "$@"