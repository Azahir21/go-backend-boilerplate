package httpresp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard JSON response structure for the API
type Response struct {
	Status  string      `json:"status"`            // HTTP status text (e.g., "OK", "Bad Request")
	Message string      `json:"message"`           // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // Response payload (optional)
}

// JSON sends a JSON response with a standard structure.
// It automatically sets the HTTP status text based on the status code.
func JSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  http.StatusText(status),
		Message: message,
		Data:    data,
	})
}

// Success sends a 200 OK response with the given data
func Success(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, message, data)
}

// Created sends a 201 Created response with the given data
func Created(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusCreated, message, data)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	JSON(c, http.StatusBadRequest, message, nil)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	JSON(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	JSON(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, message, nil)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string) {
	JSON(c, http.StatusInternalServerError, message, nil)
}
