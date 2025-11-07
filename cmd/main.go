package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/infra/cache"
	"github.com/azahir21/go-backend-boilerplate/infra/db"
	"github.com/azahir21/go-backend-boilerplate/infra/external"
	"github.com/azahir21/go-backend-boilerplate/infra/storage"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	if cfg.Server.Env == "production" || cfg.Server.Env == "staging" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	client, err := db.NewEntClient(log, cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer client.Close()

	// Initialize cache
	_, err = cache.NewCache(log, cfg.Cache)
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}

	// Initialize storage
	_, err = storage.NewStorage(ctx, log, cfg.Storage)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	// Initialize email client
	_, err = external.NewEmailClient(log, cfg.Email)
	if err != nil {
		log.Fatalf("failed to initialize email client: %v", err)
	}

	// Initialize repositories
	userRepo := userRepoImpl.NewUserRepository(client)

	// Initialize usecases
	userUsecase := userUsecase.NewUserUsecase(userRepo)

	// Create HTTP and gRPC servers
	var httpServer *http.Server
	var grpcServer *grpc.Server

	if cfg.Server.HTTP.Enable {
		httpServer = service.NewRestServer(log, cfg.Server.HTTP, userUsecase)
	}

	if cfg.Server.GRPC.Enable {
		grpcServer = service.NewGrpcServer(log, cfg.Server.GRPC, userUsecase)
	}

	// Determine the port to listen on
	var listenPort string
	if cfg.Server.HTTP.Enable {
		listenPort = cfg.Server.HTTP.Port
	} else if cfg.Server.GRPC.Enable {
		listenPort = cfg.Server.GRPC.Port
	} else {
		log.Fatalf("No server enabled. Please enable at least one of HTTP or gRPC.")
	}

	// Create a http.Handler that can multiplex between gRPC and HTTP requests
	var handler http.Handler

	if cfg.Server.HTTP.Enable && cfg.Server.GRPC.Enable {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r)
			} else {
				httpServer.Handler.ServeHTTP(w, r)
			}
		})
	} else if cfg.Server.HTTP.Enable {
		handler = httpServer.Handler
	} else if cfg.Server.GRPC.Enable {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			grpcServer.ServeHTTP(w, r)
		})
	} else {
		log.Fatalf("No server enabled. Please enable at least one of HTTP or gRPC.")
	}

	// Create a listener for the shared port
	lis, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start the main server
	mainServer := &http.Server{
		Addr:    ":" + listenPort,
		Handler: handler,
	}

	go func() {
		log.Infof("Main server starting on :%s", listenPort)
		if err := mainServer.Serve(lis); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start main server: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down servers...")
	if httpServer != nil {
		httpServer.Shutdown(ctx)
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	log.Info("Servers stopped.")
}
