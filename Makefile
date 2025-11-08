.PHONY: help deps build clean dev-backend dev-frontend dev test docker-up docker-down install

# Default target
help:
	@echo "FlintRoute - BGP Management System"
	@echo ""
	@echo "Available targets:"
	@echo "  make deps          - Install all dependencies"
	@echo "  make build         - Build backend binary"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make dev-backend   - Run backend in development mode"
	@echo "  make dev-frontend  - Run frontend in development mode"
	@echo "  make dev           - Run both backend and frontend"
	@echo "  make test          - Run tests"
	@echo "  make docker-up     - Start FRR test environment"
	@echo "  make docker-down   - Stop FRR test environment"
	@echo "  make install       - Install the application"

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build backend
build:
	@echo "Building backend..."
	mkdir -p bin
	go build -o bin/flintroute ./cmd/flintroute
	@echo "Build complete: bin/flintroute"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf data/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/
	@echo "Clean complete"

# Run backend in development mode
dev-backend:
	@echo "Starting backend in development mode..."
	@if [ ! -f configs/config.yaml ]; then \
		echo "Creating config.yaml from example..."; \
		cp configs/config.example.yaml configs/config.yaml; \
	fi
	go run ./cmd/flintroute

# Run frontend in development mode
dev-frontend:
	@echo "Starting frontend in development mode..."
	cd frontend && npm run dev

# Run both backend and frontend
dev:
	@echo "Starting FlintRoute in development mode..."
	@make -j2 dev-backend dev-frontend

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete"

# Start Docker Compose for FRR testing
docker-up:
	@echo "Starting FRR test environment..."
	docker-compose up -d
	@echo "FRR is starting... Use 'docker-compose logs -f' to view logs"

# Stop Docker Compose
docker-down:
	@echo "Stopping FRR test environment..."
	docker-compose down
	@echo "FRR stopped"

# Install the application
install: build
	@echo "Installing FlintRoute..."
	sudo cp bin/flintroute /usr/local/bin/
	@echo "Installation complete. Run 'flintroute' to start."

# Development with hot reload (requires air)
dev-hot:
	@echo "Starting backend with hot reload..."
	@if ! command -v air > /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air