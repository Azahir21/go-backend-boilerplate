package http

import (
	"github.com/gin-gonic/gin"
)

// HttpRouter is an interface for HTTP routers that can register routes.
type HttpRouter interface {
	RegisterRoutes(engine *gin.Engine)
}

// NewServer creates a new Gin engine and registers routes from provided HttpRouters.
func NewServer(httpRouters ...HttpRouter) *gin.Engine {
	engine := gin.Default()

	for _, router := range httpRouters {
		router.RegisterRoutes(engine)
	}

	return engine
}
