package main

import (
	graphqlapp "github.com/azahir21/go-backend-boilerplate/cmd/app/graphql"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	if err := graphqlapp.Run(log); err != nil {
		log.Fatalf("GraphQL application failed to start: %v", err)
	}
}
