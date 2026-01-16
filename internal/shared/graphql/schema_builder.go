//go:build graphql
// +build graphql

package graphql

import (
	"github.com/graphql-go/graphql"
)

// SchemaBuilder defines an interface for modules to contribute to the overall GraphQL schema.
type SchemaBuilder interface {
	BuildQueryFields() graphql.Fields
	BuildMutationFields() graphql.Fields
}

// NewRootSchema creates a new GraphQL schema by combining fields from multiple SchemaBuilders.
func NewRootSchema(builders []SchemaBuilder) (graphql.Schema, error) {
	queryFields := make(graphql.Fields)
	mutationFields := make(graphql.Fields)

	for _, builder := range builders {
		for name, field := range builder.BuildQueryFields() {
			queryFields[name] = field
		}
		for name, field := range builder.BuildMutationFields() {
			mutationFields[name] = field
		}
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: queryFields,
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: mutationFields,
	})

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rootQuery,
			Mutation: rootMutation,
		},
	)
}
