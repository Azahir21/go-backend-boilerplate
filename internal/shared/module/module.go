package module

import (
	sharedGraphQL "github.com/azahir21/go-backend-boilerplate/internal/shared/graphql"
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
	"google.golang.org/grpc"
)

// Module defines the interface that each feature module must implement.
// It provides a unified way to register handlers across different delivery mechanisms.
type Module interface {
	// Name returns the module name for logging/debugging purposes.
	Name() string
}

// HTTPModule is implemented by modules that provide HTTP handlers.
type HTTPModule interface {
	Module
	// HTTPHandler returns the HTTP router for this module.
	HTTPHandler() sharedHttp.HttpRouter
}

// GRPCModule is implemented by modules that provide gRPC handlers.
type GRPCModule interface {
	Module
	// RegisterGRPC registers the module's gRPC services on the given server.
	RegisterGRPC(server *grpc.Server)
}

// GraphQLModule is implemented by modules that provide GraphQL handlers.
type GraphQLModule interface {
	Module
	// GraphQLSchemaBuilder returns the schema builder for this module.
	GraphQLSchemaBuilder() sharedGraphQL.SchemaBuilder
}
