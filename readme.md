# Go Backend Boilerplate

[![Go Report Card](https://goreportcard.com/badge/github.com/azahir21/go-backend-boilerplate)](https://goreportcard.com/report/github.com/azahir21/go-backend-boilerplate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust and scalable Go backend boilerplate designed to jumpstart your next project. It comes with out-of-the-box support for REST, gRPC, and GraphQL APIs, leveraging modern Go practices and essential infrastructure components.

## Features

-   **Multiple API Gateways**:
    -   **REST API**: Built with Gin framework, including Swagger documentation.
    -   **gRPC API**: For high-performance, language-agnostic communication.
    -   **GraphQL API**: Powered by `graphql-go` for flexible data fetching, with GraphiQL playground.
-   **Database Integration**:
    -   **Ent ORM**: Type-safe and powerful ORM for database interactions.
    -   Supports PostgreSQL, MySQL, and SQLite.
    -   Automatic database migrations.
-   **Authentication & Authorization**:
    -   JWT-based authentication.
    -   Role-based access control (Basic Admin/User roles).
-   **Caching**:
    -   Flexible caching layer with support for Redis and Ristretto (in-memory).
-   **File Storage**:
    -   Pluggable storage module with support for Local filesystem, AWS S3, and Google Cloud Storage (GCS).
-   **Email Services**:
    -   Transactional email support with SMTP and SendGrid integrations.
-   **Structured Logging**: Implemented with `logrus` for clear and customizable logging.
-   **Configuration Management**: Centralized configuration using Viper, supporting YAML files and environment variables.
-   **Unit of Work Pattern**: Ensures atomic database operations.
-   **Graceful Shutdown**: Handles application shutdown cleanly for all running services.
-   **Docker Support**: Ready-to-use Dockerfile for containerization.
-   **Project Structure**: Clear and maintainable directory layout.


## Build Tags: Modular Delivery Layers

This boilerplate supports **compile-time selection** of delivery layers using Go build tags. This allows you to:
- Build only what you need (reducing binary size and dependencies)
- Skip unnecessary tooling (e.g., no protoc needed for REST-only builds)
- Choose any combination of REST, gRPC, and GraphQL

### Available Build Combinations

| Build Tags | Delivery Layers | Build Command | Binary Size* |
|------------|----------------|---------------|--------------|
| `rest` | REST only | `go build -tags rest` | ~94MB |
| `grpc` | gRPC only | `go build -tags grpc` | ~82MB |
| `graphql` | GraphQL only | `go build -tags graphql` | ~85MB |
| `rest,grpc` | REST + gRPC | `go build -tags "rest,grpc"` | ~95MB |
| `rest,graphql` | REST + GraphQL | `go build -tags "rest,graphql"` | ~95MB |
| `grpc,graphql` | gRPC + GraphQL | `go build -tags "grpc,graphql"` | ~85MB |
| `rest,grpc,graphql` | All layers | `go build -tags "rest,grpc,graphql"` | ~95MB |

*Approximate binary sizes (unstripped, debug symbols included)

### Quick Start with Build Tags

#### 1. REST-only (Recommended for beginners)
```bash
# No protoc or GraphQL dependencies needed
make build-rest
./bin/go-backend-boilerplate-rest
```

#### 2. gRPC-only
```bash
# Requires protoc and generated proto files
make setup-grpc  # Install protoc tools and generate proto files
make build-grpc
./bin/go-backend-boilerplate-grpc
```

#### 3. All delivery layers
```bash
# Requires all dependencies
make setup       # Install all tools
make setup-grpc  # Install gRPC tools
make build-all
./bin/go-backend-boilerplate-all
```

### Makefile Targets

```bash
# Build targets
make build-rest          # REST-only
make build-grpc          # gRPC-only
make build-graphql       # GraphQL-only
make build-rest-grpc     # REST + gRPC
make build-rest-graphql  # REST + GraphQL
make build-grpc-graphql  # gRPC + GraphQL
make build-all           # All delivery layers

# Setup targets
make setup               # Basic setup (REST dependencies)
make setup-grpc          # Install gRPC/protoc tools
make generate-proto      # Generate protobuf files (needed for gRPC)
```

### Configuration vs Build Tags

**Build tags** control **what code is compiled**.  
**Config files** control **what servers start at runtime**.

Example: Build with REST + gRPC support, but only run REST:
```yaml
# config.yaml
server:
  http_server:
    enable: true
  grpc_server:
    enable: false  # Won't start, but code is compiled
```

Build command:
```bash
go build -tags "rest,grpc" -o app ./cmd
```

### Required Dependencies by Build Tag

| Build Tag | Required Dependencies | Setup Command |
|-----------|----------------------|---------------|
| `rest` | Go, Gin, Swagger | `make setup` |
| `grpc` | Go, gRPC, protoc, protoc-gen-go* | `make setup-grpc` |
| `graphql` | Go, graphql-go | `make setup` |

*protoc must be installed separately: https://grpc.io/docs/protoc-installation/

### Migration from Previous Versions

If you're upgrading from a version without build tags:

**Old way (all layers always compiled):**
```bash
go build -o app ./cmd
```

**New way (explicit layer selection):**
```bash
# Choose one:
go build -tags rest -o app ./cmd              # REST only
go build -tags "rest,grpc,graphql" -o app ./cmd  # All layers (equivalent to old behavior)
```

For full backward compatibility (all layers), use: `make build` or `make build-all`

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

*   Go (version 1.22 or higher)
*   Docker & Docker Compose (optional, for database/redis setup)
*   Git

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/go-backend-boilerplate.git
cd go-backend-boilerplate
```

### 2. Configuration

The project uses `config.yaml` for base configuration and can be overridden by environment-specific files (e.g., `config.development.yaml`, `config.production.yaml`) or environment variables.

Copy the example environment file:
```bash
cp .env.example .env
```
Edit `.env` to set your environment variables, especially for JWT secret and database credentials.

Example `config.development.yaml` overrides:
```yaml
server:
  env: "development"
  http_server:
    enable: true
    port: "8080"
    cors_origins: ["http://localhost:3000"]
    read_timeout: "5s"
    write_timeout: "10s"
    idle_timeout: "120s"
    startup_banner: true
  grpc_server:
    enable: false
  graphql_server:
    enable: true
    port: "8081"
    cors_origins: ["http://localhost:3000", "http://localhost:8081/graphql/playground"]
    read_timeout: "5s"
    write_timeout: "10s"
    idle_timeout: "120s"
    startup_banner: true

database:
  driver: "sqlite3"
  name: "./data/app.db" # For SQLite, this is the file path
  auto_migrate: true
```

### 3. Database Setup (using Docker Compose for PostgreSQL/MySQL/Redis)

For development, you can use Docker Compose to spin up a PostgreSQL, MySQL, or Redis instance.

```bash
# To start PostgreSQL and Redis:
docker-compose up -d postgres redis

# Or for MySQL and Redis:
# docker-compose up -d mysql redis
```

Update your `config.yaml` or environment variables with the correct database and Redis connection details. If using `sqlite3`, no Docker setup is strictly needed as it uses a local file.

### 4. Ent Migrations

The boilerplate uses Ent ORM. Schema migrations are handled automatically on application startup if `database.auto_migrate` is set to `true` in your configuration.

To explicitly generate Ent code and run migrations:

```bash
# Generate Ent code
go generate ./ent

# Running migrations (handled by auto_migrate on app start, but good to know)
# go run -mod=mod entgo.io/ent/cmd/ent migrate --path ./migrations --dialect <your-db-driver>
```

### 5. Run the Application

#### Development (using `air` for live reload)

```bash
# Install air if you haven't already
go install github.com/cosmtrek/air@latest

# Run the application with live reload
air
```

#### Standard Go Run

```bash
go run cmd/main.go
```

The application will start the enabled servers (HTTP/REST, gRPC, GraphQL) on their configured ports.

-   **REST API**: Typically on `http://localhost:8080/api/v1`
-   **Swagger UI**: `http://localhost:8080/swagger/index.html`
-   **GraphQL Playground**: `http://localhost:8081/graphql/playground` (if enabled on port `8081`)

## Project Structure

```
.
├── .air.toml                  # Air configuration for live reload
├── .env.example               # Example environment variables
├── .gitignore
├── buf.yaml                   # Buf configuration for Protobuf
├── config.development.yaml    # Development environment configuration
├── config.yaml                # Base configuration
├── docker-compose.yaml        # Docker Compose for dev services (DB, Redis)
├── Dockerfile                 # Dockerfile for the application
├── go.mod                     # Go modules
├── go.sum
├── LICENSE
├── makefile                   # Common commands (build, run, test etc.)
├── readme.md
├── cmd/                       # Application entry points
│   ├── main.go                # Main application entry
│   └── app/                   # Application core logic and setup
│       └── app.go
│   └── service/               # Server implementations (REST, gRPC, GraphQL)
│       ├── graphql_service.go
│       ├── grpc_server.go
│       └── rest_server.go
├── docs/                      # API Documentation (Swagger)
├── ent/                       # Ent ORM generated code and schema definitions
│   └── schema/                # Database schema definitions
│       └── user.go
├── infrastructure/            # Infrastructure concerns (DB, Cache, Storage, External services)
│   ├── cache/
│   ├── db/
│   ├── external/
│   └── storage/
├── internal/                  # Internal business logic and application features
│   ├── shared/                # Shared utilities, entities, errors, middlewares
│   │   ├── cache/
│   │   ├── entity/
│   │   ├── errors/
│   │   ├── helper/
│   │   ├── http/
│   │   ├── middleware/
│   │   ├── storage/
│   │   └── unitofwork/
│   └── user/                  # User domain module (delivery, repository, usecase)
│       ├── delivery/          # API delivery (HTTP, gRPC, GraphQL handlers)
│       ├── repository/
│       └── usecase/
├── migrations/                # Database migration scripts
├── pkg/                       # Reusable packages/libraries (config, logger, httpresp)
├── proto/                     # Protobuf definitions and generated Go code
└── web/                       # Static web assets
    ├── playground.html
    └── css/
    └── js/
```

## API Endpoints

The boilerplate can enable HTTP (REST), gRPC, and GraphQL servers.

### REST API

-   **Base Path**: `/api/v1`
-   **Swagger Documentation**: `http://localhost:8080/swagger/index.html`
-   **Example Endpoints**:
    -   `GET /api/v1/ping`: Health check.
    -   `POST /api/v1/auth/register`: Register a new user.
    -   `POST /api/v1/auth/login`: User login (returns JWT token).
    -   `GET /api/v1/auth/profile`: Get user profile (requires JWT).
    -   `GET /api/v1/admin/test`: Admin-only example (requires admin JWT).

### gRPC API

-   **Port**: Configurable (e.g., `8090`)
-   **Service**: Defined in `proto/user.proto`
-   **Generated Code**: `proto/user.pb.go`, `proto/user_grpc.pb.go`
-   Use `grpc_cli` or a gRPC client library to interact.

### GraphQL API

-   **Port**: Configurable (e.g., `8081`)
-   **Endpoint**: `/graphql`
-   **Playground**: `http://localhost:8081/graphql/playground`
-   **Example Query (from playground default)**:
    ```graphql
    query RebelsShipsQuery {
      rebels {
        name
        ships(first: 1) {
          edges {
            node {
              name
            }
          }
        }
      }
    }
    ```
    (Note: This example query might need to be adapted to the actual schema defined in `internal/user/delivery/graphql/user_schema.go`)

## Authentication

The boilerplate uses JWT (JSON Web Tokens) for authentication.

-   **Login/Register**: Upon successful login or registration, a JWT token is returned.
-   **Protected Endpoints**: Include the JWT in the `Authorization` header as `Bearer <token>`.
-   **Middleware**: `middleware.AuthMiddleware()` verifies tokens. `middleware.AdminMiddleware()` checks for admin role.

## Error Handling

API responses for errors are standardized using `pkg/httpresp.JSON`.

## Testing

To run unit and integration tests:

```bash
go test ./...
```

For specific packages:
```bash
go test ./internal/user/...
```

## Deployment

The provided `Dockerfile` allows for easy containerization of the application.

```bash
docker build -t go-backend-boilerplate .
docker run -p 8080:8080 -p 8090:8090 -p 8081:8081 go-backend-boilerplate
```

Remember to bind mount your `config.yaml` and `.env` or set appropriate environment variables in your deployment environment.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
