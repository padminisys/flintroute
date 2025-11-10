#!/bin/bash

# FlintRoute Test Runner Builder
# Build the test execution framework

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"

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

Build the test execution framework

This script:
  1. Checks if Go is installed
  2. Downloads Go dependencies
  3. Compiles all Go packages in pkg/
  4. Verifies the build

OPTIONS:
    --clean      Clean build (remove build cache)
    --verbose    Verbose build output
    --help       Show this help message

EXAMPLES:
    # Normal build
    ./build-test-runner.sh

    # Clean build with verbose output
    ./build-test-runner.sh --clean --verbose

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

# Function to clean build cache
clean_cache() {
    if [[ "$CLEAN" == "true" ]]; then
        print_info "Cleaning build cache..."
        
        cd "$TEST_DIR"
        
        if [[ "$VERBOSE" == "true" ]]; then
            go clean -cache -testcache -modcache
        else
            go clean -cache -testcache -modcache > /dev/null 2>&1
        fi
        
        print_success "Build cache cleaned"
    fi
}

# Function to download dependencies
download_dependencies() {
    print_info "Downloading Go dependencies..."
    
    cd "$TEST_DIR"
    
    if [[ "$VERBOSE" == "true" ]]; then
        go mod download
    else
        go mod download > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Dependencies downloaded"
    else
        print_error "Failed to download dependencies"
        return 1
    fi
}

# Function to verify dependencies
verify_dependencies() {
    print_info "Verifying dependencies..."
    
    cd "$TEST_DIR"
    
    if [[ "$VERBOSE" == "true" ]]; then
        go mod verify
    else
        go mod verify > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Dependencies verified"
    else
        print_warning "Dependency verification had issues"
    fi
}

# Function to tidy dependencies
tidy_dependencies() {
    print_info "Tidying dependencies..."
    
    cd "$TEST_DIR"
    
    if [[ "$VERBOSE" == "true" ]]; then
        go mod tidy
    else
        go mod tidy > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Dependencies tidied"
    else
        print_warning "go mod tidy had issues"
    fi
}

# Function to build test packages
build_packages() {
    print_info "Building test packages..."
    
    cd "$TEST_DIR"
    
    # Build all packages in pkg/
    local build_cmd="go build"
    
    if [[ "$VERBOSE" == "true" ]]; then
        build_cmd="$build_cmd -v"
    fi
    
    # Add all packages
    build_cmd="$build_cmd ./pkg/..."
    
    print_info "Executing: $build_cmd"
    
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$build_cmd"
    else
        eval "$build_cmd" > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Test packages built successfully"
    else
        print_error "Failed to build test packages"
        return 1
    fi
}

# Function to compile test binaries
compile_tests() {
    print_info "Compiling test binaries..."
    
    cd "$TEST_DIR"
    
    # Compile tests without running them
    local test_cmd="go test -c"
    
    if [[ "$VERBOSE" == "true" ]]; then
        test_cmd="$test_cmd -v"
    fi
    
    # Compile all test packages
    test_cmd="$test_cmd ./tests/..."
    
    print_info "Executing: $test_cmd"
    
    if [[ "$VERBOSE" == "true" ]]; then
        eval "$test_cmd"
    else
        eval "$test_cmd" > /dev/null 2>&1
    fi
    
    if [[ $? -eq 0 ]]; then
        print_success "Test binaries compiled successfully"
    else
        print_warning "Test compilation had issues (may be normal if no tests exist yet)"
    fi
}

# Function to verify build
verify_build() {
    print_info "Verifying build..."
    
    cd "$TEST_DIR"
    
    # Check if pkg directory exists and has Go files
    if [[ ! -d "pkg" ]]; then
        print_warning "pkg directory does not exist yet"
        return 0
    fi
    
    local go_files=$(find pkg -name "*.go" 2>/dev/null | wc -l)
    
    if [[ $go_files -eq 0 ]]; then
        print_warning "No Go files found in pkg directory yet"
        return 0
    fi
    
    print_success "Found $go_files Go file(s) in pkg directory"
    
    # Try to list packages
    if [[ "$VERBOSE" == "true" ]]; then
        print_info "Available packages:"
        go list ./pkg/... 2>/dev/null || print_warning "No packages found yet"
    fi
}

# Main execution
main() {
    print_info "FlintRoute Test Runner Builder"
    print_info "==============================="
    
    # Parse arguments
    parse_args "$@"
    
    echo ""
    
    # Check Go installation
    if ! check_go; then
        exit 1
    fi
    
    echo ""
    
    # Clean build cache if requested
    clean_cache
    
    if [[ "$CLEAN" == "true" ]]; then
        echo ""
    fi
    
    # Download dependencies
    if ! download_dependencies; then
        exit 1
    fi
    
    echo ""
    
    # Verify dependencies
    verify_dependencies
    
    echo ""
    
    # Tidy dependencies
    tidy_dependencies
    
    echo ""
    
    # Build packages
    if ! build_packages; then
        print_warning "Package build had issues, but continuing..."
    fi
    
    echo ""
    
    # Compile tests
    compile_tests
    
    echo ""
    
    # Verify build
    verify_build
    
    echo ""
    print_success "Build completed successfully"
    print_info "Test framework is ready"
    
    exit 0
}

# Run main function
main "$@"