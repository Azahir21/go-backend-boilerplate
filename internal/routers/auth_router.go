package routers

import (
	"github.com/azahir21/go-backend-boilerplate/internal/container"
	"github.com/azahir21/go-backend-boilerplate/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
    // Get container instance
    c := container.GetInstance()
    
    // Get handler from container
    h := c.Handlers.AuthHandler
    
    // Public routes
    api := r.Group("/api/v1")
    {
        api.GET("/ping", h.Ping)
        
        // Auth routes
        auth := api.Group("/auth")
        {
            auth.POST("/register", h.Register)
            auth.POST("/login", h.Login)
            auth.GET("/profile", middleware.AuthMiddleware(), h.GetProfile)
        }
    }
    
    // Protected routes
    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())
    {
        // Admin routes
        admin := protected.Group("/admin")
        admin.Use(middleware.AdminMiddleware())
        {
            admin.GET("/test", h.AdminOnly)
        }
    }
}