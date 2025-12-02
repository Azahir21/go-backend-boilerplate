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
	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/cache"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/db"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/external"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/storage"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/azahir21/go-backend-boilerplate/pkg/observability"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Application holds all application-wide dependencies.
type Application struct {
	Log            *logrus.Logger
	Config         *config.Config
	DBClient       *ent.Client
	Cache          cache.Cache
	Storage        storage.Storage
	EmailClient    external.EmailClient
	Dependencies   *module.Dependencies
	HTTPModules    []module.HTTPModule
	GRPCModules    []module.GRPCModule
	GraphQLModules []module.GraphQLModule
	HTTPServer     *http.Server
	GRPCServer     *grpc.Server
	GraphQLServer  *http.Server
	OtelProvider   *observability.Provider
	Metrics        *observability.Metrics
	DBStatsCollector *observability.DBStatsCollector
}

// NewApplication initializes and returns a new Application instance.
func NewApplication(log *logrus.Logger) (*Application, error) {
	// Load configuration
	cfg, err := config.LoadConfig(log)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Env == "production" || cfg.Server.Env == "staging" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize JWT helper
	helper.InitJWT(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// Initialize database
	dbClient, err := db.NewEntClient(log, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize cache
	appCache, err := cache.NewCache(log, cfg.Cache)
	if err != nil {
		dbClient.Close()
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	// Initialize storage
	appStorage, err := storage.NewStorage(context.Background(), log, cfg.Storage)
	if err != nil {
		dbClient.Close()
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize email client
	emailClient, err := external.NewEmailClient(log, cfg.Email)
	if err != nil {
		dbClient.Close()
		return nil, fmt.Errorf("failed to initialize email client: %w", err)
	}

	// Initialize OpenTelemetry observability
	ctx := context.Background()
	otelProvider, err := observability.NewProvider(ctx, observability.Config{
		ServiceName:    "go-backend-boilerplate",
		ServiceVersion: "1.0.0",
		Environment:    cfg.Server.Env,
		TempoEndpoint:  getEnvOrDefault("TEMPO_ENDPOINT", "localhost:4318"),
		EnableTracing:  true,
		EnableMetrics:  true,
	}, log)
	if err != nil {
		log.Warnf("Failed to initialize observability: %v. Continuing without full observability.", err)
	}

	// Initialize custom metrics
	metrics, err := observability.NewMetrics(log)
	if err != nil {
		log.Warnf("Failed to initialize custom metrics: %v. Continuing without custom metrics.", err)
	}

	// Initialize database stats collector
	var dbStatsCollector *observability.DBStatsCollector
	if metrics != nil {
		dbStatsCollector = observability.NewDBStatsCollector(dbClient, metrics, log)
		dbStatsCollector.Start()
		log.Info("Database stats collector started")
	}

	// Initialize unit of work
	uow := unitofwork.NewUnitOfWork(dbClient)

	// Create shared dependencies for all modules
	deps := &module.Dependencies{
		Log:         log,
		DBClient:    dbClient,
		Cache:       appCache,
		Storage:     appStorage,
		EmailClient: emailClient,
		UoW:         uow,
	}

	// Register modules - add new modules here
	httpModules, grpcModules, graphqlModules := registerModules(deps)

	return &Application{
		Log:              log,
		Config:           cfg,
		DBClient:         dbClient,
		Cache:            appCache,
		Storage:          appStorage,
		EmailClient:      emailClient,
		Dependencies:     deps,
		HTTPModules:      httpModules,
		GRPCModules:      grpcModules,
		GraphQLModules:   graphqlModules,
		OtelProvider:     otelProvider,
		Metrics:          metrics,
		DBStatsCollector: dbStatsCollector,
	}, nil
}

func Run(log *logrus.Logger) error {
	app, err := NewApplication(log)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer app.DBClient.Close() // Ensure DB client is closed

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup servers
	if err := app.setupServers(); err != nil {
		return fmt.Errorf("failed to setup servers: %w", err)
	}

	// Start servers
	app.startServers(ctx)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	app.Log.Info("Shutting down servers...")
	if app.HTTPServer != nil {
		app.HTTPServer.Shutdown(ctx)
	}
	if app.GRPCServer != nil {
		app.GRPCServer.GracefulStop()
	}
	if app.GraphQLServer != nil {
		app.GraphQLServer.Shutdown(ctx)
	}
	
	// Stop database stats collector
	if app.DBStatsCollector != nil {
		app.DBStatsCollector.Stop()
		app.Log.Info("Database stats collector stopped")
	}
	
	// Shutdown observability
	if app.OtelProvider != nil {
		if err := app.OtelProvider.Shutdown(ctx); err != nil {
			app.Log.Errorf("Error shutting down observability: %v", err)
		}
	}
	
	app.Log.Info("Servers stopped.")
	return nil
}

// setupServers creates HTTP, gRPC, and GraphQL servers based on configuration.
func (app *Application) setupServers() error {
	if app.Config.Server.HTTP.Enable {
		srv, err := service.NewRestServer(app.Log, app.Config.Server.HTTP, app.HTTPModules)
		if err != nil {
			return fmt.Errorf("failed to create HTTP server: %w", err)
		}
		app.HTTPServer = srv
	}

	if app.Config.Server.GRPC.Enable {
		grpcServer, err := service.NewGrpcServer(app.Log, app.Config.Server.GRPC, app.GRPCModules)
		if err != nil {
			return fmt.Errorf("failed to create gRPC server: %w", err)
		}
		app.GRPCServer = grpcServer
	}

	if app.Config.Server.GraphQL.Enable {
		srv, err := service.NewGraphQLServer(app.Log, app.Config.Server.GraphQL, app.GraphQLModules)
		if err != nil {
			return fmt.Errorf("failed to create GraphQL server: %w", err)
		}
		app.GraphQLServer = srv
	}

	if !app.Config.Server.HTTP.Enable && !app.Config.Server.GRPC.Enable && !app.Config.Server.GraphQL.Enable {
		return errors.New("no server enabled. Please enable at least one of HTTP, gRPC or GraphQL")
	}
	return nil
}

// startServers starts the enabled servers in goroutines.
func (app *Application) startServers(ctx context.Context) {
	if app.Config.Server.HTTP.Enable {
		httpLis, err := net.Listen("tcp", ":"+app.Config.Server.HTTP.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on HTTP port %s: %v", app.Config.Server.HTTP.Port, err)
		} else {
			go func() {
				app.Log.Infof("HTTP server starting on :%s", app.Config.Server.HTTP.Port)
				if err := app.HTTPServer.Serve(httpLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start HTTP server: %v", err)
				}
			}()
		}
	}

	if app.Config.Server.GRPC.Enable {
		grpcLis, err := net.Listen("tcp", ":"+app.Config.Server.GRPC.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on gRPC port %s: %v", app.Config.Server.GRPC.Port, err)
		} else {
			go func() {
				app.Log.Infof("gRPC server starting on :%s", app.Config.Server.GRPC.Port)
				if err := app.GRPCServer.Serve(grpcLis); err != nil && err != grpc.ErrServerStopped {
					app.Log.Errorf("failed to start gRPC server: %v", err)
				}
			}()
		}
	}

	if app.Config.Server.GraphQL.Enable {
		graphqlLis, err := net.Listen("tcp", ":"+app.Config.Server.GraphQL.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on GraphQL port %s: %v", app.Config.Server.GraphQL.Port, err)
		} else {
			go func() {
				app.Log.Infof("GraphQL server starting on :%s", app.Config.Server.GraphQL.Port)
				if err := app.GraphQLServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start GraphQL server: %v", err)
				}
			}()
		}
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
