//go:build grpc
// +build grpc

package app

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"google.golang.org/grpc"
)

var grpcServer *grpc.Server
var grpcModulesVar []module.GRPCModule

// initGRPCModules registers gRPC modules when the grpc build tag is enabled.
func initGRPCModules(deps *module.Dependencies) {
	grpcModulesVar = registerGRPCModules(deps)
}

// setupGRPCServer sets up the gRPC server when the grpc build tag is enabled.
func (app *Application) setupGRPCServer() bool {
	if !app.Config.Server.GRPC.Enable {
		return false
	}

	srv, err := service.NewGrpcServer(app.Log, app.Config.Server.GRPC, grpcModulesVar)
	if err != nil {
		app.Log.Errorf("failed to create gRPC server: %v", err)
		return false
	}
	grpcServer = srv
	return true
}

// startGRPCServer starts the gRPC server when the grpc build tag is enabled.
func (app *Application) startGRPCServer(ctx context.Context) {
	if app.Config.Server.GRPC.Enable && grpcServer != nil {
		grpcLis, err := listenTCP(":" + app.Config.Server.GRPC.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on gRPC port %s: %v", app.Config.Server.GRPC.Port, err)
		} else {
			go func() {
				if err := grpcServer.Serve(grpcLis); err != nil && err != grpc.ErrServerStopped {
					app.Log.Errorf("failed to start gRPC server: %v", err)
				}
			}()
		}
	}
}

// shutdownGRPCServer shuts down the gRPC server when the grpc build tag is enabled.
func (app *Application) shutdownGRPCServer(ctx context.Context) {
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}
