/*
Package apperr provides structured error handling with status codes, messages,
error details, and stacktraces for Go web applications.

# Overview

The apperr package offers a comprehensive error handling solution that:
  - Provides standardized error statuses mapped to HTTP codes
  - Captures stacktraces automatically for debugging
  - Hides sensitive error details in production environments
  - Integrates seamlessly with Gin framework

# Error Statuses

The package supports the following error statuses:

	bad_gateway          -> 502
	bad_request          -> 400
	conflict             -> 409
	forbidden            -> 403
	internal_server      -> 500
	method_not_allowed   -> 405
	not_found            -> 404
	not_implemented      -> 501
	service_unavailable  -> 503
	timeout              -> 504
	too_many_requests    -> 429
	unauthorized         -> 401
	unprocessable_entity -> 422

# Basic Usage

Creating errors:

	// Simple error creation
	err := apperr.NotFound("user not found")

	// With additional detail
	err := apperr.BadRequest("validation failed").WithDetail("email field is required")

	// Wrapping an existing error
	err := apperr.Wrap(dbErr, apperr.StatusInternalServer, "failed to fetch user")

	// With cause (sets detail automatically from cause)
	err := apperr.InternalServer("database error").WithCause(originalErr)

# Configuration

Configure the error responder based on environment:

	// In your application setup (e.g., main.go or app.go)
	import "github.com/azahir21/go-backend-boilerplate/pkg/apperr"

	// Using environment name
	apperr.SetDefaultConfig(apperr.ConfigFromEnv(cfg.Server.Env))

	// Or manually
	apperr.SetDefaultConfig(apperr.DevelopmentConfig()) // Shows detail and stacktrace
	apperr.SetDefaultConfig(apperr.DefaultConfig())     // Hides detail and stacktrace (production)

# Responding to Errors

Using the package-level functions:

	func GetUser(c *gin.Context) {
		user, err := userService.Get(id)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				apperr.RespondNotFound(c, "User not found")
				return
			}
			apperr.RespondError(c, err) // Automatically converts to AppError
			return
		}
		// success response...
	}

Using AppError directly:

	func CreateUser(c *gin.Context) {
		err := userService.Create(user)
		if err != nil {
			appErr := apperr.BadRequest("Invalid user data").WithCause(err)
			apperr.Respond(c, appErr)
			return
		}
		// success response...
	}

# Middleware Integration

Add the error handler middleware to your Gin router:

	import (
		"github.com/azahir21/go-backend-boilerplate/pkg/apperr"
		"github.com/gin-gonic/gin"
	)

	func setupRouter(log *logrus.Logger, env string) *gin.Engine {
		r := gin.New()

		// Add error handling middleware
		config := apperr.ConfigFromEnv(env)
		r.Use(apperr.ErrorHandlerMiddleware(log, config))

		// ... register routes
		return r
	}

Using AbortWithError in handlers:

	func DeleteUser(c *gin.Context) {
		if !hasPermission(c) {
			apperr.AbortWithError(c, apperr.Forbidden("You don't have permission"))
			return
		}
		// ...
	}

# Response Format

In development environment (detail and stacktrace visible):

	{
		"status": "bad_request",
		"message": "Validation failed",
		"detail": "email field is required",
		"stacktrace": "main.CreateUser\n\t/app/handlers/user.go:42\n..."
	}

In production environment (detail and stacktrace hidden):

	{
		"status": "bad_request",
		"message": "Validation failed"
	}

# Error Checking

Check if an error has a specific status:

	if apperr.Is(err, apperr.StatusNotFound) {
		// handle not found case
	}

Convert any error to AppError:

	appErr := apperr.AsAppError(err)
	// If err is already an AppError, returns it unchanged
	// Otherwise, wraps it as an internal_server error
*/
package apperr
