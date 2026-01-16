package rest

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
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Application holds all application-wide dependencies for REST delivery.
type Application struct {
	Log          *logrus.Logger
	Config       *config.Config
	DBClient     *ent.Client
	Cache        cache.Cache
	Storage      storage.Storage
	EmailClient  external.EmailClient
	Dependencies *module.Dependencies
	HTTPModules  []module.HTTPModule
	HTTPServer   *http.Server
}

// NewApplication initializes and returns a new REST Application instance.
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

	// Register only HTTP modules (REST delivery)
	var httpModules []module.HTTPModule
	httpModules = append(httpModules, userConfig.NewHTTPConfig(deps))

	return &Application{
		Log:          log,
		Config:       cfg,
		DBClient:     dbClient,
		Cache:        appCache,
		Storage:      appStorage,
		EmailClient:  emailClient,
		Dependencies: deps,
		HTTPModules:  httpModules,
	}, nil
}

func Run(log *logrus.Logger) error {
	app, err := NewApplication(log)
	if err != nil {
		return fmt.Errorf("failed to initialize REST application: %w", err)
	}
	defer app.DBClient.Close()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup HTTP server
	if err := app.setupServer(); err != nil {
		return fmt.Errorf("failed to setup HTTP server: %w", err)
	}

	// Start HTTP server
	app.startServer(ctx)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	app.Log.Info("Shutting down HTTP server...")
	if app.HTTPServer != nil {
		app.HTTPServer.Shutdown(ctx)
	}
	app.Log.Info("HTTP server stopped.")
	return nil
}

// setupServer creates the HTTP server based on configuration.
func (app *Application) setupServer() error {
	if !app.Config.Server.HTTP.Enable {
		return errors.New("HTTP server is disabled in configuration")
	}

	srv, err := service.NewRestServer(app.Log, app.Config.Server.HTTP, app.HTTPModules)
	if err != nil {
		return fmt.Errorf("failed to create HTTP server: %w", err)
	}
	app.HTTPServer = srv
	return nil
}

// startServer starts the HTTP server in a goroutine.
func (app *Application) startServer(ctx context.Context) {
	if app.Config.Server.HTTP.Enable {
		httpLis, err := net.Listen("tcp", ":"+app.Config.Server.HTTP.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on HTTP port %s: %v", app.Config.Server.HTTP.Port, err)
		} else {
			go func() {
				if err := app.HTTPServer.Serve(httpLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start HTTP server: %v", err)
				}
			}()
		}
	}
}
