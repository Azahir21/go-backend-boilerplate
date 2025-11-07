package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/cache"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/db"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/external"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/storage"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Run(log *logrus.Logger) error {
	// Load configuration
	cfg, err := config.LoadConfig(log)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
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
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer client.Close()

	// Initialize cache
	_, err = cache.NewCache(log, cfg.Cache)
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize storage
	_, err = storage.NewStorage(ctx, log, cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize email client
	_, err = external.NewEmailClient(log, cfg.Email)
	if err != nil {
		return fmt.Errorf("failed to initialize email client: %w", err)
	}

	// Initialize repositories
	userRepo := userRepoImpl.NewUserRepository(client)

	// Initialize usecases
	userUsecase := userUsecase.NewUserUsecase(userRepo)

	// Create HTTP, gRPC and GraphQL servers
	var httpServer *http.Server
	var grpcServer *grpc.Server
	var graphqlServer *http.Server

	if cfg.Server.HTTP.Enable {
		httpServer = service.NewRestServer(log, cfg.Server.HTTP, userUsecase)
	}

	if cfg.Server.GRPC.Enable {
		grpcServer = service.NewGrpcServer(log, cfg.Server.GRPC, userUsecase)
	}

	if cfg.Server.GraphQL.Enable {
		graphqlServer = service.NewGraphQLServer(log, cfg.Server.GraphQL, userUsecase)
	}

	// Start servers individually
	if cfg.Server.HTTP.Enable {
		httpLis, err := net.Listen("tcp", ":"+cfg.Server.HTTP.Port)
		if err != nil {
			return fmt.Errorf("failed to listen on HTTP port %s: %w", cfg.Server.HTTP.Port, err)
		}
		go func() {
			log.Infof("HTTP server starting on :%s", cfg.Server.HTTP.Port)
			if err := httpServer.Serve(httpLis); err != nil && err != http.ErrServerClosed {
				log.Errorf("failed to start HTTP server: %v", err)
			}
		}()
	}

	if cfg.Server.GRPC.Enable {
		grpcLis, err := net.Listen("tcp", ":"+cfg.Server.GRPC.Port)
		if err != nil {
			return fmt.Errorf("failed to listen on gRPC port %s: %w", cfg.Server.GRPC.Port, err)
		}
		go func() {
			log.Infof("gRPC server starting on :%s", cfg.Server.GRPC.Port)
			if err := grpcServer.Serve(grpcLis); err != nil && err != grpc.ErrServerStopped {
				log.Errorf("failed to start gRPC server: %v", err)
			}
		}()
	}

	if cfg.Server.GraphQL.Enable {
		graphqlLis, err := net.Listen("tcp", ":"+cfg.Server.GraphQL.Port)
		if err != nil {
			return fmt.Errorf("failed to listen on GraphQL port %s: %w", cfg.Server.GraphQL.Port, err)
		}
		go func() {
			log.Infof("GraphQL server starting on :%s", cfg.Server.GraphQL.Port)
			if err := graphqlServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
				log.Errorf("failed to start GraphQL server: %v", err)
			}
		}()
	}

	if !cfg.Server.HTTP.Enable && !cfg.Server.GRPC.Enable && !cfg.Server.GraphQL.Enable {
		return errors.New("no server enabled. Please enable at least one of HTTP, gRPC or GraphQL")
	}

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
	if graphqlServer != nil {
		graphqlServer.Shutdown(ctx)
	}
	log.Info("Servers stopped.")
	return nil
}
