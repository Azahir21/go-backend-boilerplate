//go:build grpc && graphql
// +build grpc,graphql

package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"google.golang.org/grpc"
)

// CombinedModules holds both gRPC and GraphQL modules.
type CombinedModules struct {
	GRPCModules    []module.GRPCModule
	GraphQLModules []module.GraphQLModule
	GRPCServer     *grpc.Server
	GraphQLServer  *http.Server
}

// combinedModules stores both gRPC and GraphQL modules when both build tags are enabled.
var combinedModules *CombinedModules

// registerConditionalModules registers both gRPC and GraphQL modules when both build tags are enabled.
func (app *Application) registerConditionalModules(deps *module.Dependencies) {
	combinedModules = &CombinedModules{
		GRPCModules:    registerGRPCModules(deps),
		GraphQLModules: registerGraphQLModules(deps),
	}
}

// setupConditionalServers sets up both gRPC and GraphQL servers when both build tags are enabled.
func (app *Application) setupConditionalServers() error {
	if app.Config.Server.GRPC.Enable && combinedModules != nil {
		grpcServer, err := service.NewGrpcServer(app.Log, app.Config.Server.GRPC, combinedModules.GRPCModules)
		if err != nil {
			return fmt.Errorf("failed to create gRPC server: %w", err)
		}
		combinedModules.GRPCServer = grpcServer
	}

	if app.Config.Server.GraphQL.Enable && combinedModules != nil {
		srv, err := service.NewGraphQLServer(app.Log, app.Config.Server.GraphQL, combinedModules.GraphQLModules)
		if err != nil {
			return fmt.Errorf("failed to create GraphQL server: %w", err)
		}
		combinedModules.GraphQLServer = srv
	}

	return nil
}

// startConditionalServers starts both gRPC and GraphQL servers when both build tags are enabled.
func (app *Application) startConditionalServers(ctx context.Context) {
	// Start gRPC server
	if app.Config.Server.GRPC.Enable && combinedModules != nil && combinedModules.GRPCServer != nil {
		grpcLis, err := net.Listen("tcp", ":"+app.Config.Server.GRPC.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on gRPC port %s: %v", app.Config.Server.GRPC.Port, err)
		} else {
			go func() {
				if err := combinedModules.GRPCServer.Serve(grpcLis); err != nil && err != grpc.ErrServerStopped {
					app.Log.Errorf("failed to start gRPC server: %v", err)
				}
			}()
		}
	}

	// Start GraphQL server
	if app.Config.Server.GraphQL.Enable && combinedModules != nil && combinedModules.GraphQLServer != nil {
		graphqlLis, err := net.Listen("tcp", ":"+app.Config.Server.GraphQL.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on GraphQL port %s: %v", app.Config.Server.GraphQL.Port, err)
		} else {
			go func() {
				if err := combinedModules.GraphQLServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start GraphQL server: %v", err)
				}
			}()
		}
	}
}

// shutdownConditionalServers shuts down both gRPC and GraphQL servers when both build tags are enabled.
func (app *Application) shutdownConditionalServers(ctx context.Context) {
	if combinedModules != nil {
		if combinedModules.GRPCServer != nil {
			combinedModules.GRPCServer.GracefulStop()
		}
		if combinedModules.GraphQLServer != nil {
			combinedModules.GraphQLServer.Shutdown(ctx)
		}
	}
}

// hasGRPCServer returns true when grpc tag is enabled and server is configured.
func (app *Application) hasGRPCServer() bool {
	return app.Config.Server.GRPC.Enable
}

// hasGraphQLServer returns true when graphql tag is enabled and server is configured.
func (app *Application) hasGraphQLServer() bool {
	return app.Config.Server.GraphQL.Enable
}
