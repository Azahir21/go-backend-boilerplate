package http

import (
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/pkg/httpresp"
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

// BindJSON is a helper function to bind JSON requests and handle errors.
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		httpresp.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
		return false
	}
	return true
}
