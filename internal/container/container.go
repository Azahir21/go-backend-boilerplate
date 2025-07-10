package container

import (
	"github.com/azahir21/go-backend-boilerplate/internal/config"
	"github.com/azahir21/go-backend-boilerplate/internal/handler"
	"github.com/azahir21/go-backend-boilerplate/internal/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/repository"
	"github.com/azahir21/go-backend-boilerplate/internal/service"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
    Config      *config.Config
    DB          *gorm.DB
    Log         *logrus.Logger
    Repositories *Repositories
    Services     *Services
    Handlers     *Handlers
}

// Repositories holds all repository instances
type Repositories struct {
    UserRepository repository.UserRepository
}

// Services holds all service instances
type Services struct {
    AuthService service.AuthService
}

// Handlers holds all handler instances
type Handlers struct {
    AuthHandler *handler.AuthHandler
}

// Singleton instance
var instance *Container

// GetInstance returns the singleton container instance
func GetInstance() *Container {
    return instance
}

// Initialize creates and configures the application container
func Initialize(cfg *config.Config, db *gorm.DB, log *logrus.Logger) {
    if instance == nil {
        instance = &Container{
            Config: cfg,
            DB:     db,
            Log:    log,
        }

        // Initialize JWT helper with config
        helper.InitJWT(cfg.JWTSecret, cfg.JWTExpiryHours)
        
        // Initialize repositories
        instance.initRepositories()
        
        // Initialize services
        instance.initServices()
        
        // Initialize handlers
        instance.initHandlers()
    }
}

// Initialize all repositories
func (c *Container) initRepositories() {
    c.Repositories = &Repositories{
        UserRepository: repository.NewUserRepository(c.DB),
    }
}

// Initialize all services
func (c *Container) initServices() {
    c.Services = &Services{
        AuthService: service.NewAuthService(c.Repositories.UserRepository),
    }
}

// Initialize all handlers
func (c *Container) initHandlers() {
    c.Handlers = &Handlers{
        AuthHandler: handler.NewHandler(c.DB, c.Log, c.Services.AuthService),
    }
}