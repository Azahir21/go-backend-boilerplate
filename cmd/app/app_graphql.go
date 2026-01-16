//go:build graphql
// +build graphql

package app

import (
	"context"
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

var graphqlServer *http.Server
var graphqlModulesVar []module.GraphQLModule

// initGraphQLModules registers GraphQL modules when the graphql build tag is enabled.
func initGraphQLModules(deps *module.Dependencies) {
	graphqlModulesVar = registerGraphQLModules(deps)
}

// setupGraphQLServer sets up the GraphQL server when the graphql build tag is enabled.
func (app *Application) setupGraphQLServer() bool {
	if !app.Config.Server.GraphQL.Enable {
		return false
	}

	srv, err := service.NewGraphQLServer(app.Log, app.Config.Server.GraphQL, graphqlModulesVar)
	if err != nil {
		app.Log.Errorf("failed to create GraphQL server: %v", err)
		return false
	}
	graphqlServer = srv
	return true
}

// startGraphQLServer starts the GraphQL server when the graphql build tag is enabled.
func (app *Application) startGraphQLServer(ctx context.Context) {
	if app.Config.Server.GraphQL.Enable && graphqlServer != nil {
		graphqlLis, err := listenTCP(":" + app.Config.Server.GraphQL.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on GraphQL port %s: %v", app.Config.Server.GraphQL.Port, err)
		} else {
			go func() {
				if err := graphqlServer.Serve(graphqlLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start GraphQL server: %v", err)
				}
			}()
		}
	}
}

// shutdownGraphQLServer shuts down the GraphQL server when the graphql build tag is enabled.
func (app *Application) shutdownGraphQLServer(ctx context.Context) {
	if graphqlServer != nil {
		graphqlServer.Shutdown(ctx)
	}
}
