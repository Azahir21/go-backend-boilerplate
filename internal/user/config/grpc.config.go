package config

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	grpcDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/grpc"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/proto"
	"google.golang.org/grpc"
)

// GRPCConfig provides dependency injection for the user module's gRPC handler.
type GRPCConfig struct {
	deps    *module.Dependencies
	handler *grpcDelivery.UserHandler
}

// NewGRPCConfig creates a new GRPCConfig with all dependencies injected.
func NewGRPCConfig(deps *module.Dependencies) *GRPCConfig {
	// Initialize repository
	userRepo := userRepoImpl.NewUserRepository(deps.DBClient)

	// Initialize unit of work
	uow := unitofwork.NewUnitOfWork(deps.DBClient)

	// Initialize usecase
	usecase := userUsecase.NewUserUsecase(userRepo, uow)

	// Initialize handler
	handler := grpcDelivery.NewUserHandler(deps.Log, usecase)

	return &GRPCConfig{
		deps:    deps,
		handler: handler,
	}
}

// Name returns the module name.
func (c *GRPCConfig) Name() string {
	return "user"
}

// RegisterGRPC registers the user module's gRPC services on the given server.
func (c *GRPCConfig) RegisterGRPC(server *grpc.Server) {
	proto.RegisterUserServiceServer(server, c.handler)
}
