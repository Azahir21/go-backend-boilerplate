package handler

import (
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/internal/model"
	"github.com/azahir21/go-backend-boilerplate/internal/service"
	"github.com/azahir21/go-backend-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthHandler struct {
    DB          *gorm.DB
    Log         *logrus.Logger
    AuthService service.AuthService
}

func NewHandler(db *gorm.DB, log *logrus.Logger, authService service.AuthService) *AuthHandler {
    return &AuthHandler{
        DB:          db,
        Log:         log,
        AuthService: authService,
    }
}

// Ping godoc
// @Summary Ping endpoint
// @Description Health check endpoint
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *AuthHandler) Ping(c *gin.Context) {
    response.JSON(c, http.StatusOK, "pong", gin.H{"message": "pong"})
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration request"
// @Success 201 {object} response.Response{data=model.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    var req model.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    result, err := h.AuthService.Register(&req)
    if err != nil {
        response.JSON(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    response.JSON(c, http.StatusCreated, "User registered successfully", result)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=model.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req model.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    result, err := h.AuthService.Login(&req)
    if err != nil {
        response.JSON(c, http.StatusUnauthorized, err.Error(), nil)
        return
    }

    response.JSON(c, http.StatusOK, "Login successful", result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        response.JSON(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }

    user, err := h.AuthService.GetProfile(userID.(uint))
    if err != nil {
        response.JSON(c, http.StatusNotFound, "User not found", nil)
        return
    }

    response.JSON(c, http.StatusOK, "Profile retrieved successfully", user)
}

// AdminOnly godoc
// @Summary Admin only endpoint ??
// @Description Example endpoint that requires admin access
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /admin/test [get]
func (h *AuthHandler) AdminOnly(c *gin.Context) {
    response.JSON(c, http.StatusOK, "Admin access granted", gin.H{
        "message": "This is an admin-only endpoint",
        "user":    c.GetString("username"),
    })
}