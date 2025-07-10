package middleware

import (
	"net/http"
	"strings"

	"github.com/azahir21/go-backend-boilerplate/internal/helper"
	"github.com/azahir21/go-backend-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.JSON(c, http.StatusUnauthorized, "Authorization header is required", nil)
            c.Abort()
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := helper.ValidateToken(tokenString)
        if err != nil {
            response.JSON(c, http.StatusUnauthorized, "Invalid token", nil)
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != "admin" {
            response.JSON(c, http.StatusForbidden, "Admin access required", nil)
            c.Abort()
            return
        }
        c.Next()
    }
}