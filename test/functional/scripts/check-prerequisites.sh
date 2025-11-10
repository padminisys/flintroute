#!/bin/bash

# FlintRoute Test Prerequisites Checker
# Check if all prerequisites are met for running tests

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Required Go version
REQUIRED_GO_VERSION="1.21"

# Required ports
REQUIRED_PORTS=(8080 50051 51051)

# Required disk space (in MB)
REQUIRED_DISK_SPACE=100

# Function to print colored messages
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Function to show usage
show_help() {
    cat << EOF
Usage: $(basename "$0")

Check if all prerequisites are met for running tests

CHECKS:
    - Go version (>= $REQUIRED_GO_VERSION)
    - Required ports available (${REQUIRED_PORTS[@]})
    - Disk space (>= ${REQUIRED_DISK_SPACE}MB)
    - Write permissions
    - Required commands (make, sqlite3, etc.)

EXIT CODES:
    0 - All prerequisites met
    1 - Missing prerequisites

EOF
}

# Check for help flag
for arg in "$@"; do
    if [[ "$arg" == "--help" ]]; then
        show_help
        exit 0
    fi
done

# Function to compare versions
version_ge() {
    # Returns 0 if $1 >= $2
    printf '%s\n%s\n' "$2" "$1" | sort -V -C
}

# Function to check Go version
check_go_version() {
    print_info "Checking Go version..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        return 1
    fi
    
    local go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
    
    if version_ge "$go_version" "$REQUIRED_GO_VERSION"; then
        print_success "Go version $go_version (>= $REQUIRED_GO_VERSION required)"
        return 0
    else
        print_error "Go version $go_version is too old (>= $REQUIRED_GO_VERSION required)"
        return 1
    fi
}

# Function to check if port is available
check_port() {
    local port=$1
    
    if command -v lsof &> /dev/null; then
        if lsof -ti:$port &> /dev/null; then
            return 1
        fi
    elif command -v netstat &> /dev/null; then
        if netstat -tuln | grep -q ":$port "; then
            return 1
        fi
    elif command -v ss &> /dev/null; then
        if ss -tuln | grep -q ":$port "; then
            return 1
        fi
    else
        # Fallback: try to bind to the port
        if (echo > /dev/tcp/127.0.0.1/$port) &> /dev/null; then
            return 1
        fi
    fi
    
    return 0
}

# Function to check required ports
check_ports() {
    print_info "Checking required ports..."
    
    local all_available=true
    
    for port in "${REQUIRED_PORTS[@]}"; do
        if check_port "$port"; then
            print_success "Port $port is available"
        else
            print_error "Port $port is already in use"
            all_available=false
        fi
    done
    
    if [[ "$all_available" == "true" ]]; then
        return 0
    else
        return 1
    fi
}

# Function to check disk space
check_disk_space() {
    print_info "Checking disk space..."
    
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local test_dir="$(dirname "$script_dir")"
    
    # Get available space in MB
    local available_space=0
    if command -v df &> /dev/null; then
        available_space=$(df -m "$test_dir" | awk 'NR==2 {print $4}')
    else
        print_warning "Cannot check disk space (df command not found)"
        return 0
    fi
    
    if [[ $available_space -ge $REQUIRED_DISK_SPACE ]]; then
        print_success "Disk space: ${available_space}MB available (>= ${REQUIRED_DISK_SPACE}MB required)"
        return 0
    else
        print_error "Disk space: ${available_space}MB available (>= ${REQUIRED_DISK_SPACE}MB required)"
        return 1
    fi
}

# Function to check write permissions
check_permissions() {
    print_info "Checking write permissions..."
    
    local script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local test_dir="$(dirname "$script_dir")"
    
    local dirs_to_check=("$test_dir/logs" "$test_dir/results" "$test_dir/tmp")
    local all_writable=true
    
    for dir in "${dirs_to_check[@]}"; do
        # Create directory if it doesn't exist
        mkdir -p "$dir" 2>/dev/null || true
        
        if [[ -w "$dir" ]]; then
            print_success "Write permission: $dir"
        else
            print_error "No write permission: $dir"
            all_writable=false
        fi
    done
    
    if [[ "$all_writable" == "true" ]]; then
        return 0
    else
        return 1
    fi
}

# Function to check required commands
check_commands() {
    print_info "Checking required commands..."
    
    local required_commands=("make" "git")
    local optional_commands=("sqlite3" "curl" "wget" "nc" "lsof")
    
    local all_required=true
    
    # Check required commands
    for cmd in "${required_commands[@]}"; do
        if command -v "$cmd" &> /dev/null; then
            print_success "Required command: $cmd"
        else
            print_error "Required command not found: $cmd"
            all_required=false
        fi
    done
    
    # Check optional commands
    for cmd in "${optional_commands[@]}"; do
        if command -v "$cmd" &> /dev/null; then
            print_success "Optional command: $cmd"
        else
            print_warning "Optional command not found: $cmd (recommended but not required)"
        fi
    done
    
    if [[ "$all_required" == "true" ]]; then
        return 0
    else
        return 1
    fi
}

# Function to check protobuf compiler
check_protoc() {
    print_info "Checking protobuf compiler..."
    
    if command -v protoc &> /dev/null; then
        local protoc_version=$(protoc --version | grep -oP '\d+\.\d+' || echo "unknown")
        print_success "protoc version $protoc_version"
        return 0
    else
        print_warning "protoc not found (required for rebuilding proto files)"
        return 0  # Not critical for running tests
    fi
}

# Function to check Go tools
check_go_tools() {
    print_info "Checking Go tools..."
    
    local tools=("protoc-gen-go" "protoc-gen-go-grpc")
    local all_found=true
    
    for tool in "${tools[@]}"; do
        if command -v "$tool" &> /dev/null; then
            print_success "Go tool: $tool"
        else
            print_warning "Go tool not found: $tool (required for rebuilding proto files)"
            all_found=false
        fi
    done
    
    return 0  # Not critical for running tests
}

# Main execution
main() {
    print_info "FlintRoute Test Prerequisites Check"
    print_info "===================================="
    echo ""
    
    local exit_code=0
    
    # Check Go version
    if ! check_go_version; then
        exit_code=1
    fi
    echo ""
    
    # Check required ports
    if ! check_ports; then
        exit_code=1
    fi
    echo ""
    
    # Check disk space
    if ! check_disk_space; then
        exit_code=1
    fi
    echo ""
    
    # Check write permissions
    if ! check_permissions; then
        exit_code=1
    fi
    echo ""
    
    # Check required commands
    if ! check_commands; then
        exit_code=1
    fi
    echo ""
    
    # Check protobuf compiler (optional)
    check_protoc
    echo ""
    
    # Check Go tools (optional)
    check_go_tools
    echo ""
    
    # Final result
    echo "========================================"
    if [[ $exit_code -eq 0 ]]; then
        print_success "All prerequisites met"
        print_info "System is ready for running tests"
    else
        print_error "Some prerequisites are missing"
        print_info "Please install missing requirements before running tests"
    fi
    
    exit $exit_code
}

# Run main function
main "$@"