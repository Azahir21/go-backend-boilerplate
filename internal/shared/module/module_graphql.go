//go:build graphql
// +build graphql

package module

import (
	sharedGraphQL "github.com/azahir21/go-backend-boilerplate/internal/shared/graphql"
)

// GraphQLModule is implemented by modules that provide GraphQL handlers.
type GraphQLModule interface {
	Module
	// GraphQLSchemaBuilder returns the schema builder for this module.
	GraphQLSchemaBuilder() sharedGraphQL.SchemaBuilder
}
