.PHONY: help build run dev test test-coverage lint fmt vet clean docker-build docker-up docker-down docker-logs docker-shell docker-clean goose-create goose-up goose-down goose-status swag setup generate

PROTO_ROOT=internal
PROTO_FILES=$(shell find $(PROTO_ROOT) -name "*.proto")

SWAG_MAIN=cmd/main.go

# Default target
help: ## Show this help message
    @echo 'Usage: make [target]'
	@echo ''
    @echo 'Targets:'
    @awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $1, $2}' $(MAKEFILE_LIST)

# Development commands
dev: ## Start development server with hot reload (requires Air)
	@echo "Starting development server with hot reload..."
	air

run: ## Start production server
	@echo "Starting production server..."
	go run cmd/main.go

run-rest: ## Start REST server
	@echo "Starting REST server..."
	go run cmd/rest/main.go

run-graphql: ## Start GraphQL server
	@echo "Starting GraphQL server..."
	go run cmd/graphql/main.go

run-grpc: ## Start gRPC server
	@echo "Starting gRPC server..."
	go run cmd/grpc/main.go

build: ## Build all three delivery binaries
	@echo "Building all delivery binaries..."
	@$(MAKE) build-rest
	@$(MAKE) build-graphql
	@$(MAKE) build-grpc

build-rest: ## Build REST-only binary
	@echo "Building REST-only application..."
	go build -o bin/go-backend-boilerplate-rest ./cmd/rest

build-graphql: ## Build GraphQL-only binary
	@echo "Building GraphQL-only application..."
	go build -o bin/go-backend-boilerplate-graphql ./cmd/graphql

build-grpc: ## Build gRPC-only binary
	@echo "Building gRPC-only application..."
	go build -o bin/go-backend-boilerplate-grpc ./cmd/grpc

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
	swag init \
		-g $(SWAG_MAIN) \
		-dir . \
		--exclude internal/*/delivery/grpc,internal/*/delivery/grpc/**,internal/*/delivery/graphql,internal/*/delivery/graphql/** \
		--parseInternal

# Generate Proto
proto:
	@echo "Generating protobuf code..."
	for proto in $(PROTO_FILES); do \
		proto_dir=$$(dirname $$proto); \
		gen_dir=$$(echo $$proto_dir | sed 's|/proto$$|/gen|'); \
		mkdir -p $$gen_dir; \
		protoc \
			--proto_path=$$proto_dir \
			--go_out=$$gen_dir \
			--go_opt=paths=source_relative \
			--go-grpc_out=$$gen_dir \
			--go-grpc_opt=paths=source_relative \
			$$proto ; \
	done

# Code generation
generate: ## Generate code for ent and protobuf
	@echo "Generating ent code..."
	go generate ./ent
	@echo "Generating protobuf code..."
	$(MAKE) proto
	@echo "Generating swag documentation..."
	$(MAKE) swag


goose-create: ## Create a new Goose migration file
	@echo "Creating new Goose migration..."
	goose create ${NAME} sql

goose-up: ## Apply pending Goose migrations
	@echo "Applying Goose migrations..."
	goose up

goose-down: ## Rollback the last Goose migration
	@echo "Rolling back last Goose migration..."
	goose down
goose-status: ## Check Goose migration status
	@echo "Checking Goose migration status..."
	@echo "Make sure the 'GOOSE_DB_STRING', 'GOOSE_DRIVER', and 'GOOSE_MIGRATION_DIR' environment variables are set."
	@echo "Checking .env and GOOSE_DB_STRING..."
	@if [ -f .env ]; then \
	  echo ".env found — loading variables..."; \
	  set -a; . .env; set +a; \
	  goose status; \
	else \
	  if [ -n "$$GOOSE_DB_STRING" ]; then \
		echo "GOOSE_DB_STRING is set in environment — running goose status..."; \
		goose status; \
	  else \
		echo ".env not found and GOOSE_DB_STRING is not set."; \
		echo "Create a .env from .env.example or export GOOSE_DB_STRING in your shell."; \
		exit 1; \
	  fi; \
	fi

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
	@if [ ! -f .env ]; then cp .env.example .env; echo "Please edit .env file with your configuration"; fi
	go mod tidy
	@echo "Installing development tools..."
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Generating Swagger documentation..."
	$(MAKE) swag
	@echo "Generating ent and protobuf code..."
	make generate
	@echo "Development environment setup complete!"

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean