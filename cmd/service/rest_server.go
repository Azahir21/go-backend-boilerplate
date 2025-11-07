package service

import (
	"fmt"
	"net/http"
	"time"

	restDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRestServer(log *logrus.Logger, cfg config.HTTPServerConfig, userUsecase userUsecase.UserUsecase) *http.Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	if cfg.StartupBanner {
		fmt.Println("🚀 Starting HTTP server...")
	}

	// Initialize handlers
	userHandler := restDelivery.NewUserHandler(log, userUsecase)

	// Initialize Gin server and register routes
	router := sharedHttp.NewServer(userHandler)

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	readTimeout, err := time.ParseDuration(cfg.ReadTimeout)
	if err != nil {
		log.Fatalf("invalid read timeout duration: %v", err)
	}

	writeTimeout, err := time.ParseDuration(cfg.WriteTimeout)
	if err != nil {
		log.Fatalf("invalid write timeout duration: %v", err)
	}

	idleTimeout, err := time.ParseDuration(cfg.IdleTimeout)
	if err != nil {
		log.Fatalf("invalid idle timeout duration: %v", err)
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return server
}
