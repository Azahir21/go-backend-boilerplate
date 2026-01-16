//go:build !graphql
// +build !graphql

package app

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

// initGraphQLModules is a no-op when graphql tag is not enabled.
func initGraphQLModules(deps *module.Dependencies) {
	// No GraphQL modules to initialize
}

// setupGraphQLServer returns false when graphql tag is not enabled.
func (app *Application) setupGraphQLServer() bool {
	return false
}

// startGraphQLServer is a no-op when graphql tag is not enabled.
func (app *Application) startGraphQLServer(ctx context.Context) {
	// No GraphQL server to start
}

// shutdownGraphQLServer is a no-op when graphql tag is not enabled.
func (app *Application) shutdownGraphQLServer(ctx context.Context) {
	// No GraphQL server to shutdown
}
