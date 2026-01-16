//go:build graphql
// +build graphql

package app

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
)

// registerGraphQLModules returns GraphQL modules when the graphql build tag is enabled.
func registerGraphQLModules(deps *module.Dependencies) []module.GraphQLModule {
	var graphqlModules []module.GraphQLModule
	
	// User module
	graphqlModules = append(graphqlModules, userConfig.NewGraphQLConfig(deps))

	// Add more modules here as needed:
	// graphqlModules = append(graphqlModules, productConfig.NewGraphQLConfig(deps))

	return graphqlModules
}
