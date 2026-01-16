//go:build !grpc
// +build !grpc

package app

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

// initGRPCModules is a no-op when grpc tag is not enabled.
func initGRPCModules(deps *module.Dependencies) {
	// No gRPC modules to initialize
}

// setupGRPCServer returns false when grpc tag is not enabled.
func (app *Application) setupGRPCServer() bool {
	return false
}

// startGRPCServer is a no-op when grpc tag is not enabled.
func (app *Application) startGRPCServer(ctx context.Context) {
	// No gRPC server to start
}

// shutdownGRPCServer is a no-op when grpc tag is not enabled.
func (app *Application) shutdownGRPCServer(ctx context.Context) {
	// No gRPC server to shutdown
}
