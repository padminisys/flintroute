#!/bin/bash

# FlintRoute Test Results Cleanup
# Clean previous test results

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
RESULTS_DIR="$TEST_DIR/results"

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

Clean previous test results

This script removes:
  - JSON test results (*.json)
  - XML test results (*.xml)
  - HTML test reports (*.html)
  - Test summary files (*.txt)

The .gitkeep file is preserved to maintain the directory structure.

OPTIONS:
    --help    Show this help message

EXAMPLES:
    # Clean all test results
    ./cleanup-results.sh

EXIT CODES:
    0 - Cleanup successful
    1 - Cleanup failed

EOF
}

# Check for help flag
for arg in "$@"; do
    if [[ "$arg" == "--help" ]]; then
        show_help
        exit 0
    fi
done

# Function to count files
count_files() {
    local pattern=$1
    local count=0
    
    if compgen -G "$RESULTS_DIR/$pattern" > /dev/null 2>&1; then
        count=$(ls -1 "$RESULTS_DIR"/$pattern 2>/dev/null | wc -l)
    fi
    
    echo "$count"
}

# Function to cleanup results
cleanup_results() {
    print_info "Cleaning test results directory: $RESULTS_DIR"
    
    if [[ ! -d "$RESULTS_DIR" ]]; then
        print_warning "Results directory does not exist"
        return 0
    fi
    
    local total_removed=0
    
    # Count files before cleanup
    local json_count=$(count_files "*.json")
    local xml_count=$(count_files "*.xml")
    local html_count=$(count_files "*.html")
    local txt_count=$(count_files "*.txt")
    
    print_info "Found files to remove:"
    print_info "  - JSON files: $json_count"
    print_info "  - XML files: $xml_count"
    print_info "  - HTML files: $html_count"
    print_info "  - Text files: $txt_count"
    
    # Remove JSON files
    if [[ $json_count -gt 0 ]]; then
        rm -f "$RESULTS_DIR"/*.json
        print_info "Removed $json_count JSON file(s)"
        total_removed=$((total_removed + json_count))
    fi
    
    # Remove XML files
    if [[ $xml_count -gt 0 ]]; then
        rm -f "$RESULTS_DIR"/*.xml
        print_info "Removed $xml_count XML file(s)"
        total_removed=$((total_removed + xml_count))
    fi
    
    # Remove HTML files
    if [[ $html_count -gt 0 ]]; then
        rm -f "$RESULTS_DIR"/*.html
        print_info "Removed $html_count HTML file(s)"
        total_removed=$((total_removed + html_count))
    fi
    
    # Remove text files
    if [[ $txt_count -gt 0 ]]; then
        rm -f "$RESULTS_DIR"/*.txt
        print_info "Removed $txt_count text file(s)"
        total_removed=$((total_removed + txt_count))
    fi
    
    # Verify .gitkeep is still there
    if [[ ! -f "$RESULTS_DIR/.gitkeep" ]]; then
        touch "$RESULTS_DIR/.gitkeep"
        print_info "Recreated .gitkeep file"
    fi
    
    if [[ $total_removed -eq 0 ]]; then
        print_info "No files to remove"
    else
        print_success "Removed $total_removed file(s) total"
    fi
}

# Main execution
main() {
    print_info "FlintRoute Test Results Cleanup"
    print_info "================================"
    echo ""
    
    if cleanup_results; then
        echo ""
        print_success "Test results cleanup completed successfully"
        exit 0
    else
        echo ""
        print_error "Test results cleanup failed"
        exit 1
    fi
}

# Run main function
main "$@"