#!/bin/bash

# FlintRoute Test Database Cleanup
# Clean test database with various options

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$(dirname "$SCRIPT_DIR")"
DB_FILE="$TEST_DIR/tmp/test.db"

# Default mode
MODE="full"

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

Clean test database with various cleanup modes

OPTIONS:
    --drop-tables    Drop all tables, keep database file
    --clear-data     Delete all data, keep schema
    --full           Remove database file completely (default)
    --help           Show this help message

MODES:
    drop-tables  - Drops all tables but keeps the database file
                   Useful for schema changes
    
    clear-data   - Deletes all data from tables but keeps schema
                   Useful for resetting test data
    
    full         - Removes the entire database file
                   Complete cleanup (default)

EXAMPLES:
    # Full cleanup (remove database file)
    ./cleanup-db.sh

    # Clear all data but keep schema
    ./cleanup-db.sh --clear-data

    # Drop all tables
    ./cleanup-db.sh --drop-tables

EXIT CODES:
    0 - Cleanup successful
    1 - Cleanup failed

EOF
}

# Parse arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --drop-tables)
                MODE="drop-tables"
                shift
                ;;
            --clear-data)
                MODE="clear-data"
                shift
                ;;
            --full)
                MODE="full"
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

# Function to check if database exists
check_database() {
    if [[ ! -f "$DB_FILE" ]]; then
        print_warning "Database file does not exist: $DB_FILE"
        return 1
    fi
    return 0
}

# Function to check if database is in use
check_database_in_use() {
    if command -v lsof > /dev/null 2>&1; then
        if lsof "$DB_FILE" > /dev/null 2>&1; then
            print_error "Database is currently in use by:"
            lsof "$DB_FILE"
            print_error "Please stop all services before cleaning the database"
            return 1
        fi
    fi
    return 0
}

# Function to drop all tables
drop_tables() {
    print_info "Dropping all tables from database..."
    
    if ! check_database; then
        print_success "No database to clean"
        return 0
    fi
    
    if ! check_database_in_use; then
        return 1
    fi
    
    # Use sqlite3 to drop all tables
    if command -v sqlite3 > /dev/null 2>&1; then
        # Get list of tables and drop them
        local tables=$(sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
        
        if [[ -z "$tables" ]]; then
            print_info "No tables to drop"
            return 0
        fi
        
        for table in $tables; do
            print_info "Dropping table: $table"
            sqlite3 "$DB_FILE" "DROP TABLE IF EXISTS $table;" || {
                print_error "Failed to drop table: $table"
                return 1
            }
        done
        
        print_success "All tables dropped"
    else
        print_warning "sqlite3 not found, falling back to full cleanup"
        rm -f "$DB_FILE"
        print_success "Database file removed"
    fi
}

# Function to clear all data
clear_data() {
    print_info "Clearing all data from database..."
    
    if ! check_database; then
        print_success "No database to clean"
        return 0
    fi
    
    if ! check_database_in_use; then
        return 1
    fi
    
    # Use sqlite3 to delete all data
    if command -v sqlite3 > /dev/null 2>&1; then
        # Get list of tables
        local tables=$(sqlite3 "$DB_FILE" "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
        
        if [[ -z "$tables" ]]; then
            print_info "No tables with data"
            return 0
        fi
        
        # Delete data from each table
        for table in $tables; do
            print_info "Clearing data from table: $table"
            sqlite3 "$DB_FILE" "DELETE FROM $table;" || {
                print_error "Failed to clear data from table: $table"
                return 1
            }
        done
        
        # Vacuum to reclaim space
        print_info "Vacuuming database..."
        sqlite3 "$DB_FILE" "VACUUM;" || true
        
        print_success "All data cleared"
    else
        print_warning "sqlite3 not found, falling back to full cleanup"
        rm -f "$DB_FILE"
        print_success "Database file removed"
    fi
}

# Function to perform full cleanup
full_cleanup() {
    print_info "Performing full database cleanup..."
    
    if ! check_database; then
        print_success "No database to clean"
        return 0
    fi
    
    if ! check_database_in_use; then
        return 1
    fi
    
    # Remove database file
    rm -f "$DB_FILE"
    
    # Also remove any journal or WAL files
    rm -f "${DB_FILE}-journal"
    rm -f "${DB_FILE}-wal"
    rm -f "${DB_FILE}-shm"
    
    print_success "Database file removed"
}

# Main execution
main() {
    print_info "FlintRoute Test Database Cleanup"
    print_info "================================="
    
    # Parse arguments
    parse_args "$@"
    
    print_info "Cleanup mode: $MODE"
    print_info "Database file: $DB_FILE"
    echo ""
    
    # Perform cleanup based on mode
    case "$MODE" in
        drop-tables)
            if drop_tables; then
                print_success "Database tables dropped successfully"
                exit 0
            else
                print_error "Failed to drop database tables"
                exit 1
            fi
            ;;
        clear-data)
            if clear_data; then
                print_success "Database data cleared successfully"
                exit 0
            else
                print_error "Failed to clear database data"
                exit 1
            fi
            ;;
        full)
            if full_cleanup; then
                print_success "Database cleanup completed successfully"
                exit 0
            else
                print_error "Failed to cleanup database"
                exit 1
            fi
            ;;
        *)
            print_error "Invalid mode: $MODE"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"