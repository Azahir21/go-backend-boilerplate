//go:build rest
// +build rest

package module

import (
	sharedHttp "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
)

// HTTPModule is implemented by modules that provide HTTP handlers.
type HTTPModule interface {
	Module
	// HTTPHandler returns the HTTP router for this module.
	HTTPHandler() sharedHttp.HttpRouter
}
