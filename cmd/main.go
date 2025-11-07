package main

import (
	"github.com/azahir21/go-backend-boilerplate/cmd/app"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"
)

// @title Headcount Checker API
// @version 1.0
// @description API Documentation for Headcount Checker Backend
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
	if err := app.Run(log); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
