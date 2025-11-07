package main

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/app"
	"github.com/azahir21/go-backend-boilerplate/internal/delivery/http"
	"github.com/azahir21/go-backend-boilerplate/internal/repository/implementation"
	"github.com/azahir21/go-backend-boilerplate/internal/usecase"
	"github.com/azahir21/go-backend-boilerplate/infra/db"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Headcount Checker API
// @version 1.0
// @description API Documentation for Headcount Checker Backend
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger
	log := logger.NewLogger()

	// Load configuration
	cfg, err := config.LoadConfig(log)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	client, err := db.NewEntClient(log, cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer client.Close()

	// Run migrations
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := implementation.NewUserRepository(client)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize handlers
	userHandler := http.NewUserHandler(log, userUsecase)

	// Initialize Gin server
	server := app.NewServer(userHandler)

	// Swagger endpoint
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))


	// Start HTTP server
	log.Info("HTTP server starting on :" + cfg.Server.HTTPPort)
	if err := server.Run(":" + cfg.Server.HTTPPort); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}