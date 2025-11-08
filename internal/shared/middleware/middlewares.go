package middleware

import (
	"net/http"
	"strings"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/helper"
	"github.com/azahir21/go-backend-boilerplate/pkg/httpresp"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDKey           = "user_id"
	usernameKey         = "username"
	roleKey             = "role"
)

// AuthMiddleware validates JWT tokens and sets user information in the context.
// It expects the Authorization header with format: "Bearer <token>".
// On success, it sets user_id, username, and role in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(authorizationHeader)
		if authHeader == "" {
			httpresp.JSON(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			httpresp.JSON(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			httpresp.JSON(c, http.StatusUnauthorized, "Token is empty", nil)
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(tokenString)
		if err != nil {
			httpresp.JSON(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Set user information in context
		c.Set(userIDKey, claims.UserID)
		c.Set(usernameKey, claims.Username)
		c.Set(roleKey, claims.Role)
		c.Next()
	}
}

// AdminMiddleware checks if the authenticated user has admin role.
// This middleware should be used after AuthMiddleware.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(roleKey)
		if !exists {
			httpresp.JSON(c, http.StatusUnauthorized, "User not authenticated", nil)
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			httpresp.JSON(c, http.StatusForbidden, "Admin access required", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
