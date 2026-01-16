//go:build graphql
// +build graphql

package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

// GraphQLModules holds GraphQL-specific modules.
type GraphQLModules struct {
	Modules       []module.GraphQLModule
	GraphQLServer *http.Server
}

// graphqlModules stores GraphQL modules when the build tag is enabled.
var graphqlModules *GraphQLModules

// registerConditionalModules registers GraphQL modules when the graphql build tag is enabled.
func (app *Application) registerConditionalModules(deps *module.Dependencies) {
	graphqlModules = &GraphQLModules{
		Modules: registerGraphQLModules(deps),
	}
}

// setupConditionalServers sets up the GraphQL server when the graphql build tag is enabled.
func (app *Application) setupConditionalServers() error {
	if app.Config.Server.GraphQL.Enable && graphqlModules != nil {
		srv, err := service.NewGraphQLServer(app.Log, app.Config.Server.GraphQL, graphqlModules.Modules)
		if err != nil {
			return fmt.Errorf("failed to create GraphQL server: %w", err)
		}
		graphqlModules.GraphQLServer = srv
	}
	return nil
}

// startConditionalServers starts the GraphQL server when the graphql build tag is enabled.
func (app *Application) startConditionalServers(ctx context.Context) {
	if app.Config.Server.GraphQL.Enable && graphqlModules != nil && graphqlModules.GraphQLServer != nil {
		graphqlLis, err := net.Listen("tcp", ":"+app.Config.Server.GraphQL.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on GraphQL port %s: %v", app.Config.Server.GraphQL.Port, err)
		} else {
			go func() {
				if err := graphqlModules.GraphQLServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start GraphQL server: %v", err)
				}
			}()
		}
	}
}

// shutdownConditionalServers shuts down the GraphQL server when the graphql build tag is enabled.
func (app *Application) shutdownConditionalServers(ctx context.Context) {
	if graphqlModules != nil && graphqlModules.GraphQLServer != nil {
		graphqlModules.GraphQLServer.Shutdown(ctx)
	}
}

// hasGRPCServer returns false since this file only handles graphql.
func (app *Application) hasGRPCServer() bool {
	return false
}

// hasGraphQLServer returns true when graphql tag is enabled and server is configured.
func (app *Application) hasGraphQLServer() bool {
	return app.Config.Server.GraphQL.Enable
}
