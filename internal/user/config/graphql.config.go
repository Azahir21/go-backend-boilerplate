package config

import (
	sharedGraphQL "github.com/azahir21/go-backend-boilerplate/internal/shared/graphql"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	graphqlDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/graphql"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
)

// GraphQLConfig provides dependency injection for the user module's GraphQL handler.
type GraphQLConfig struct {
	deps          *module.Dependencies
	schemaBuilder sharedGraphQL.SchemaBuilder
}

// NewGraphQLConfig creates a new GraphQLConfig with all dependencies injected.
func NewGraphQLConfig(deps *module.Dependencies) *GraphQLConfig {
	// Initialize repository
	userRepo := userRepoImpl.NewUserRepository(deps.DBClient)

	// Initialize unit of work
	uow := unitofwork.NewUnitOfWork(deps.DBClient)

	// Initialize usecase
	usecase := userUsecase.NewUserUsecase(userRepo, uow)

	// Initialize schema builder
	schemaBuilder := graphqlDelivery.NewUserSchemaBuilder(deps.Log, usecase)

	return &GraphQLConfig{
		deps:          deps,
		schemaBuilder: schemaBuilder,
	}
}

// Name returns the module name.
func (c *GraphQLConfig) Name() string {
	return "user"
}

// GraphQLSchemaBuilder returns the GraphQL schema builder for the user module.
func (c *GraphQLConfig) GraphQLSchemaBuilder() sharedGraphQL.SchemaBuilder {
	return c.schemaBuilder
}
