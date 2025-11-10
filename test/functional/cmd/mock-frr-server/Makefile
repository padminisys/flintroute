.PHONY: all build test clean run help

# Binary name
BINARY=mock-frr-server

# Default configuration
CONFIG=../../config/mock-frr-config.yaml

all: build

# Build the mock server
build:
	@echo "Building mock FRR server..."
	@go build -o $(BINARY)
	@echo "Build complete: $(BINARY)"

# Run tests
test: build
	@echo "Running tests..."
	@./test-server.sh

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY)
	@rm -f /tmp/mock-frr-server.log
	@echo "Clean complete"

# Run the server
run: build
	@echo "Starting mock FRR server..."
	@./$(BINARY) -config $(CONFIG)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

# Show help
help:
	@echo "Mock FRR Server Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the mock server binary"
	@echo "  make test     - Run tests"
	@echo "  make run      - Build and run the server"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make help     - Show this help message"
	@echo ""
	@echo "Configuration:"
	@echo "  CONFIG=$(CONFIG)"