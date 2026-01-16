package apperr

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandlerMiddleware returns a Gin middleware that catches panics and handles errors.
// It converts panics to InternalServer errors and logs them appropriately.
func ErrorHandlerMiddleware(log *logrus.Logger, config Config) gin.HandlerFunc {
	responder := NewResponder(config)

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var appErr *AppError

				switch v := r.(type) {
				case *AppError:
					appErr = v
				case error:
					appErr = InternalServer("An unexpected error occurred").WithCause(v)
				case string:
					appErr = InternalServer(v)
				default:
					appErr = InternalServer("An unexpected error occurred")
				}

				// Log the error with full details
				log.WithFields(logrus.Fields{
					"status":     appErr.Status,
					"message":    appErr.Message,
					"detail":     appErr.Detail,
					"stacktrace": appErr.Stacktrace,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				responder.Respond(c, appErr)
				c.Abort()
			}
		}()

		c.Next()

		// Check if there are any errors set in the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			appErr := AsAppError(err)

			log.WithFields(logrus.Fields{
				"status":     appErr.Status,
				"message":    appErr.Message,
				"detail":     appErr.Detail,
				"path":       c.Request.URL.Path,
				"method":     c.Request.Method,
			}).Error("Request error")

			responder.Respond(c, appErr)
		}
	}
}

// RecoveryMiddleware is a simpler middleware that only handles panic recovery.
// Use this if you want minimal panic handling without full error management.
func RecoveryMiddleware(log *logrus.Logger, config Config) gin.HandlerFunc {
	responder := NewResponder(config)

	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				appErr := InternalServer("An unexpected error occurred")

				if err, ok := r.(error); ok {
					appErr = appErr.WithCause(err)
				}

				log.WithFields(logrus.Fields{
					"panic":      r,
					"stacktrace": appErr.Stacktrace,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				responder.Respond(c, appErr)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// AbortWithError is a helper that sets an error in the Gin context and aborts.
// Use this in handlers to trigger error handling by the middleware.
func AbortWithError(c *gin.Context, err *AppError) {
	_ = c.Error(err)
	c.Abort()
}
