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
    -   **SQL Databases**: Ent ORM with support for PostgreSQL, MySQL, SQLite, and MariaDB.
    -   **MongoDB**: Native MongoDB driver support for NoSQL workloads.
    -   **Multi-Database Support**: Run SQL and MongoDB simultaneously for polyglot persistence.
    -   **Optional Infrastructure**: All databases are optional and can be disabled via configuration.
    -   Automatic database migrations (SQL only).
-   **Authentication & Authorization**:
    -   JWT-based authentication.
    -   Role-based access control (Basic Admin/User roles).
-   **Caching**:
    -   Flexible caching layer with support for Redis and Ristretto (in-memory).
    -   Optional: Can be disabled if not needed.
-   **File Storage**:
    -   Pluggable storage module with support for Local filesystem, AWS S3, and Google Cloud Storage (GCS).
    -   Optional: Can be disabled if not needed.
-   **Email Services**:
    -   Transactional email support with SMTP and SendGrid integrations.
    -   Optional: Can be disabled if not needed.
-   **Optional Infrastructure**: All infrastructure components (database, cache, storage, email) can be individually enabled/disabled.
-   **Structured Logging**: Implemented with `logrus` for clear and customizable logging.
-   **Configuration Management**: Centralized configuration using Viper, supporting YAML files and environment variables.
-   **Unit of Work Pattern**: Ensures atomic database operations.
-   **Graceful Shutdown**: Handles application shutdown cleanly for all running services.
-   **Docker Support**: Ready-to-use Dockerfile for containerization.
-   **Project Structure**: Clear and maintainable directory layout.


## Delivery Layer Selection

This boilerplate supports **three delivery mechanisms**: REST, gRPC, and GraphQL. 

### Compile-Time Selection (3 Separate Binaries)

Choose your delivery layer at **build time** by selecting the appropriate binary:

| Binary | Delivery Layer | Build Command | Use Case |
|--------|----------------|---------------|----------|
| `rest` | REST/HTTP API | `make build-rest` | RESTful services, web APIs |
| `graphql` | GraphQL API | `make build-graphql` | Flexible queries, frontend-driven APIs |
| `grpc` | gRPC API | `make build-grpc` | High-performance, microservices |

**Benefits:**
- ✅ **True compile-time exclusion** - only your chosen delivery layer is compiled
- ✅ **Smaller binaries** - no unused code in your binary
- ✅ **Clear dependencies** - REST doesn't need proto, gRPC doesn't need Swagger
- ✅ **Clean IDE experience** - no build tag issues with gopls/VS Code

### Quick Start

**Build and run REST API:**
```bash
make build-rest
./bin/go-backend-boilerplate-rest
```

**Build and run GraphQL API:**
```bash
make build-graphql
./bin/go-backend-boilerplate-graphql
```

**Build and run gRPC API:**
```bash
# Requires protoc to be installed
make build-grpc
./bin/go-backend-boilerplate-grpc
```

**Build all three:**
```bash
make build  # Creates all three binaries
```

### Runtime Configuration

Within each binary, you can still use `config.yaml` to control server behavior:

```yaml
server:
  http_server:
    enable: true    # For REST binary
    port: "8080"
  grpc_server:
    enable: true    # For gRPC binary
    port: "50051"
  graphql_server:
    enable: true    # For GraphQL binary
    port: "8081"
```

**Note:** Each binary only respects its relevant server configuration. The REST binary ignores `grpc_server` and `graphql_server` settings.



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

## Optional Infrastructure Configuration

All infrastructure components (database, cache, storage, email) are **optional** and can be individually enabled or disabled via configuration. This allows you to:

- **Start quickly** without configuring unused services
- **Run minimal setups** for development or testing
- **Use only what you need** in production

### Enabling/Disabling Infrastructure

Each infrastructure component has an `enable` flag in the configuration:

```yaml
database:
  enable: true  # Set to false to disable SQL database
  driver: postgres
  host: localhost
  port: 5432
  # ... other database config

mongo:
  enable: false  # Set to true to enable MongoDB
  uri: mongodb://localhost:27017
  database: myapp_db
  # ... other MongoDB config

cache:
  enable: true  # Set to false to disable cache
  type: ristretto
  # ... cache config

storage:
  enable: true  # Set to false to disable storage
  type: local
  # ... storage config

email:
  enable: false  # Set to false to disable email
  type: smtp
  # ... email config
```

### Example Configurations

**Minimal Setup (No Infrastructure)**:
```yaml
database:
  enable: false
mongo:
  enable: false
cache:
  enable: false
storage:
  enable: false
email:
  enable: false
```

**SQL Database Only**:
```yaml
database:
  enable: true
  driver: postgres
  # ... postgres config
mongo:
  enable: false
cache:
  enable: false
storage:
  enable: false
email:
  enable: false
```

**MongoDB Only**:
```yaml
database:
  enable: false
mongo:
  enable: true
  uri: mongodb://localhost:27017
  database: myapp_db
cache:
  enable: false
storage:
  enable: false
email:
  enable: false
```

**Polyglot Persistence (SQL + MongoDB)**:
```yaml
database:
  enable: true
  driver: postgres
  host: localhost
  port: 5432
  # ... postgres config for transactional data

mongo:
  enable: true
  uri: mongodb://localhost:27017
  database: analytics_db
  # ... MongoDB for logs/events/analytics

cache:
  enable: true
  type: redis
  # ... redis cache

storage:
  enable: true
  type: s3
  # ... s3 storage

email:
  enable: true
  type: smtp
  # ... smtp email
```

### MongoDB Configuration

MongoDB support is built-in and can be enabled alongside SQL databases or used independently:

```yaml
mongo:
  enable: true
  uri: mongodb://localhost:27017  # MongoDB connection URI
  database: myapp_db              # Database name
  username: ""                    # Optional: MongoDB username
  password: ""                    # Optional: MongoDB password
  auth_source: admin              # Optional: Authentication database
  connect_timeout_ms: 10000       # Connection timeout in milliseconds
  max_pool_size: 100              # Maximum connection pool size
  min_pool_size: 10               # Minimum connection pool size
```

**Using MongoDB in Your Code**:

The MongoDB client is available in the module dependencies:

```go
// In your module config
func NewMyModule(deps *module.Dependencies) *MyModule {
    if deps.MongoClient != nil {
        // MongoDB is enabled, use it
        collection := deps.MongoClient.Collection("my_collection")
        // ... perform MongoDB operations
    }
    // ... rest of your code
}
```

### Behavior When Infrastructure is Disabled

When an infrastructure component is disabled:
- ✅ Application starts successfully without it
- ✅ No connection attempts are made
- ✅ Clear logs indicate which components are disabled
- ✅ Modules requiring that infrastructure are automatically skipped

**Example Log Output (Minimal Setup)**:
```
INFO[2024-01-16T15:26:59Z] SQL Database is disabled, skipping initialization 
INFO[2024-01-16T15:26:59Z] MongoDB is disabled, skipping initialization 
INFO[2024-01-16T15:26:59Z] Cache is disabled, skipping initialization   
INFO[2024-01-16T15:26:59Z] Storage is disabled, skipping initialization 
INFO[2024-01-16T15:26:59Z] Email client is disabled, skipping initialization 
INFO[2024-01-16T15:26:59Z] 🚀 Starting HTTP server on :8080
```

### 3. Database Setup

#### SQL Databases (PostgreSQL/MySQL/SQLite)

For development, you can use Docker Compose to spin up a PostgreSQL, MySQL, or Redis instance.

```bash
# To start PostgreSQL and Redis:
docker-compose up -d postgres redis

# Or for MySQL and Redis:
# docker-compose up -d mysql redis
```

Update your `config.yaml` or environment variables with the correct database and Redis connection details. If using `sqlite3`, no Docker setup is strictly needed as it uses a local file.

#### MongoDB

To use MongoDB, you can either:

1. **Run MongoDB locally**:
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or add to docker-compose.yaml:
# mongodb:
#   image: mongo:latest
#   ports:
#     - "27017:27017"
#   environment:
#     MONGO_INITDB_ROOT_USERNAME: admin
#     MONGO_INITDB_ROOT_PASSWORD: password
```

2. **Use a cloud MongoDB service** (MongoDB Atlas, etc.) and configure the URI in your `config.yaml`

3. **Disable MongoDB** if not needed by setting `mongo.enable: false`

Update your configuration to enable MongoDB:
```yaml
mongo:
  enable: true
  uri: mongodb://localhost:27017
  database: myapp_db
```

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
│   │   ├── database.go        # SQL database setup (Ent)
│   │   └── mongo/             # MongoDB integration
│   │       ├── client.go      # MongoDB client
│   │       └── health.go      # MongoDB health checks
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
