package service

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/azahir21/go-backend-boilerplate/docs" // Import generated docs
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
	restDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// parseDuration is a helper to parse duration strings and return an error.
func parseDuration(durationStr, fieldName string) (time.Duration, error) {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s duration: %w", fieldName, err)
	}
	return d, nil
}

func NewRestServer(log *logrus.Logger, cfg config.HTTPServerConfig, userUsecase userUsecase.UserUsecase) (*http.Server, error) {
	// Gin mode is set in cmd/app/app.go based on environment.

	if cfg.StartupBanner {
		log.Infof("HTTP server starting on :%s", cfg.Port)
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:"+cfg.Port+"/swagger/doc.json")))

	readTimeout, err := parseDuration(cfg.ReadTimeout, "read timeout")
	if err != nil {
		return nil, err
	}

	writeTimeout, err := parseDuration(cfg.WriteTimeout, "write timeout")
	if err != nil {
		return nil, err
	}

	idleTimeout, err := parseDuration(cfg.IdleTimeout, "idle timeout")
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return server, nil
}
