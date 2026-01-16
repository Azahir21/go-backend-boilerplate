package apperr

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the JSON structure for error responses.
type ErrorResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

// Config holds configuration for error response handling.
type Config struct {
	// ShowDetail controls whether error details are included in the response.
	ShowDetail bool
	// ShowStacktrace controls whether stacktrace is included in the response.
	ShowStacktrace bool
}

// DefaultConfig returns the default configuration for error responses.
// In production environments, detail and stacktrace should be hidden.
func DefaultConfig() Config {
	return Config{
		ShowDetail:     false,
		ShowStacktrace: false,
	}
}

// DevelopmentConfig returns a configuration suitable for development environments.
// Both detail and stacktrace are visible for debugging.
func DevelopmentConfig() Config {
	return Config{
		ShowDetail:     true,
		ShowStacktrace: true,
	}
}

// ConfigFromEnv returns a configuration based on the environment.
// For "production" and "staging", sensitive information is hidden.
// For "development" and other environments, debugging information is shown.
func ConfigFromEnv(env string) Config {
	switch env {
	case "production", "staging":
		return DefaultConfig()
	default:
		return DevelopmentConfig()
	}
}

// Responder handles error responses with the appropriate configuration.
type Responder struct {
	config Config
}

// NewResponder creates a new Responder with the given configuration.
func NewResponder(config Config) *Responder {
	return &Responder{config: config}
}

// Respond sends an error response based on the AppError and configuration.
func (r *Responder) Respond(c *gin.Context, err *AppError) {
	if err == nil {
		return
	}

	response := ErrorResponse{
		Status:  err.Status.String(),
		Message: err.Message,
	}

	if r.config.ShowDetail {
		response.Detail = err.Detail
	}

	if r.config.ShowStacktrace {
		response.Stacktrace = err.Stacktrace
	}

	c.JSON(err.Status.HTTPCode(), response)
}

// RespondError converts a generic error to AppError and sends the response.
func (r *Responder) RespondError(c *gin.Context, err error) {
	r.Respond(c, AsAppError(err))
}

// --- Package-level responder for convenience ---

var defaultResponder = NewResponder(DefaultConfig())

// SetDefaultConfig sets the default configuration for the package-level responder.
func SetDefaultConfig(config Config) {
	defaultResponder = NewResponder(config)
}

// Respond sends an error response using the default responder.
func Respond(c *gin.Context, err *AppError) {
	defaultResponder.Respond(c, err)
}

// RespondError converts a generic error to AppError and sends the response.
func RespondError(c *gin.Context, err error) {
	defaultResponder.RespondError(c, err)
}

// --- Convenience functions that create and respond in one call ---

// RespondBadGateway creates and responds with a BadGateway error.
func RespondBadGateway(c *gin.Context, message string) {
	Respond(c, BadGateway(message))
}

// RespondBadRequest creates and responds with a BadRequest error.
func RespondBadRequest(c *gin.Context, message string) {
	Respond(c, BadRequest(message))
}

// RespondConflict creates and responds with a Conflict error.
func RespondConflict(c *gin.Context, message string) {
	Respond(c, Conflict(message))
}

// RespondForbidden creates and responds with a Forbidden error.
func RespondForbidden(c *gin.Context, message string) {
	Respond(c, Forbidden(message))
}

// RespondInternalServer creates and responds with an InternalServer error.
func RespondInternalServer(c *gin.Context, message string) {
	Respond(c, InternalServer(message))
}

// RespondMethodNotAllowed creates and responds with a MethodNotAllowed error.
func RespondMethodNotAllowed(c *gin.Context, message string) {
	Respond(c, MethodNotAllowed(message))
}

// RespondNotFound creates and responds with a NotFound error.
func RespondNotFound(c *gin.Context, message string) {
	Respond(c, NotFound(message))
}

// RespondNotImplemented creates and responds with a NotImplemented error.
func RespondNotImplemented(c *gin.Context, message string) {
	Respond(c, NotImplemented(message))
}

// RespondServiceUnavailable creates and responds with a ServiceUnavailable error.
func RespondServiceUnavailable(c *gin.Context, message string) {
	Respond(c, ServiceUnavailable(message))
}

// RespondTimeout creates and responds with a Timeout error.
func RespondTimeout(c *gin.Context, message string) {
	Respond(c, Timeout(message))
}

// RespondTooManyRequests creates and responds with a TooManyRequests error.
func RespondTooManyRequests(c *gin.Context, message string) {
	Respond(c, TooManyRequests(message))
}

// RespondUnauthorized creates and responds with an Unauthorized error.
func RespondUnauthorized(c *gin.Context, message string) {
	Respond(c, Unauthorized(message))
}

// RespondUnprocessableEntity creates and responds with an UnprocessableEntity error.
func RespondUnprocessableEntity(c *gin.Context, message string) {
	Respond(c, UnprocessableEntity(message))
}
