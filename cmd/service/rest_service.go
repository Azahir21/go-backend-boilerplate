package service

import (
	"net/http"
	"time"

	_ "github.com/azahir21/go-backend-boilerplate/docs" // Import generated docs
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRestServer(log *logrus.Logger, cfg config.HTTPServerConfig, modules []module.HTTPModule) (*http.Server, error) {
	// Gin mode is set in cmd/app/app.go based on environment.
	if cfg.StartupBanner {
		log.Infof("HTTP server starting on :%s", cfg.Port)
	}

	// Collect HTTP handlers from all modules
	var httpRouters []sharedHttp.HttpRouter
	for _, m := range modules {
		log.Infof("Registering HTTP routes for module: %s", m.Name())
		httpRouters = append(httpRouters, m.HTTPHandler())
	}

	// Initialize Gin server and register routes
	router := sharedHttp.NewServer(httpRouters...)

	// Add OpenTelemetry middleware for tracing
	router.Use(otelgin.Middleware("go-backend-boilerplate"))

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CorsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
