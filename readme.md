# Go Backend Boilerplate

A production-ready RESTful API boilerplate using Go, Gin, GORM, JWT authentication, and Swagger documentation.

## Features

- **Clean Architecture**: Repository, Service, Handler pattern
- **Authentication**: JWT-based authentication system
- **Authorization**: Role-based access control (admin/user roles)
- **API Documentation**: Auto-generated Swagger documentation
- **Environment Configuration**: Using .env files for different environments
- **Database**: PostgreSQL with GORM ORM
- **Dependency Injection**: Singleton container pattern
- **Logging**: Structured JSON logging with Logrus
- **Hot Reload**: Development with Air for live reloading
- **Middleware**: Authentication and role-based middleware
- **Response**: Standardized API response format

## Project Structure

```
├── .air.toml               # Air configuration for hot reload
├── .env                    # Environment variables
├── .env.example            # Example environment variables
├── .gitignore              # Git ignore file
├── go.mod                  # Go modules
├── go.sum                  # Go modules checksums
├── main.go                 # Entry point
├── docs/                   # Auto-generated Swagger docs
├── internal/               # Private application code
│   ├── config/             # Configuration
│   ├── container/          # Dependency injection container
│   ├── handler/            # HTTP handlers
│   ├── helper/             # Helper functions
│   ├── middleware/         # HTTP middleware
│   ├── model/              # Data models
│   ├── repository/         # Data access layer
│   ├── routers/            # HTTP routes
│   └── service/            # Business logic
└── pkg/                    # Public libraries
    ├── logger/             # Logging package
    └── response/           # HTTP response formatter
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL
- [Air](https://github.com/cosmtrek/air) for hot reloading (optional)
- [Swag](https://github.com/swaggo/swag) for Swagger documentation

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-backend-boilerplate.git
   cd go-backend-boilerplate
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Install Swag for Swagger documentation:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

4. Install Air for hot reloading (optional):

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

5. Create a PostgreSQL database:

   ```bash
   createdb headcount_checker
   ```

6. Copy the example environment file and update it with your configuration:

   ```bash
   cp .env.example .env
   ```

7. Generate Swagger documentation:

   ```bash
   swag init
   ```

8. Start the server:

   ```bash
   # Using standard Go run
   go run main.go

   # Or using Air for hot reloading
   air
   ```

## Configuration

The application is configured using environment variables. You can set them in the `.env` file or directly in your environment.

### Environment Variables

```
# Database Configuration
DB_HOST=localhost            # Database host
DB_PORT=5432                 # Database port
DB_USER=postgres             # Database user
DB_PASSWORD=123456           # Database password
DB_NAME=headcount_checker    # Database name
DB_SSL_MODE=disable          # SSL mode (disable, require, verify-ca, verify-full)

# JWT Configuration
JWT_SECRET=your-secret-key   # Secret key for JWT signing
JWT_EXPIRY_HOURS=72          # JWT token expiry in hours

# Server Configuration
SERVER_PORT=8080             # Server port
SERVER_ENV=development       # Environment (development, production)

# Default Admin User
DEFAULT_ADMIN_USERNAME=admin             # Default admin username
DEFAULT_ADMIN_EMAIL=admin@example.com    # Default admin email
DEFAULT_ADMIN_PASSWORD=admin123          # Default admin password
```

## API Documentation

The API documentation is automatically generated using Swagger. After starting the server, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Authentication

The boilerplate includes a complete JWT-based authentication system:

1. **Register**: Create a new user account

   ```
   POST /api/v1/auth/register
   ```

2. **Login**: Authenticate and get a JWT token

   ```
   POST /api/v1/auth/login
   ```

3. **Get Profile**: Get the current user's profile (requires authentication)

   ```
   GET /api/v1/auth/profile
   ```

4. **Admin Only**: Example endpoint that requires admin role
   ```
   GET /api/v1/admin/test
   ```

### Authentication Flow

1. Register a new user or use the default admin account
2. Login to get a JWT token
3. Include the token in the Authorization header for protected routes:
   ```
   Authorization: Bearer your.jwt.token
   ```

## Development Workflow

### Using Air for Hot Reload

The project is configured to use Air for hot reloading during development. When you make changes to your code, Air will automatically rebuild and restart the server.

```bash
# Start the server with hot reloading
air
```

### Generating Swagger Documentation

When you modify API endpoints, you need to regenerate the Swagger documentation:

```bash
swag init
```

This is automatically done when using Air due to the pre_cmd configuration in `.air.toml`.

## Directory Structure Explained

### `/internal`

Contains application code that's not meant to be imported by other applications.

- **`/config`**: Application configuration and environment loading
- **`/container`**: Dependency injection container for managing application dependencies
- **`/handler`**: HTTP handlers that receive HTTP requests and return HTTP responses
- **`/helper`**: Helper functions like JWT utilities and password hashing
- **`/middleware`**: HTTP middleware for authentication and authorization
- **`/model`**: Data models representing database entities and API request/response objects
- **`/repository`**: Data access layer for interacting with the database
- **`/routers`**: HTTP route definitions
- **`/service`**: Business logic layer

### `/pkg`

Contains code that can be imported by other applications.

- **`/logger`**: Logging utilities
- **`/response`**: Standardized API response formatter

## License

This project is licensed under the MIT License - see the LICENSE file for details.
