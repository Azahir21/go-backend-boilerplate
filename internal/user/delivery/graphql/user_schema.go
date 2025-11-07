package graphql

import (
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

func NewGraphQLSchema(log *logrus.Logger, userUsecase usecase.UserUsecase) (graphql.Schema, error) {
	userResolver := NewUserResolver(log, userUsecase)

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: userResolver.GetUserResolver,
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
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
				Resolve: userResolver.RegisterUserResolver,
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
				Resolve: userResolver.LoginUserResolver,
			},
		},
	})

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rootQuery,
			Mutation: rootMutation,
		},
	)
}