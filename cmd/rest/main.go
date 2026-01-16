package main

import (
	restapp "github.com/azahir21/go-backend-boilerplate/cmd/app/rest"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"
)

// @title REST API
// @version 1.0
// @description REST API Documentation for Go Backend Boilerplate
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log := logger.NewLogger()
	if err := restapp.Run(log); err != nil {
		log.Fatalf("REST application failed to start: %v", err)
	}
}
