# Go Backend Boilerplate

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/Azahir21/go-backend-boilerplate)

> 💡 **Have questions about this repo?** Click the badge above to ask the DeepWiki AI assistant anything about this codebase. It has analyzed the entire repository and can answer specific questions about how the code works!

A **beginner-friendly**, production-ready RESTful API boilerplate using Go, Gin, GORM, JWT authentication, and clean architecture patterns.

> 🎯 **Perfect for newcomers** who want to learn modern Go backend development with best practices built-in!

## 🚀 Features

- **🏗️ Clean Architecture**: Easy-to-understand Repository, Service, Handler pattern with dependency injection
- **🔐 JWT Authentication**: Secure login system with role-based access control (admin/user)
- **📊 Database**: PostgreSQL with Ent ORM (powerful Go ORM)
- **📚 Auto Documentation**: Interactive Swagger/OpenAPI documentation that updates automatically
- **🔒 Built-in Security**: Helmet, CORS, rate limiting, and input validation
- **📝 Smart Logging**: Structured JSON logging with Logrus to track what's happening
- **🔧 Type Safety**: Full Go type safety with struct validation
- **🐳 Docker Ready**: Production-ready Docker setup with one-command deployment
- **🔄 Hot Reload**: Development server automatically restarts with Air when you make changes
- **✅ Data Validation**: Request validation using Go validator to ensure data integrity
- **📋 Code Quality**: Go fmt, Go vet, and golangci-lint for clean code

## 🎓 How This Boilerplate Works

### 📊 Data Flow Explained (For Beginners)

When someone makes a request to your API, here's exactly what happens:

```
1. 📨 HTTP Request → 2. 🛡️ Middleware → 3. 🎯 Handler → 4. 🧠 Service → 5. 💾 Repository → 6. 🗄️ Database → 7. 📤 Response
```

**Step-by-Step Breakdown:**

1. **📨 HTTP Request arrives** (e.g., user wants to login)
2. **🛡️ Middleware checks** authentication, validates data, handles CORS
3. **🎯 Handler receives** the request and extracts user data
4. **🧠 Service processes** business logic (e.g., "is this password correct?")
5. **💾 Repository handles** database operations (e.g., "find user in database")
6. **🗄️ Database returns** the data (user information)
7. **📤 Response flows back** through the same chain to the user

### 🏗️ Architecture Made Simple

```
📁 Your API Structure:
├── 🛣️ routers/         → "Which URLs your API responds to"
├── 🎮 handler/         → "Handles HTTP requests and responses"
├── 🧠 service/         → "Your business logic lives here"
├── 💾 repository/      → "Talks to the database"
├── 📋 model/           → "Defines what your data looks like"
├── 🛡️ middleware/      → "Security, validation, and logging"
├── 🔧 helper/          → "Helper functions you'll use everywhere"
└── 🔌 container/       → "Manages dependencies (like a smart organizer)"
```

## 📁 Project Structure

```
go-backend-boilerplate/
├── .env.example                # 🔧 Copy this to .env and add your settings
├── .gitignore                  # 🚫 Tells Git what files to ignore
├── .air.toml                   # 🔄 Hot reload configuration with Air
├── Dockerfile                  # 🐳 Docker configuration for your app
├── docker-compose.yml          # 🐳 One-command setup with Docker
├── Makefile                   # 🚀 Simple commands (make dev, make test)
├── go.mod                     # 📦 Go module dependencies
├── go.sum                     # 🔒 Dependency checksums
├── main.go                    # 🚪 The entry point - starts your API
├── db/migrations/              # 🗄️ Goose migration files
├── ent/                        # 📊 Ent ORM schema definitions and generated code
│   ├── schema/                # 📝 Ent schema definitions
│   └── generate.go            # ⚙️ Ent code generation
├── docs/                      # 📚 Auto-generated Swagger documentation
│   ├── docs.go                # 📄 Swagger docs configuration
│   ├── swagger.json           # 📋 API specification in JSON
│   └── swagger.yaml           # 📋 API specification in YAML
└── internal/                  # 💻 All your private application code
    ├── config/                # ⚙️ Configuration settings
    │   └── config.go          # 📝 Environment variables setup
    ├── delivery/              # 🚚 API delivery layer (HTTP and gRPC handlers)
    │   ├── http/              # 🌐 HTTP handlers (Gin)
    │   └── grpc/              # 📡 gRPC handlers
    ├── domain/                # 📦 Core business entities and interfaces
    ├── helper/                # 🔨 Helper functions
    │   ├── jwt_helper.go      # 🔐 JWT token utilities
    │   └── password_helper.go # 🔒 Password hashing utilities
    ├── middleware/            # 🛡️ Request processing pipeline
    │   ├── auth_middleware.go # 🔐 Check if user is logged in
    │   └── admin_middleware.go # 👑 Check if user is admin
    ├── repository/            # 💾 Database operations
    │   ├── implementation/    # 📝 Repository implementations
    │   └── user_repo.go       # 👤 User repository interface
    └── usecase/               # 🧠 Business logic
        └── user_usecase.go    # 🔐 User-related business rules
└── pkg/                       # 💻 Public libraries (can be imported by other apps)
    ├── config/                # 📝 Configuration loading
    ├── httpresp/              # 📤 HTTP response helpers
    ├── logger/                # 📝 Logging utilities
    └── response/              # 📤 Standardized response formatter
```

## 🚀 Getting Started (Step by Step)

### Prerequisites (What You Need First)

- **Go 1.21+** → [Download here](https://golang.org/dl/) (Go programming language)
- **PostgreSQL 13+** → [Download here](https://postgresql.org) (Database)
- **Air** → `go install github.com/cosmtrek/air@latest` (Hot reload tool)
- **Swag** → `go install github.com/swaggo/swag/cmd/swag@latest` (Swagger docs generator)
- **Goose** → `go install github.com/pressly/goose/v3/cmd/goose@latest` (Database migration tool)
- **Docker & Docker Compose** → [Download here](https://docker.com) (Optional but recommended)

### 🎯 Quick Start (Recommended for Beginners)

**Option 1: Docker (Easiest - Everything is set up for you!)**

```bash
# 1. Get the code
git clone https://github.com/Azahir21/go-backend-boilerplate.git
cd go-backend-boilerplate

# 2. Set up your environment variables
cp .env.example .env
# Edit .env file with your settings (see Configuration section below)

# 3. Start everything with one command!
make docker-up

# 4. Open your browser and go to:
# http://localhost:8080/swagger/index.html (to see your API documentation)
# http://localhost:8080/ping (to test if it works)
```

**Option 2: Local Development (More control)**

```bash
# 1. Get the code
git clone https://github.com/Azahir21/go-backend-boilerplate.git
cd go-backend-boilerplate

# 2. Install all dependencies
go mod tidy

# 3. Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest

# 4. Set up your environment variables
cp .env.example .env
# Edit .env file with your database details

# 5. Create PostgreSQL database
createdb headcount_checker  # or use your preferred method

# 6. Generate Ent code
go generate ./ent

# 7. Create initial Goose migration file
make goose-create name=initial_schema
# Manually add the SQL for your schema (e.g., CREATE TABLE users...) to the generated file.

# 8. Apply migrations
make goose-up

# 9. Generate Swagger documentation
swag init

# 10. Start your API server with hot reload
air

# 11. Test it works:
# Open http://localhost:8080/swagger/index.html in your browser
```

## ⚙️ Configuration (Environment Variables)

Create a `.env` file and customize these settings:

```bash
# 🗄️ Database Settings (where your data is stored)
DB_HOST=localhost                    # Database server location
DB_PORT=5432                         # Database port
DB_USER=postgres                     # Database username
DB_PASSWORD=postgres                 # Database password
DB_NAME=headcount_checker            # Your database name
DB_SSL_MODE=disable                  # SSL mode (disable, require, verify-ca, verify-full)
AUTO_MIGRATE=false                   # Enable/disable Ent's automatic schema migration

# 🔐 Security Settings (keep these secret!)
JWT_SECRET=change-this-to-something-very-secret  # Used to encrypt tokens
JWT_EXPIRY_HOURS=72                 # How long login tokens last

# 🚀 Server Settings
SERVER_PORT=8080                    # What port your API runs on
SERVER_ENV=development              # development or production
ENABLE_HTTP=true                    # Enable/disable HTTP server
ENABLE_GRPC=true                    # Enable/disable gRPC server

# 👑 Default Admin User (created automatically)
DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_PASSWORD=admin123
```

### 🗄️ Database Setup Made Simple

This boilerplate uses **Ent** - a powerful ORM library for Go that makes database operations easy.

The database models are defined as Go structs in `ent/schema`:

```go
// ent/schema/user.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique(),
		field.String("email").
			Unique(),
		field.String("password"),
		field.String("role").
			Default("user"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
```

**Goose Migrations:**

This project uses [Goose](https://github.com/pressly/goose) for database migrations. Migrations are SQL files located in the `db/migrations` directory.

1.  **Define your Ent schema** in `ent/schema/`.
2.  **Generate Ent code**: `go generate ./ent`
3.  **Create a new Goose migration file**: `make goose-create name=<migration_name>`
4.  **Manually add SQL DDL** (Data Definition Language) statements to the `up` section of the generated migration file, based on your Ent schema. You can generate initial SQL from Ent by running a temporary Go program or by inspecting your Ent schema.
5.  **Apply migrations**: `make goose-up`

## 📚 API Documentation & Testing

### 🌐 Interactive Documentation

Once your server is running, visit these URLs:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
  - 🎮 Interactive API playground - test endpoints directly!
- **Health Check**: `http://localhost:8080/ping`
  - 🔍 Quick test to see if your server is running

### 🧪 Quick Test - Is It Working?

```bash
# Test the health check
curl http://localhost:8080/ping

# Expected response:
{"status": "success", "message": "Server is running", "timestamp": "2024-01-15T10:30:00Z"}
```

## 🔐 Authentication System (How Login Works)

### 📋 Available Endpoints

| Method | Endpoint                | What It Does            | Need Login? | Need Admin? |
| ------ | ----------------------- | ----------------------- | ----------- | ----------- |
| POST   | `/api/v1/auth/register` | Create new user account | ❌          | ❌          |
| POST   | `/api/v1/auth/login`    | Login and get token     | ❌          | ❌          |
| GET    | `/api/v1/auth/profile`  | Get your user info      | ✅          | ❌          |
| GET    | `/api/v1/admin/test`    | Admin-only test         | ✅          | ✅          |

### 🎯 How to Use the Authentication System

#### 1. 📝 Register a New User

```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**What happens:**

- Password gets encrypted and stored safely using bcrypt
- User account is created with "user" role
- You get a success message

#### 2. 🔑 Login to Get Your Token

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    }
}  }
```

**💡 Save that `token` - you'll need it for protected endpoints!**

#### 3. 🔒 Access Protected Endpoints

```bash
curl -X GET "http://localhost:8080/api/v1/auth/profile" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 👑 Default Admin Account

For testing, there's a pre-created admin account:

- **Username:** `admin`
- **Password:** `admin123`
- **Email:** `admin@example.com`

## 🛠️ Development Commands (Your Toolkit)

### 📜 Available Commands

```bash
# 🚀 Running the server
make dev                # Start development server with hot reload (using Air)
make run                # Start production server
make build              # Build the application binary

# 📦 Dependencies & Setup
go mod tidy             # Install/update all dependencies
go mod download         # Download dependencies
swag init               # Generate Swagger documentation

# 🗄️ Database operations
make goose-create name=<migration_name> # Create a new Goose migration file
make goose-up                           # Apply pending Goose migrations
make goose-down                         # Rollback the last Goose migration
make goose-status                       # Check Goose migration status

# 🧪 Testing & Quality
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make lint               # Check code quality with golangci-lint
make fmt                # Format code with go fmt
make vet                # Analyze code with go vet

# 🐳 Docker commands
make docker-build       # Build Docker image
make docker-up          # Start everything with Docker
make docker-down        # Stop Docker containers
make docker-logs        # See what's happening in containers
make docker-shell       # Access container terminal
make docker-clean       # Remove Docker containers and images
```

## 🔄 Working With the Database

### 🗃️ Understanding Models (Database Tables)

This boilerplate uses **Ent** - a powerful ORM library for Go. Your database models are defined as Go structs in `ent/schema`:

```go
// ent/schema/user.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique(),
		field.String("email").
			Unique(),
		field.String("password"),
		field.String("role").
			Default("user"),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).UpdateDefault(time.Now),
		field.Time("deleted_at").
			Optional().Nillable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
```

**Goose Migrations:**

This project uses [Goose](https://github.com/pressly/goose) for database migrations. Migrations are SQL files located in the `db/migrations` directory.

**Workflow for Database Changes:**

1.  **Modify your Ent schema** in `ent/schema/` (e.g., add a new field to `user.go`).
2.  **Generate Ent code**: `go generate ./ent`
3.  **Create a new Goose migration file**: `make goose-create name=<descriptive_migration_name>`
4.  **Manually write SQL DDL** (Data Definition Language) statements in the `up` section of the newly generated migration file to reflect your Ent schema changes. For the `down` section, write the SQL to revert those changes.
5.  **Apply migrations**: `make goose-up`

## 📈 Adding New Features (Step-by-Step Guide)

Let's say you want to add a "Posts" feature where users can create blog posts:

### 1. 🗃️ Create the Ent Schema (Database Table)

```go
// ent/schema/post.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.Text("content").NotEmpty(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

func (Post) Edges() []ent.Edge {
	return nil
}
```

### 2. 💾 Generate Ent Code

```bash
go generate ./ent
```

### 3. 💾 Create Repository (Database Operations)

```go
// internal/repository/post_repo.go
package repository

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/internal/domain"
)

type PostRepository interface {
	Create(ctx context.Context, post *domain.Post) error
	FindByID(ctx context.Context, id int) (*domain.Post, error)
}

type postRepository struct {
	client *ent.Client
}

func NewPostRepository(client *ent.Client) PostRepository {
	return &postRepository{client: client}
}

func (r *postRepository) Create(ctx context.Context, post *domain.Post) error {
	_, err := r.client.Post.Create().
		SetTitle(post.Title).
		SetContent(post.Content).
		Save(ctx)
	return err
}

func (r *postRepository) FindByID(ctx context.Context, id int) (*domain.Post, error) {
	entPost, err := r.client.Post.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainPost(entPost), nil
}

func toDomainPost(entPost *ent.Post) *domain.Post {
	return &domain.Post{
		ID:        entPost.ID,
		Title:     entPost.Title,
		Content:   entPost.Content,
		CreatedAt: entPost.CreatedAt,
		UpdatedAt: entPost.UpdatedAt,
	}
}
```

### 4. 🧠 Create Usecase (Business Logic)

```go
// internal/usecase/post_usecase.go
package usecase

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/domain"
	"github.com/azahir21/go-backend-boilerplate/internal/repository"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, req *domain.CreatePostRequest) (*domain.Post, error)
	GetPost(ctx context.Context, id int) (*domain.Post, error)
}

type postUsecase struct {
	postRepo repository.PostRepository
}

func NewPostUsecase(postRepo repository.PostRepository) PostUsecase {
	return &postUsecase{postRepo: postRepo}
}

func (u *postUsecase) CreatePost(ctx context.Context, req *domain.CreatePostRequest) (*domain.Post, error) {
	post := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
	}

	if err := u.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (u *postUsecase) GetPost(ctx context.Context, id int) (*domain.Post, error) {
	return u.postRepo.FindByID(ctx, id)
}
```

### 5. 🎮 Create Handler (Handle HTTP Requests)

```go
// internal/delivery/http/post_handler.go
package http

import (
	"net/http"
	"strconv"

	"github.com/azahir21/go-backend-boilerplate/internal/domain"
	"github.com/azahir21/go-backend-boilerplate/internal/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	Log         *logrus.Logger
	PostUsecase usecase.PostUsecase
}

func NewPostHandler(log *logrus.Logger, postUsecase usecase.PostUsecase) *PostHandler {
	return &PostHandler{
		Log:         log,
		PostUsecase: postUsecase,
	}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param request body domain.CreatePostRequest true "Post data"
// @Success 201 {object} httpresp.Response{data=domain.Post}
// @Failure 400 {object} httpresp.Response
// @Failure 500 {object} httpresp.Response
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req domain.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpresp.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	post, err := h.PostUsecase.CreatePost(c.Request.Context(), &req)
	if err != nil {
		httpresp.JSON(c, http.StatusInternalServerError, "Failed to create post", err.Error())
		return
	}

	httpresp.JSON(c, http.StatusCreated, "Post created successfully", post)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a single blog post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} httpresp.Response{data=domain.Post}
// @Failure 400 {object} httpresp.Response
// @Failure 404 {object} httpresp.Response
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpresp.JSON(c, http.StatusBadRequest, "Invalid post ID", nil)
		return
	}

	post, err := h.PostUsecase.GetPost(c.Request.Context(), id)
	if err != nil {
		httpresp.JSON(c, http.StatusNotFound, "Post not found", nil)
		return
	}

	httpresp.JSON(c, http.StatusOK, "Post retrieved successfully", post)
}
```

### 6. 🛣️ Register Routes (API Endpoints)

```go
// internal/app/server.go (add to existing NewServer function)
// ...
// Add post routes
// postHandler := http.NewPostHandler(log, postUsecase)
// server := app.NewServer(userHandler, postHandler)
// ...
```

### 7. 🔌 Register in Dependency Container

```go
// internal/container/container.go (add to existing container)
// ...
// func (c *Container) GetPostRepository() repository.PostRepository {
// 	if c.postRepo == nil {
// 		c.postRepo = implementation.NewPostRepository(c.GetDB())
// 	}
// 	return c.postRepo
// }
//
// func (c *Container) GetPostUsecase() usecase.PostUsecase {
// 	if c.postUsecase == nil {
// 		c.postUsecase = usecase.NewPostUsecase(c.GetPostRepository())
// 	}
// 	return c.postUsecase
// }
//
// func (c *Container) GetPostHandler() *httpDelivery.PostHandler {
// 	if c.postHandler == nil {
// 		c.postHandler = httpDelivery.NewPostHandler(c.GetLogger(), c.GetPostUsecase())
// 	}
// 	return c.postHandler
// }
// ...
```

### 8. 🔄 Run Goose Migration

```bash
# Create a new migration file
make goose-create name=create_posts_table

# Manually add the SQL for creating the posts table to the generated file.
# Example SQL (based on Ent schema):
# -- +goose Up
# -- SQL in section 'Up' is executed when this migration is applied
# CREATE TABLE posts (
#     id SERIAL PRIMARY KEY,
#     title VARCHAR(255) NOT NULL,
#     content TEXT NOT NULL,
#     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
#     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
# );
#
# -- +goose Down
# -- SQL in section 'Down' is executed when this migration is rolled back
# DROP TABLE posts;

# Apply migrations
make goose-up
```

Now you have a complete Posts feature! 🎉

## 🧪 Testing Your API

### 🎯 Manual Testing with Swagger

Once your server is running, visit these URLs:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
  - 🎮 Interactive API playground - test endpoints directly!
- **Health Check**: `http://localhost:8080/ping`
  - 🔍 Quick test to see if your server is running

### 🤖 Automated Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test file
go test ./internal/service -v

# Run tests with verbose output
go test -v ./...
```

**Example test:**

```go
// internal/service/auth_service_test.go
func TestAuthService_Register(t *testing.T) {
    // Setup test database and repository
    db := setupTestDB()
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo)

    req := &domain.RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "testpass123",
    }

    authResponse, err := authService.Register(context.Background(), req)

    assert.NoError(t, err)
    assert.NotNil(t, authResponse)
    assert.Equal(t, "testuser", authResponse.User.Username)
    assert.Equal(t, "test@example.com", authResponse.User.Email)
}
```

## 🚀 Deployment (Going Live)

### 🐳 Docker Deployment (Recommended)

```bash
# 1. Build your application
make docker-build

# 2. Start in production mode
make docker-up

# 3. Your API is now live at http://your-server:8080
```

### 🖥️ Manual Deployment

```bash
# 1. Install dependencies
go mod download

# 2. Build the application
make build

# 3. Set up production environment
cp .env.example .env
# Edit .env with production settings

# 4. Start production server
./main  # or your binary name
```

### 🔧 Production Configuration

```bash
# Production .env settings
SERVER_ENV=production
JWT_SECRET=your-super-secure-jwt-secret-key
DB_SSL_MODE=require
# ... other production settings
```

## 🏗️ Architecture Deep Dive

### 🔄 Request Lifecycle (What Happens When Someone Calls Your API)

```
📨 HTTP Request (e.g., POST /api/v1/auth/login)
     ↓
🛡️ Middleware (CORS, Auth, Validation)
     ↓
🛣️ Router (matches /api/v1/auth/login to login handler)
     ↓
🎮 Handler (extracts username/password, calls usecase)
     ↓
🧠 Usecase (business logic: "check if password is correct")
     ↓
💾 Repository (database query: "find user by username")
     ↓
🗄️ Database (PostgreSQL returns user data)
     ↓
📤 Response (success/error flows back to user)
```

### 🧩 Component Responsibilities

- **🛣️ Routers** → "Which function handles this URL?"
- **🎮 Handlers** → "Extract data from request, call usecase, format response"
- **🧠 Usecases** → "Business rules and logic"
- **💾 Repositories** → "How to get/save data from database"
- **📦 Domains** → "Go structs and data structures for business entities"
- **🛡️ Middleware** → "Security, validation, logging, error handling"
- **🔌 Container** → "Manages dependencies (like a smart organizer)"

### 🔌 Dependency Injection Made Simple

Think of the container as a smart organizer that creates and manages all your app components:

```go
// internal/container/container.go
type Container struct {
    db          *ent.Client
    userRepo    repository.UserRepository
    userUsecase usecase.UserUsecase
    userHandler httpDelivery.UserHandler
    // ... other components
}

func (c *Container) GetUserHandler() *httpDelivery.UserHandler {
    if c.userHandler == nil {
        userRepo := implementation.NewUserRepository(c.GetDB())
        userUsecase := usecase.NewUserUsecase(userRepo)
        c.userHandler = *httpDelivery.NewUserHandler(c.GetLogger(), userUsecase)
    }
    return &c.userHandler
}
```

## 📋 Best Practices & Tips

### ✅ Code Quality Tips

1. **Use Go's type system effectively**

   ```go
   // Good: Type-safe function with clear return types
   func (r *userRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
       entUser, err := r.client.User.Get(ctx, int(id))
       if err != nil {
           return nil, err
       }
       return toDomainUser(entUser), nil
   }
   ```

2. **Handle errors gracefully**

   ```go
   // Good: Proper error handling
   user, err := u.userUsecase.Register(ctx, req)
   if err != nil {
       h.Log.WithError(err).Error("Failed to register user")
       httpresp.JSON(c, http.StatusInternalServerError, "Failed to register user", err.Error())
       return
   }
   httpresp.JSON(c, http.StatusCreated, "User registered successfully", user)
   ```

3. **Use struct validation**
   ```go
   // Good: Validate using struct tags
   type RegisterRequest struct {
       Username string `json:"username" binding:"required,min=3,max=50"`
       Email    string `json:"email" binding:"required,email"`
       Password string `json:"password" binding:"required,min=6"`
   }
   ```

### 🔒 Security Best Practices

1. **Never store plain text passwords** (use bcrypt)
2. **Always validate input data** with struct validation tags
3. **Use environment variables** for secrets
4. **Implement rate limiting** to prevent abuse
5. **Keep JWT secrets secure** and rotate them

### 📈 Performance Tips

1. **Use database indexes** for frequently queried fields
2. **Implement pagination** for large datasets
3. **Use Ent eager loading** to avoid N+1 queries
4. **Use goroutines** for concurrent operations when appropriate
5. **Monitor your API performance** with middleware

## 🆘 Troubleshooting

### ❌ Common Issues & Solutions

**Problem:** `cannot find module` errors
**Solution:** Run `go mod tidy` and make sure you're in the project directory

**Problem:** Database connection error
**Solution:** Check your `.env` file has correct database credentials and PostgreSQL is running

**Problem:** JWT token invalid
**Solution:** Check if `JWT_SECRET` in `.env` matches between token creation and validation

**Problem:** Swagger docs not updating
**Solution:** Run `swag init` to regenerate documentation

### 🔍 Debugging Tips

1. **Check the logs:**

   ```bash
   # Local development
   tail -f logs/app.log

   # Docker
   make docker-logs
   ```

2. **Test individual components:**

   ```bash
   # Test database connection
   go run main.go -test-db

   # Test specific endpoint
   curl -v http://localhost:8080/ping
   ```

3. **Use Go's built-in debugging:**

   ```go
   import "log"

   log.Printf("Debug info: userID=%d, userData=%+v", userID, userData)
   // Or use Delve debugger: dlv debug
   ```

## 🤝 Contributing

Want to improve this boilerplate? Here's how:

1. **Fork the repository** on GitHub
2. **Create a feature branch:** `git checkout -b feature/amazing-feature`
3. **Make your changes** and test them
4. **Commit your changes:** `git commit -m 'Add some amazing feature'`
5. **Push to the branch:** `git push origin feature/amazing-feature`
6. **Open a Pull Request** and describe what you've added

### 🐛 Found a Bug?

1. Check if the issue already exists
2. Create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Your Go version and environment details

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help, please:

1. **Check the API documentation:** `http://localhost:8080/swagger/index.html`
2. **Look through existing issues:** [GitHub Issues](https://github.com/Azahir21/go-backend-boilerplate/issues)
3. **Create a new issue** if needed
4. **Join our community discussions**

## 🙏 Acknowledgments

This boilerplate is built on top of amazing open-source projects:

- **[Gin](https://gin-gonic.com/)** → High-performance HTTP web framework written in Go
- **[Ent](https://entgo.io/)** → An entity framework for Go
- **[Goose](https://github.com/pressly/goose)** → Database migration tool
- **[JWT-Go](https://golang-jwt/jwt)** → Go implementation of JSON Web Tokens
- **[Logrus](https://github.com/sirupsen/logrus)** → Structured logger for Go
- **[Swag](https://github.com/swaggo/swag)** → Automatically generate RESTful API documentation
- **[Air](https://github.com/cosmtrek/air)** → Live reload for Go apps
- **[Validator](https://github.com/go-playground/validator)** → Go Struct and Field validation

---

**🎉 Happy coding! If this boilerplate helped you, consider giving it a star ⭐**

**💬 Questions? Open an issue or discussion - we're here to help!**

**Made with ❤️ by the Go community**

## 🚀 Getting Started (Step by Step)

### Prerequisites (What You Need First)

- **Go 1.21+** → [Download here](https://golang.org/dl/) (Go programming language)
- **PostgreSQL 13+** → [Download here](https://postgresql.org) (Database)
- **Air** → `go install github.com/cosmtrek/air@latest` (Hot reload tool)
- **Swag** → `go install github.com/swaggo/swag/cmd/swag@latest` (Swagger docs generator)
- **Docker & Docker Compose** → [Download here](https://docker.com) (Optional but recommended)

### 🎯 Quick Start (Recommended for Beginners)

**Option 1: Docker (Easiest - Everything is set up for you!)**

```bash
# 1. Get the code
git clone https://github.com/Azahir21/go-backend-boilerplate.git
cd go-backend-boilerplate

# 2. Set up your environment variables
cp .env.example .env
# Edit .env file with your settings (see Configuration section below)

# 3. Start everything with one command!
make docker-up

# 4. Open your browser and go to:
# http://localhost:8080/swagger/index.html (to see your API documentation)
# http://localhost:8080/ping (to test if it works)
```

**Option 2: Local Development (More control)**

```bash
# 1. Get the code
git clone https://github.com/Azahir21/go-backend-boilerplate.git
cd go-backend-boilerplate

# 2. Install all dependencies
go mod tidy

# 3. Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 4. Set up your environment variables
cp .env.example .env
# Edit .env file with your database details

# 5. Create PostgreSQL database
createdb go_boilerplate  # or use your preferred method

# 6. Generate Swagger documentation
swag init

# 7. Start your API server with hot reload
air

# 8. Test it works:
# Open http://localhost:8080/swagger/index.html in your browser
```

## ⚙️ Configuration (Environment Variables)

Create a `.env` file and customize these settings:

```bash
# 🗄️ Database Settings (where your data is stored)
DB_HOST=localhost                    # Database server location
DB_PORT=5432                         # Database port
DB_USER=postgres                     # Database username
DB_PASSWORD=postgres                 # Database password
DB_NAME=go_boilerplate              # Your database name
DB_SSL_MODE=disable                  # SSL mode (disable, require, verify-ca, verify-full)

# 🔐 Security Settings (keep these secret!)
JWT_SECRET=change-this-to-something-very-secret  # Used to encrypt tokens
JWT_EXPIRY_HOURS=72                 # How long login tokens last

# 🚀 Server Settings
SERVER_PORT=8080                    # What port your API runs on
SERVER_ENV=development              # development or production

# 👑 Default Admin User (created automatically)
DEFAULT_ADMIN_USERNAME=admin
DEFAULT_ADMIN_EMAIL=admin@example.com
DEFAULT_ADMIN_PASSWORD=admin123
```

### 🗄️ Database Setup Made Simple

This boilerplate uses **GORM** - a fantastic ORM library for Go that makes database operations easy.

The database models are defined as Go structs:

```go
// internal/model/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`  // "-" hides password in JSON
    Role      string    `json:"role" gorm:"default:user"`  // "user" or "admin"
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

## 📚 API Documentation & Testing

### 🌐 Interactive Documentation

Once your server is running, visit these URLs:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
  - 🎮 Interactive API playground - test endpoints directly!
- **Health Check**: `http://localhost:8080/ping`
  - 🔍 Quick test to see if your server is running

### 🧪 Quick Test - Is It Working?

```bash
# Test the health check
curl http://localhost:8080/ping

# Expected response:
{"status": "success", "message": "Server is running", "timestamp": "2024-01-15T10:30:00Z"}
```

## 🔐 Authentication System (How Login Works)

### 📋 Available Endpoints

| Method | Endpoint                | What It Does            | Need Login? | Need Admin? |
| ------ | ----------------------- | ----------------------- | ----------- | ----------- |
| POST   | `/api/v1/auth/register` | Create new user account | ❌          | ❌          |
| POST   | `/api/v1/auth/login`    | Login and get token     | ❌          | ❌          |
| GET    | `/api/v1/auth/profile`  | Get your user info      | ✅          | ❌          |
| GET    | `/api/v1/admin/test`    | Admin-only test         | ✅          | ✅          |

### 🎯 How to Use the Authentication System

#### 1. 📝 Register a New User

```bash
curl -X POST "http://localhost:8080/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**What happens:**

- Password gets encrypted and stored safely using bcrypt
- User account is created with "user" role
- You get a success message

#### 2. 🔑 Login to Get Your Token

```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

**Response:**

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    }
  }
}
```

**💡 Save that `token` - you'll need it for protected endpoints!**

#### 3. 🔒 Access Protected Endpoints

```bash
curl -X GET "http://localhost:8080/api/v1/auth/profile" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 👑 Default Admin Account

For testing, there's a pre-created admin account:

- **Username:** `admin`
- **Password:** `admin123`
- **Email:** `admin@example.com`

## 🛠️ Development Commands (Your Toolkit)

### 📜 Available Commands

```bash
# 🚀 Running the server
make dev                # Start development server with hot reload (using Air)
make run                # Start production server
make build              # Build the application binary

# 📦 Dependencies & Setup
go mod tidy             # Install/update all dependencies
go mod download         # Download dependencies
swag init               # Generate Swagger documentation

# 🗄️ Database operations
make migrate            # Auto-migrate database schema (GORM)
make seed               # Add sample data (including admin user)

# 🧪 Testing & Quality
make test               # Run all tests
make test-coverage      # Run tests with coverage report
make lint               # Check code quality with golangci-lint
make fmt                # Format code with go fmt
make vet                # Analyze code with go vet

# 🐳 Docker commands
make docker-build       # Build Docker image
make docker-up          # Start everything with Docker
make docker-down        # Stop Docker containers
make docker-logs        # See what's happening in containers
make docker-shell       # Access container terminal
make docker-clean       # Remove Docker containers and images
```

## 🔄 Working With the Database

### 🗃️ Understanding Models (Database Tables)

Models define what your database tables look like using Go structs and GORM tags:

```go
// internal/model/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Username  string    `json:"username" gorm:"unique;not null" validate:"required,min=3,max=20"`
    Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
    Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
    Role      string    `json:"role" gorm:"default:user"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

### 🔧 Simple Migrations (GORM Auto-Migration!)

GORM provides automatic migration - it creates and updates your database schema based on your models:

1. **Define your models** in `internal/model/`
2. **Add migration in your main application:**

```go
// This automatically creates/updates tables based on your structs
func migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Post{}, // Add new models here
        // Add more models as needed
    )
}
```

3. **Run migrations:**

```bash
make migrate
```

## 📈 Adding New Features (Step-by-Step Guide)

Let's say you want to add a "Posts" feature where users can create blog posts:

### 1. 🗃️ Create the Model (Database Table)

```go
// internal/model/post.go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Title     string    `json:"title" gorm:"not null" validate:"required,min=1,max=200"`
    Content   string    `json:"content" gorm:"type:text;not null" validate:"required"`
    AuthorID  uint      `json:"author_id" gorm:"not null"`
    Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Request/Response structs
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=1,max=200"`
    Content string `json:"content" validate:"required"`
}

type PostResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    AuthorID  uint      `json:"author_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. 💾 Create Repository (Database Operations)

```go
// internal/repository/post_repository.go
package repository

import (
    "go-backend-boilerplate/internal/model"
    "gorm.io/gorm"
)

type PostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
    return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *model.Post) error {
    return r.db.Create(post).Error
}

func (r *PostRepository) GetUserPosts(authorID uint) ([]model.Post, error) {
    var posts []model.Post
    err := r.db.Where("author_id = ?", authorID).Order("created_at desc").Find(&posts).Error
    return posts, err
}

func (r *PostRepository) GetPostByID(id uint) (*model.Post, error) {
    var post model.Post
    err := r.db.Preload("Author").First(&post, id).Error
    return &post, err
}
```

### 3. 🧠 Create Service (Business Logic)

```go
// internal/service/post_service.go
package service

import (
    "errors"
    "go-backend-boilerplate/internal/model"
    "go-backend-boilerplate/internal/repository"
)

type PostService struct {
    postRepo *repository.PostRepository
}

func NewPostService(postRepo *repository.PostRepository) *PostService {
    return &PostService{postRepo: postRepo}
}

func (s *PostService) CreatePost(req *model.CreatePostRequest, authorID uint) (*model.PostResponse, error) {
    // Business logic validation
    if len(req.Title) == 0 {
        return nil, errors.New("post title cannot be empty")
    }

    post := &model.Post{
        Title:    req.Title,
        Content:  req.Content,
        AuthorID: authorID,
    }

    if err := s.postRepo.CreatePost(post); err != nil {
        return nil, err
    }

    return &model.PostResponse{
        ID:        post.ID,
        Title:     post.Title,
        Content:   post.Content,
        AuthorID:  post.AuthorID,
        CreatedAt: post.CreatedAt,
        UpdatedAt: post.UpdatedAt,
    }, nil
}

func (s *PostService) GetUserPosts(authorID uint) ([]model.PostResponse, error) {
    posts, err := s.postRepo.GetUserPosts(authorID)
    if err != nil {
        return nil, err
    }

    var responses []model.PostResponse
    for _, post := range posts {
        responses = append(responses, model.PostResponse{
            ID:        post.ID,
            Title:     post.Title,
            Content:   post.Content,
            AuthorID:  post.AuthorID,
            CreatedAt: post.CreatedAt,
            UpdatedAt: post.UpdatedAt,
        })
    }

    return responses, nil
}
```

### 4. 🎮 Create Handler (Handle HTTP Requests)

```go
// internal/handler/post_handler.go
package handler

import (
    "net/http"
    "strconv"

    "go-backend-boilerplate/internal/model"
    "go-backend-boilerplate/internal/service"
    "go-backend-boilerplate/pkg/response"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type PostHandler struct {
    postService *service.PostService
    validator   *validator.Validate
}

func NewPostHandler(postService *service.PostService, validator *validator.Validate) *PostHandler {
    return &PostHandler{
        postService: postService,
        validator:   validator,
    }
}

// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body model.CreatePostRequest true "Post data"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response
// @Router /api/v1/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
    var req model.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid JSON format", err.Error())
        return
    }

    if err := h.validator.Struct(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
        return
    }

    // Get user ID from JWT middleware
    userID, exists := c.Get("user_id")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }

    post, err := h.postService.CreatePost(&req, userID.(uint))
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create post", err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "Post created successfully", post)
}

// @Summary Get user's posts
// @Description Get all posts created by the current user
// @Tags posts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /api/v1/posts/my-posts [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        response.Error(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }

    posts, err := h.postService.GetUserPosts(userID.(uint))
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to get posts", err.Error())
        return
    }

    response.Success(c, http.StatusOK, "Posts retrieved successfully", posts)
}
```

### 5. 🛣️ Create Routes (API Endpoints)

```go
// internal/routers/post.go
package routers

import (
    "go-backend-boilerplate/internal/handler"
    "go-backend-boilerplate/internal/middleware"

    "github.com/gin-gonic/gin"
)

func SetupPostRoutes(r *gin.Engine, postHandler *handler.PostHandler, authMiddleware *middleware.AuthMiddleware) {
    api := r.Group("/api/v1")

    posts := api.Group("/posts")
    posts.Use(authMiddleware.RequireAuth()) // All post routes require authentication
    {
        posts.POST("/", postHandler.CreatePost)
        posts.GET("/my-posts", postHandler.GetUserPosts)
    }
}
```

### 6. 🔌 Register in Dependency Container

```go
// internal/container/container.go (add to existing container)
func (c *Container) GetPostHandler() *handler.PostHandler {
    if c.postHandler == nil {
        postRepo := repository.NewPostRepository(c.GetDB())
        postService := service.NewPostService(postRepo)
        c.postHandler = handler.NewPostHandler(postService, c.GetValidator())
    }
    return c.postHandler
}
```

### 7. 🔌 Register Routes in Main App

```go
// main.go (add to your route setup)
func setupRoutes(container *container.Container) *gin.Engine {
    r := gin.Default()

    // ... existing routes ...

    // Add post routes
    routers.SetupPostRoutes(r, container.GetPostHandler(), container.GetAuthMiddleware())

    return r
}
```

### 8. 🔄 Run Migration

```bash
# Add Post model to your migration and run
make migrate
```

Now you have a complete Posts feature! 🎉

## 🧪 Testing Your API

### 🎯 Manual Testing with Swagger

1. **Start your server:** `make dev` or `air`
2. **Open Swagger:** `http://localhost:8080/swagger/index.html`
3. **Try the endpoints:**
   - Click on an endpoint
   - Click "Try it out"
   - Fill in the data
   - Click "Execute"

### 🤖 Automated Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test file
go test ./internal/service -v

# Run tests with verbose output
go test -v ./...
```

**Example test:**

```go
// internal/service/auth_service_test.go
func TestAuthService_Register(t *testing.T) {
    // Setup test database and repository
    db := setupTestDB()
    userRepo := repository.NewUserRepository(db)
    authService := service.NewAuthService(userRepo)

    req := &model.RegisterRequest{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "testpass123",
    }

    user, err := authService.Register(req)

    assert.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "test@example.com", user.Email)
}
```

## 🚀 Deployment (Going Live)

### 🐳 Docker Deployment (Recommended)

```bash
# 1. Build your application
make docker-build

# 2. Start in production mode
make docker-up

# 3. Your API is now live at http://your-server:8080
```

### 🖥️ Manual Deployment

```bash
# 1. Install dependencies
go mod download

# 2. Build the application
make build

# 3. Set up production environment
cp .env.example .env
# Edit .env with production settings

# 4. Start production server
./main  # or your binary name
```

### 🔧 Production Configuration

```bash
# Production .env settings
SERVER_ENV=production
JWT_SECRET=your-super-secure-jwt-secret-key
DB_SSL_MODE=require
# ... other production settings
```

## 🏗️ Architecture Deep Dive

### 🔄 Request Lifecycle (What Happens When Someone Calls Your API)

```
📨 HTTP Request (e.g., POST /api/v1/auth/login)
     ↓
🛡️ Middleware (CORS, Auth, Validation)
     ↓
🛣️ Router (matches /api/v1/auth/login to login handler)
     ↓
🎮 Handler (extracts username/password, calls service)
     ↓
🧠 Service (business logic: "check if password is correct")
     ↓
💾 Repository (database query: "find user by username")
     ↓
🗄️ Database (PostgreSQL returns user data)
     ↓
📤 Response (success/error flows back to user)
```

### 🧩 Component Responsibilities

- **🛣️ Routers** → "Which function handles this URL?"
- **🎮 Handlers** → "Extract data from request, call service, format response"
- **🧠 Services** → "Business rules and logic"
- **💾 Repositories** → "How to get/save data from database"
- **📋 Models** → "Go structs and data structures"
- **🛡️ Middleware** → "Security, validation, logging, error handling"
- **🔌 Container** → "Manages dependencies (like a smart organizer)"

### 🔌 Dependency Injection Made Simple

Think of the container as a smart organizer that creates and manages all your app components:

```go
// internal/container/container.go
type Container struct {
    db          *gorm.DB
    userRepo    *repository.UserRepository
    authService *service.AuthService
    authHandler *handler.AuthHandler
}

func (c *Container) GetAuthHandler() *handler.AuthHandler {
    if c.authHandler == nil {
        userRepo := c.GetUserRepository()
        authService := service.NewAuthService(userRepo)
        c.authHandler = handler.NewAuthHandler(authService, c.GetValidator())
    }
    return c.authHandler
}
```

## 📋 Best Practices & Tips

### ✅ Code Quality Tips

1. **Use Go's type system effectively**

   ```go
   // Good: Type-safe function with clear return types
   func GetUser(userID uint) (*model.User, error) {
       var user model.User
       err := db.First(&user, userID).Error
       return &user, err
   }
   ```

2. **Handle errors gracefully**

   ```go
   // Good: Proper error handling
   user, err := userService.CreateUser(userData)
   if err != nil {
       logger.WithError(err).Error("Failed to create user")
       response.Error(c, http.StatusInternalServerError, "Failed to create user", err.Error())
       return
   }
   response.Success(c, http.StatusCreated, "User created successfully", user)
   ```

3. **Use struct validation**
   ```go
   // Good: Validate using struct tags
   type RegisterRequest struct {
       Username string `json:"username" validate:"required,min=3,max=20"`
       Email    string `json:"email" validate:"required,email"`
       Password string `json:"password" validate:"required,min=6"`
   }
   ```

### 🔒 Security Best Practices

1. **Never store plain text passwords** (use bcrypt)
2. **Always validate input data** with struct validation tags
3. **Use environment variables** for secrets
4. **Implement rate limiting** to prevent abuse
5. **Keep JWT secrets secure** and rotate them

### 📈 Performance Tips

1. **Use database indexes** for frequently queried fields
2. **Implement pagination** for large datasets
3. **Use GORM preloading** to avoid N+1 queries
4. **Use goroutines** for concurrent operations when appropriate
5. **Monitor your API performance** with middleware

## 🆘 Troubleshooting

### ❌ Common Issues & Solutions

**Problem:** `cannot find module` errors
**Solution:** Run `go mod tidy` and make sure you're in the project directory

**Problem:** Database connection error
**Solution:** Check your `.env` file has correct database credentials and PostgreSQL is running

**Problem:** JWT token invalid
**Solution:** Check if `JWT_SECRET` in `.env` matches between token creation and validation

**Problem:** Swagger docs not updating
**Solution:** Run `swag init` to regenerate documentation

### 🔍 Debugging Tips

1. **Check the logs:**

   ```bash
   # Local development
   tail -f logs/app.log

   # Docker
   make docker-logs
   ```

2. **Test individual components:**

   ```bash
   # Test database connection
   go run main.go -test-db

   # Test specific endpoint
   curl -v http://localhost:8080/ping
   ```

3. **Use Go's built-in debugging:**

   ```go
   import "log"

   log.Printf("Debug info: userID=%d, userData=%+v", userID, userData)
   // Or use Delve debugger: dlv debug
   ```

## 🤝 Contributing

Want to improve this boilerplate? Here's how:

1. **Fork the repository** on GitHub
2. **Create a feature branch:** `git checkout -b feature/amazing-feature`
3. **Make your changes** and test them
4. **Commit your changes:** `git commit -m 'Add some amazing feature'`
5. **Push to the branch:** `git push origin feature/amazing-feature`
6. **Open a Pull Request** and describe what you've added

### 🐛 Found a Bug?

1. Check if the issue already exists
2. Create a new issue with:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - Your Go version and environment details

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help, please:

1. **Check the API documentation:** `http://localhost:8080/swagger/index.html`
2. **Look through existing issues:** [GitHub Issues](https://github.com/Azahir21/go-backend-boilerplate/issues)
3. **Create a new issue** if needed
4. **Join our community discussions**

## 🙏 Acknowledgments

This boilerplate is built on top of amazing open-source projects:

- **[Gin](https://gin-gonic.com/)** → High-performance HTTP web framework written in Go
- **[GORM](https://gorm.io/)** → The fantastic ORM library for Golang
- **[JWT-Go](https://github.com/golang-jwt/jwt)** → Go implementation of JSON Web Tokens
- **[Logrus](https://github.com/sirupsen/logrus)** → Structured logger for Go
- **[Swag](https://github.com/swaggo/swag)** → Automatically generate RESTful API documentation
- **[Air](https://github.com/cosmtrek/air)** → Live reload for Go apps
- **[Validator](https://github.com/go-playground/validator)** → Go Struct and Field validation

---

**🎉 Happy coding! If this boilerplate helped you, consider giving it a star ⭐**

**💬 Questions? Open an issue or discussion - we're here to help!**

**Made with ❤️ by the Go community**
