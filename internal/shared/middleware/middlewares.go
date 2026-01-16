package middleware

import (
	"strings"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/helper"
	"github.com/azahir21/go-backend-boilerplate/pkg/apperr"
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
			apperr.RespondUnauthorized(c, "Authorization header missing")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			apperr.RespondUnauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			apperr.RespondUnauthorized(c, "Token missing in authorization header")
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(tokenString)
		if err != nil {
			apperr.RespondUnauthorized(c, "Invalid or expired token")
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
			apperr.RespondUnauthorized(c, "User not authenticated")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			apperr.RespondForbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
