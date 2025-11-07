package app

import "github.com/gin-gonic/gin"

// HttpRouter is an interface for HTTP routers that can register routes.
type HttpRouter interface {
	RegisterRoutes(group *gin.RouterGroup)
}

func NewServer(httpRouters ...HttpRouter) *gin.Engine {
	engine := gin.Default()

	// API Versioning
	v1 := engine.Group("/api/v1")
	{
		for _, router := range httpRouters {
			router.RegisterRoutes(v1)
		}
	}

	return engine
}