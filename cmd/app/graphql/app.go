package graphql

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
	"github.com/azahir21/go-backend-boilerplate/infrastructure/db/mongo"
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

// Application holds all application-wide dependencies for GraphQL delivery.
type Application struct {
	Log            *logrus.Logger
	Config         *config.Config
	DBClient       *ent.Client
	MongoClient    *mongo.Client
	Cache          cache.Cache
	Storage        storage.Storage
	EmailClient    external.EmailClient
	Dependencies   *module.Dependencies
	GraphQLModules []module.GraphQLModule
	GraphQLServer  *http.Server
}

// NewApplication initializes and returns a new GraphQL Application instance.
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

	// Initialize database (optional)
	var dbClient *ent.Client
	var uow unitofwork.UnitOfWork
	if cfg.DB.Enable {
		var err error
		dbClient, err = db.NewEntClient(log, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
		uow = unitofwork.NewUnitOfWork(dbClient)
	} else {
		log.Info("SQL Database is disabled, skipping initialization")
	}

	// Initialize MongoDB (optional)
	var mongoClient *mongo.Client
	if cfg.Mongo.Enable {
		var err error
		mongoClient, err = mongo.NewClient(log, cfg.Mongo)
		if err != nil {
			if dbClient != nil {
				dbClient.Close()
			}
			return nil, fmt.Errorf("failed to initialize MongoDB: %w", err)
		}
	} else {
		log.Info("MongoDB is disabled, skipping initialization")
	}

	// Initialize cache (optional)
	var appCache cache.Cache
	if cfg.Cache.Enable {
		var err error
		appCache, err = cache.NewCache(log, cfg.Cache)
		if err != nil {
			if dbClient != nil {
				dbClient.Close()
			}
			if mongoClient != nil {
				mongoClient.Close()
			}
			return nil, fmt.Errorf("failed to initialize cache: %w", err)
		}
	} else {
		log.Info("Cache is disabled, skipping initialization")
	}

	// Initialize storage (optional)
	var appStorage storage.Storage
	if cfg.Storage.Enable {
		var err error
		appStorage, err = storage.NewStorage(context.Background(), log, cfg.Storage)
		if err != nil {
			if dbClient != nil {
				dbClient.Close()
			}
			if mongoClient != nil {
				mongoClient.Close()
			}
			return nil, fmt.Errorf("failed to initialize storage: %w", err)
		}
	} else {
		log.Info("Storage is disabled, skipping initialization")
	}

	// Initialize email client (optional)
	var emailClient external.EmailClient
	if cfg.Email.Enable {
		var err error
		emailClient, err = external.NewEmailClient(log, cfg.Email)
		if err != nil {
			if dbClient != nil {
				dbClient.Close()
			}
			if mongoClient != nil {
				mongoClient.Close()
			}
			return nil, fmt.Errorf("failed to initialize email client: %w", err)
		}
	} else {
		log.Info("Email client is disabled, skipping initialization")
	}

	// Create shared dependencies for all modules
	deps := &module.Dependencies{
		Log:         log,
		DBClient:    dbClient,
		MongoClient: mongoClient,
		Cache:       appCache,
		Storage:     appStorage,
		EmailClient: emailClient,
		UoW:         uow,
	}

	// Register only GraphQL modules
	var graphqlModules []module.GraphQLModule
	// User module - only register if SQL database is enabled
	if dbClient != nil {
		graphqlModules = append(graphqlModules, userConfig.NewGraphQLConfig(deps))
	}

	return &Application{
		Log:            log,
		Config:         cfg,
		DBClient:       dbClient,
		MongoClient:    mongoClient,
		Cache:          appCache,
		Storage:        appStorage,
		EmailClient:    emailClient,
		Dependencies:   deps,
		GraphQLModules: graphqlModules,
	}, nil
}

func Run(log *logrus.Logger) error {
	app, err := NewApplication(log)
	if err != nil {
		return fmt.Errorf("failed to initialize GraphQL application: %w", err)
	}
	// Ensure DB clients are closed
	defer func() {
		if app.DBClient != nil {
			app.DBClient.Close()
		}
		if app.MongoClient != nil {
			app.MongoClient.Close()
		}
	}()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup GraphQL server
	if err := app.setupServer(); err != nil {
		return fmt.Errorf("failed to setup GraphQL server: %w", err)
	}

	// Start GraphQL server
	app.startServer(ctx)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	app.Log.Info("Shutting down GraphQL server...")
	if app.GraphQLServer != nil {
		app.GraphQLServer.Shutdown(ctx)
	}
	app.Log.Info("GraphQL server stopped.")
	return nil
}

// setupServer creates the GraphQL server based on configuration.
func (app *Application) setupServer() error {
	if !app.Config.Server.GraphQL.Enable {
		return errors.New("GraphQL server is disabled in configuration")
	}

	srv, err := service.NewGraphQLServer(app.Log, app.Config.Server.GraphQL, app.GraphQLModules)
	if err != nil {
		return fmt.Errorf("failed to create GraphQL server: %w", err)
	}
	app.GraphQLServer = srv
	return nil
}

// startServer starts the GraphQL server in a goroutine.
func (app *Application) startServer(ctx context.Context) {
	if app.Config.Server.GraphQL.Enable {
		graphqlLis, err := net.Listen("tcp", ":"+app.Config.Server.GraphQL.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on GraphQL port %s: %v", app.Config.Server.GraphQL.Port, err)
		} else {
			go func() {
				if err := app.GraphQLServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start GraphQL server: %v", err)
				}
			}()
		}
	}
}
