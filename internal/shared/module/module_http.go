package module

import (
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
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
