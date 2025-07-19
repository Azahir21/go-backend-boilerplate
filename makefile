.PHONY: help build run dev test test-coverage lint fmt vet clean docker-build docker-up docker-down docker-logs docker-shell docker-clean migrate seed swag setup

# Default target
help: ## Show this help message
    @echo 'Usage: make [target]'
	@echo ''
    @echo 'Targets:'
    @awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands
dev: ## Start development server with hot reload (requires Air)
	@echo "Starting development server with hot reload..."
	air

run: ## Start production server
	@echo "Starting production server..."
	go run main.go

build: ## Build the application binary
	@echo "Building application..."
	go build -o bin/go-backend-boilerplate main.go

# Testing commands
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality commands
lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code with go fmt
	@echo "Formatting code..."
	go fmt ./...

vet: ## Analyze code with go vet
	@echo "Running go vet..."
	go vet ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Documentation
swag: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	swag init

# Database commands
migrate: ## Run database migrations (GORM auto-migrate)
	@echo "Running database migrations..."
	go run scripts/migrate.go

seed: ## Seed the database with sample data
	@echo "Seeding database..."
	go run scripts/seed.go

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t go-backend-boilerplate .

docker-up: ## Start application with Docker Compose
	@echo "Starting application with Docker Compose..."
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## View Docker container logs
	@echo "Viewing container logs..."
	docker-compose logs -f app

docker-shell: ## Access application container shell
	@echo "Accessing container shell..."
	docker-compose exec app sh

docker-clean: ## Clean up Docker resources
	@echo "Cleaning up Docker resources..."
	docker-compose down -v
	docker system prune -f

# Environment setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env file with your configuration"
	go mod tidy
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger documentation..."
	swag init
	@echo "Development environment setup complete!"

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean