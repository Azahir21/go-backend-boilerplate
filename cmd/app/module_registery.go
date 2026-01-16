//go:build rest
// +build rest

package app

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
)

// registerRESTModules registers REST/HTTP modules when the rest build tag is enabled.
func registerRESTModules(deps *module.Dependencies) []module.HTTPModule {
	var httpModules []module.HTTPModule

	// User module - HTTP
	httpModules = append(httpModules, userConfig.NewHTTPConfig(deps))

	// Add more HTTP modules here as needed:
	// httpModules = append(httpModules, productConfig.NewHTTPConfig(deps))

	return httpModules
}
