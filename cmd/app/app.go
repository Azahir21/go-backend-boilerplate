package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/cache"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/db"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/external"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/storage"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	"github.com/azahir21/go-backend-boilerplate/pkg/apperr"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Application holds all application-wide dependencies.
type Application struct {
	Log          *logrus.Logger
	Config       *config.Config
	DBClient     *ent.Client
	Cache        cache.Cache
	Storage      storage.Storage
	EmailClient  external.EmailClient
	Dependencies *module.Dependencies
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

	// Initialize error responder configuration
	apperr.SetDefaultConfig(apperr.ConfigFromEnv(cfg.Server.Env))

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

	app := &Application{
		Log:          log,
		Config:       cfg,
		DBClient:     dbClient,
		Cache:        appCache,
		Storage:      appStorage,
		EmailClient:  emailClient,
		Dependencies: deps,
	}

	// Register modules based on build tags
	app.registerModules(deps)

	return app, nil
}

// registerModules calls the init functions for all enabled delivery layers.
func (app *Application) registerModules(deps *module.Dependencies) {
	// These functions are defined in build-tag-specific files
	// and will only be linked if the corresponding tag is present
	initRESTModules(deps)
	initGRPCModules(deps)
	initGraphQLModules(deps)
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
	app.shutdownServers(ctx)
	app.Log.Info("Servers stopped.")
	return nil
}

// setupServers creates servers based on configuration and build tags.
func (app *Application) setupServers() error {
	serversEnabled := false

	// Setup REST server if build tag is enabled
	if app.setupRESTServer() {
		serversEnabled = true
	}

	// Setup gRPC server if build tag is enabled
	if app.setupGRPCServer() {
		serversEnabled = true
	}

	// Setup GraphQL server if build tag is enabled
	if app.setupGraphQLServer() {
		serversEnabled = true
	}

	if !serversEnabled {
		return errors.New("no server enabled. Please enable at least one of HTTP, gRPC or GraphQL in config, and build with appropriate tags")
	}
	return nil
}

// startServers starts the enabled servers in goroutines.
func (app *Application) startServers(ctx context.Context) {
	app.startRESTServer(ctx)
	app.startGRPCServer(ctx)
	app.startGraphQLServer(ctx)
}

// shutdownServers gracefully shuts down all running servers.
func (app *Application) shutdownServers(ctx context.Context) {
	app.shutdownRESTServer(ctx)
	app.shutdownGRPCServer(ctx)
	app.shutdownGraphQLServer(ctx)
}

// Helper to listen on a TCP port
func listenTCP(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
