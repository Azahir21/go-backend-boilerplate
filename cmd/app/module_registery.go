package app

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
)

// registerModules registers all application modules and returns their handlers.
// Add new modules here when extending the application.
func registerModules(deps *module.Dependencies) ([]module.HTTPModule, []module.GRPCModule, []module.GraphQLModule) {
	var httpModules []module.HTTPModule
	var grpcModules []module.GRPCModule
	var graphqlModules []module.GraphQLModule

	// User module
	httpModules = append(httpModules, userConfig.NewHTTPConfig(deps))
	grpcModules = append(grpcModules, userConfig.NewGRPCConfig(deps))
	graphqlModules = append(graphqlModules, userConfig.NewGraphQLConfig(deps))

	// Add more modules here as needed:
	// productModule := productConfig.NewProductConfig(deps)
	// httpModules = append(httpModules, productConfig.NewHTTPConfig(deps))
	// grpcModules = append(grpcModules, productConfig.NewGRPCConfig(deps))
	// graphqlModules = append(graphqlModules, productConfig.NewGraphQLConfig(deps))

	return httpModules, grpcModules, graphqlModules
}
