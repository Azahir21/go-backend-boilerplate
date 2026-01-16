//go:build !rest
// +build !rest

package app

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

// initRESTModules is a no-op when rest tag is not enabled.
func initRESTModules(deps *module.Dependencies) {
	// No REST modules to initialize
}

// setupRESTServer returns false when rest tag is not enabled.
func (app *Application) setupRESTServer() bool {
	return false
}

// startRESTServer is a no-op when rest tag is not enabled.
func (app *Application) startRESTServer(ctx context.Context) {
	// No REST server to start
}

// shutdownRESTServer is a no-op when rest tag is not enabled.
func (app *Application) shutdownRESTServer(ctx context.Context) {
	// No REST server to shutdown
}
