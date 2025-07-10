package main

import (
	_ "github.com/azahir21/go-backend-boilerplate/docs"
	"github.com/azahir21/go-backend-boilerplate/internal/config"
	"github.com/azahir21/go-backend-boilerplate/internal/container"
	"github.com/azahir21/go-backend-boilerplate/internal/routers"
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
    cfg := config.LoadConfig(log)
    
    // Set Gin mode based on environment
    if cfg.ServerEnv == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // Initialize database
    db := config.InitDB(log, cfg)
    config.Migrate(db, log, cfg)

    // Initialize dependency container
    container.Initialize(cfg, db, log)
    
    // Initialize Gin router
    r := gin.Default()
    
    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Register routes
    routers.RegisterRoutes(r)

    // Start server
    log.Info("Server starting on :" + cfg.ServerPort)
    r.Run(":" + cfg.ServerPort)
}