#!/bin/bash

# Wait for a service to become ready
# Supports HTTP, gRPC, and TCP checks

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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
Usage: $(basename "$0") <type> <host> <port> [timeout]

Wait for a service to become ready with retry logic

ARGUMENTS:
    type      Service type: http, grpc, tcp
    host      Service hostname or IP
    port      Service port number
    timeout   Timeout in seconds (default: 30)

TYPES:
    http   - HTTP health check (GET request to /health or /)
    grpc   - gRPC port check (TCP connection)
    tcp    - TCP port check

EXAMPLES:
    # Wait for HTTP service
    ./wait-for-service.sh http localhost 8080

    # Wait for gRPC service with 60s timeout
    ./wait-for-service.sh grpc localhost 50051 60

    # Wait for TCP service
    ./wait-for-service.sh tcp localhost 5432 30

EXIT CODES:
    0 - Service is ready
    1 - Service failed to become ready within timeout

EOF
}

# Check arguments
if [[ $# -lt 3 ]]; then
    show_help
    exit 1
fi

TYPE="$1"
HOST="$2"
PORT="$3"
TIMEOUT="${4:-30}"

# Validate type
if [[ ! "$TYPE" =~ ^(http|grpc|tcp)$ ]]; then
    print_error "Invalid type: $TYPE (must be http, grpc, or tcp)"
    show_help
    exit 1
fi

# Function to check TCP connection
check_tcp() {
    local host=$1
    local port=$2
    
    if command -v nc > /dev/null 2>&1; then
        # Use netcat if available
        nc -z -w1 "$host" "$port" > /dev/null 2>&1
    elif command -v timeout > /dev/null 2>&1; then
        # Use timeout with bash TCP
        timeout 1 bash -c "cat < /dev/null > /dev/tcp/$host/$port" > /dev/null 2>&1
    else
        # Fallback to bash TCP
        (echo > /dev/tcp/"$host"/"$port") > /dev/null 2>&1
    fi
}

# Function to check HTTP service
check_http() {
    local host=$1
    local port=$2
    
    # Try /health endpoint first, then root
    if command -v curl > /dev/null 2>&1; then
        curl -sf "http://${host}:${port}/health" > /dev/null 2>&1 || \
        curl -sf "http://${host}:${port}/" > /dev/null 2>&1
    elif command -v wget > /dev/null 2>&1; then
        wget -q -O /dev/null "http://${host}:${port}/health" 2>&1 || \
        wget -q -O /dev/null "http://${host}:${port}/" 2>&1
    else
        # Fallback to TCP check
        check_tcp "$host" "$port"
    fi
}

# Function to check gRPC service
check_grpc() {
    local host=$1
    local port=$2
    
    # gRPC uses TCP, so we check if the port is open
    # For more sophisticated checks, grpcurl could be used
    check_tcp "$host" "$port"
}

# Function to wait for service
wait_for_service() {
    local type=$1
    local host=$2
    local port=$3
    local timeout=$4
    
    print_info "Waiting for $type service at ${host}:${port} (timeout: ${timeout}s)..."
    
    local start_time=$(date +%s)
    local attempt=0
    local backoff=1
    
    while true; do
        ((attempt++))
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        if [[ $elapsed -ge $timeout ]]; then
            print_error "Timeout waiting for service (${elapsed}s elapsed)"
            return 1
        fi
        
        # Check service based on type
        local check_result=false
        case "$type" in
            http)
                if check_http "$host" "$port"; then
                    check_result=true
                fi
                ;;
            grpc)
                if check_grpc "$host" "$port"; then
                    check_result=true
                fi
                ;;
            tcp)
                if check_tcp "$host" "$port"; then
                    check_result=true
                fi
                ;;
        esac
        
        if [[ "$check_result" == "true" ]]; then
            print_success "Service is ready (attempt $attempt, ${elapsed}s elapsed)"
            return 0
        fi
        
        # Exponential backoff (max 5 seconds)
        if [[ $backoff -lt 5 ]]; then
            backoff=$((backoff * 2))
            if [[ $backoff -gt 5 ]]; then
                backoff=5
            fi
        fi
        
        # Show progress every 5 attempts
        if [[ $((attempt % 5)) -eq 0 ]]; then
            print_info "Still waiting... (attempt $attempt, ${elapsed}s elapsed)"
        fi
        
        sleep "$backoff"
    done
}

# Main execution
main() {
    if ! wait_for_service "$TYPE" "$HOST" "$PORT" "$TIMEOUT"; then
        print_error "Service at ${HOST}:${PORT} is not ready"
        exit 1
    fi
    exit 0
}

# Run main function
main