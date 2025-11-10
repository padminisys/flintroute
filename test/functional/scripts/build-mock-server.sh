#!/bin/bash

# FlintRoute Mock FRR Server Builder
# Build the mock FRR server binary

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
MOCK_SERVER_DIR="$TEST_DIR/cmd/mock-frr-server"
OUTPUT_BINARY="$MOCK_SERVER_DIR/mock-frr-server"

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

Build the mock FRR server binary

This script:
  1. Checks if Go is installed
  2. Generates protobuf code (if needed)
  3. Builds the mock FRR server binary
  4. Verifies the build

OPTIONS:
    --clean      Clean build (remove old binary first)
    --verbose    Verbose build output
    --help       Show this help message

OUTPUT:
    Binary: $OUTPUT_BINARY

EXAMPLES:
    # Normal build
    ./build-mock-server.sh

    # Clean build with verbose output
    ./build-mock-server.sh --clean --verbose

EXIT CODES:
    0 - Build successful
    1 - Build failed

EOF
}

# Parse arguments
CLEAN=false
VERBOSE=false

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --clean)
                CLEAN=true
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

# Function to check Go installation
check_go() {
    print_info "Checking Go installation..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        return 1
    fi
    
    local go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "0.0")
    print_success "Go version $go_version"
}

# Function to clean old binary
clean_binary() {
    if [[ "$CLEAN" == "true" ]]; then
        print_info "Cleaning old binary..."
        
        if [[ -f "$OUTPUT_BINARY" ]]; then
            rm -f "$OUTPUT_BINARY"
            print_success "Old binary removed"
        else
            print_info "No old binary to remove"
        fi
    fi
}

# Function to check if protobuf generation is needed
check_protobuf() {
    print_info "Checking protobuf files..."
    
    local proto_file="$MOCK_SERVER_DIR/proto/frr.proto"
    local pb_go_file="$MOCK_SERVER_DIR/proto/frr.pb.go"
    local grpc_go_file="$MOCK_SERVER_DIR/proto/frr_grpc.pb.go"
    
    if [[ ! -f "$proto_file" ]]; then
        print_error "Proto file not found: $proto_file"
        return 1
    fi
    
    # Check if generated files exist and are newer than proto file
    if [[ -f "$pb_go_file" ]] && [[ -f "$grpc_go_file" ]]; then
        if [[ "$pb_go_file" -nt "$proto_file" ]] && [[ "$grpc_go_file" -nt "$proto_file" ]]; then
            print_success "Protobuf files are up to date"
            return 0
        fi
    fi
    
    print_warning "Protobuf files need to be generated"
    return 1
}

# Function to generate protobuf code
generate_protobuf() {
    if check_protobuf; then
        return 0
    fi
    
    print_info "Generating protobuf code..."
    
    # Check if protoc is available
    if ! command -v protoc &> /dev/null; then
        print_warning "protoc not found, skipping protobuf generation"
        print_warning "Using existing generated files (if any)"
        return 0
    fi
    
    # Check if Go protobuf plugins are available
    if ! command -v protoc-gen-go &> /dev/null || ! command -v protoc-gen-go-grpc &> /dev/null; then
        print_warning "protoc-gen-go or protoc-gen-go-grpc not found"
        print_warning "Using existing generated files (if any)"
        return 0
    fi
    
    # Generate protobuf code
    cd "$MOCK_SERVER_DIR"
    
    if [[ "$VERBOSE" == "true" ]]; then
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               proto/frr.proto
    else
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               proto/frr.proto > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Protobuf code generated"
    else
        print_error "Failed to generate protobuf code"
        return 1
    fi
    
    cd "$TEST_DIR"
}

# Function to build mock server
build_server() {
    print_info "Building mock FRR server..."
    
    cd "$MOCK_SERVER_DIR"
    
    # Build command
    local build_cmd="go build -o mock-frr-server"
    
    if [[ "$VERBOSE" == "true" ]]; then
        build_cmd="$build_cmd -v"
    fi
    
    # Add build flags for optimization
    build_cmd="$build_cmd -ldflags='-s -w'"
    
    # Add source files
    build_cmd="$build_cmd ."
    
    print_info "Executing: $build_cmd"
    
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$build_cmd"
    else
        eval "$build_cmd" > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Mock FRR server built successfully"
    else
        print_error "Failed to build mock FRR server"
        return 1
    fi
    
    cd "$TEST_DIR"
}

# Function to verify build
verify_build() {
    print_info "Verifying build..."
    
    if [[ ! -f "$OUTPUT_BINARY" ]]; then
        print_error "Binary not found: $OUTPUT_BINARY"
        return 1
    fi
    
    if [[ ! -x "$OUTPUT_BINARY" ]]; then
        print_error "Binary is not executable: $OUTPUT_BINARY"
        return 1
    fi
    
    # Get binary size
    local size=$(du -h "$OUTPUT_BINARY" | cut -f1)
    print_success "Binary created: $OUTPUT_BINARY ($size)"
    
    # Test if binary can show version/help
    if [[ "$VERBOSE" == "true" ]]; then
        print_info "Testing binary..."
        if "$OUTPUT_BINARY" --help > /dev/null 2>&1 || "$OUTPUT_BINARY" -h > /dev/null 2>&1; then
            print_success "Binary is functional"
        else
            print_warning "Binary help test inconclusive (may be normal)"
        fi
    fi
}

# Main execution
main() {
    print_info "FlintRoute Mock FRR Server Builder"
    print_info "==================================="
    
    # Parse arguments
    parse_args "$@"
    
    echo ""
    
    # Check Go installation
    if ! check_go; then
        exit 1
    fi
    
    echo ""
    
    # Clean old binary if requested
    clean_binary
    
    echo ""
    
    # Generate protobuf code if needed
    if ! generate_protobuf; then
        print_warning "Protobuf generation had issues, continuing anyway..."
    fi
    
    echo ""
    
    # Build server
    if ! build_server; then
        exit 1
    fi
    
    echo ""
    
    # Verify build
    if ! verify_build; then
        exit 1
    fi
    
    echo ""
    print_success "Build completed successfully"
    print_info "Binary location: $OUTPUT_BINARY"
    
    exit 0
}

# Run main function
main "$@"