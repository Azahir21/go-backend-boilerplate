//go:build !grpc && !graphql
// +build !grpc,!graphql

package app

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

// registerConditionalModules is a no-op when neither grpc nor graphql tags are enabled.
func (app *Application) registerConditionalModules(deps *module.Dependencies) {
	// No additional modules to register
}

// setupConditionalServers is a no-op when neither grpc nor graphql tags are enabled.
func (app *Application) setupConditionalServers() error {
	return nil
}

// startConditionalServers is a no-op when neither grpc nor graphql tags are enabled.
func (app *Application) startConditionalServers(ctx context.Context) {
	// No additional servers to start
}

// shutdownConditionalServers is a no-op when neither grpc nor graphql tags are enabled.
func (app *Application) shutdownConditionalServers(ctx context.Context) {
	// No additional servers to shutdown
}

// hasGRPCServer returns false when grpc tag is not enabled.
func (app *Application) hasGRPCServer() bool {
	return false
}

// hasGraphQLServer returns false when graphql tag is not enabled.
func (app *Application) hasGraphQLServer() bool {
	return false
}
