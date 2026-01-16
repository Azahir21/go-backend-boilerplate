package config

import (
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	restDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
)

// HTTPConfig provides dependency injection for the user module's HTTP handler.
type HTTPConfig struct {
	deps    *module.Dependencies
	handler *restDelivery.UserHandler
}

// NewHTTPConfig creates a new HTTPConfig with all dependencies injected.
func NewHTTPConfig(deps *module.Dependencies) *HTTPConfig {
	// Initialize repository
	userRepo := userRepoImpl.NewUserRepository(deps.DBClient)

	// Initialize unit of work
	uow := unitofwork.NewUnitOfWork(deps.DBClient)

	// Initialize usecase
	usecase := userUsecase.NewUserUsecase(userRepo, uow)

	// Initialize handler
	handler := restDelivery.NewUserHandler(deps.Log, usecase)

	return &HTTPConfig{
		deps:    deps,
		handler: handler,
	}
}

// Name returns the module name.
func (c *HTTPConfig) Name() string {
	return "user"
}

// HTTPHandler returns the HTTP handler for the user module.
func (c *HTTPConfig) HTTPHandler() sharedHttp.HttpRouter {
	return c.handler
}
