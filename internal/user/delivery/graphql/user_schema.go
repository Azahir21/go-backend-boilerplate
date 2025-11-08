package graphql

import (
	sharedGraphQL "github.com/azahir21/go-backend-boilerplate/internal/shared/graphql"
	"github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"role": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// UserSchemaBuilder implements the SchemaBuilder interface for the user module.
type UserSchemaBuilder struct {
	userResolver *UserResolver
}

// NewUserSchemaBuilder creates a new UserSchemaBuilder.
func NewUserSchemaBuilder(log *logrus.Logger, userUsecase usecase.UserUsecase) sharedGraphQL.SchemaBuilder {
	return &UserSchemaBuilder{
		userResolver: NewUserResolver(log, userUsecase),
	}
}

// BuildQueryFields returns the query fields for the user module.
func (b *UserSchemaBuilder) BuildQueryFields() graphql.Fields {
	return graphql.Fields{
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: b.userResolver.GetUserResolver,
		},
	}
}

// BuildMutationFields returns the mutation fields for the user module.
func (b *UserSchemaBuilder) BuildMutationFields() graphql.Fields {
	return graphql.Fields{
		"register": &graphql.Field{
			Type: userType, // Or a custom AuthResponse type
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: b.userResolver.RegisterUserResolver,
		},
		"login": &graphql.Field{
			Type: graphql.String, // Returns token
			Args: graphql.FieldConfigArgument{
				"username": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: b.userResolver.LoginUserResolver,
		},
	}
}
