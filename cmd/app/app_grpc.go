//go:build grpc
// +build grpc

package app

import (
	"context"
	"fmt"
	"net"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"google.golang.org/grpc"
)

// GRPCModules holds gRPC-specific modules.
type GRPCModules struct {
	Modules    []module.GRPCModule
	GRPCServer *grpc.Server
}

// grpcModules stores gRPC modules when the build tag is enabled.
var grpcModules *GRPCModules

// registerConditionalModules registers gRPC modules when the grpc build tag is enabled.
func (app *Application) registerConditionalModules(deps *module.Dependencies) {
	grpcModules = &GRPCModules{
		Modules: registerGRPCModules(deps),
	}
}

// setupConditionalServers sets up the gRPC server when the grpc build tag is enabled.
func (app *Application) setupConditionalServers() error {
	if app.Config.Server.GRPC.Enable && grpcModules != nil {
		grpcServer, err := service.NewGrpcServer(app.Log, app.Config.Server.GRPC, grpcModules.Modules)
		if err != nil {
			return fmt.Errorf("failed to create gRPC server: %w", err)
		}
		grpcModules.GRPCServer = grpcServer
	}
	return nil
}

// startConditionalServers starts the gRPC server when the grpc build tag is enabled.
func (app *Application) startConditionalServers(ctx context.Context) {
	if app.Config.Server.GRPC.Enable && grpcModules != nil && grpcModules.GRPCServer != nil {
		grpcLis, err := net.Listen("tcp", ":"+app.Config.Server.GRPC.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on gRPC port %s: %v", app.Config.Server.GRPC.Port, err)
		} else {
			go func() {
				if err := grpcModules.GRPCServer.Serve(grpcLis); err != nil && err != grpc.ErrServerStopped {
					app.Log.Errorf("failed to start gRPC server: %v", err)
				}
			}()
		}
	}
}

// shutdownConditionalServers shuts down the gRPC server when the grpc build tag is enabled.
func (app *Application) shutdownConditionalServers(ctx context.Context) {
	if grpcModules != nil && grpcModules.GRPCServer != nil {
		grpcModules.GRPCServer.GracefulStop()
	}
}

// hasGRPCServer returns true when grpc tag is enabled and server is configured.
func (app *Application) hasGRPCServer() bool {
	return app.Config.Server.GRPC.Enable
}

// hasGraphQLServer returns false since this file only handles grpc.
func (app *Application) hasGraphQLServer() bool {
	return false
}
