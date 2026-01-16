//go:build rest
// +build rest

package app

import (
	"context"
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/cmd/service"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
)

var restServer *http.Server
var restModules []module.HTTPModule

// initRESTModules registers REST modules when the rest build tag is enabled.
func initRESTModules(deps *module.Dependencies) {
	restModules = registerRESTModules(deps)
}

// setupRESTServer sets up the REST/HTTP server when the rest build tag is enabled.
func (app *Application) setupRESTServer() bool {
	if !app.Config.Server.HTTP.Enable {
		return false
	}

	srv, err := service.NewRestServer(app.Log, app.Config.Server.HTTP, restModules)
	if err != nil {
		app.Log.Errorf("failed to create HTTP server: %v", err)
		return false
	}
	restServer = srv
	return true
}

// startRESTServer starts the REST/HTTP server when the rest build tag is enabled.
func (app *Application) startRESTServer(ctx context.Context) {
	if app.Config.Server.HTTP.Enable && restServer != nil {
		httpLis, err := listenTCP(":" + app.Config.Server.HTTP.Port)
		if err != nil {
			app.Log.Errorf("failed to listen on HTTP port %s: %v", app.Config.Server.HTTP.Port, err)
		} else {
			go func() {
				if err := restServer.Serve(httpLis); err != nil && err != http.ErrServerClosed {
					app.Log.Errorf("failed to start HTTP server: %v", err)
				}
			}()
		}
	}
}

// shutdownRESTServer shuts down the REST/HTTP server when the rest build tag is enabled.
func (app *Application) shutdownRESTServer(ctx context.Context) {
	if restServer != nil {
		restServer.Shutdown(ctx)
	}
}
