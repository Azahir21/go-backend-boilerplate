package app

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
)

// registerModules registers all application modules and returns their handlers.
// Add new modules here when extending the application.
func registerModules(deps *module.Dependencies) []module.HTTPModule {
	var httpModules []module.HTTPModule

	// User module - HTTP is always registered
	httpModules = append(httpModules, userConfig.NewHTTPConfig(deps))

	// Add more HTTP modules here as needed:
	// httpModules = append(httpModules, productConfig.NewHTTPConfig(deps))

	return httpModules
}
