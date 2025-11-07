package http

import (
	"net/http"

	"github.com/azahir21/go-backend-boilerplate/internal/domain"
	"github.com/azahir21/go-backend-boilerplate/internal/middleware"
	"github.com/azahir21/go-backend-boilerplate/internal/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/httpresp"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
    Log         *logrus.Logger
    UserUsecase usecase.UserUsecase
}

func NewUserHandler(log *logrus.Logger, userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        Log:         log,
        UserUsecase: userUsecase,
    }
}

// RegisterRoutes implements the HttpRouter interface.
func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/ping", h.Ping)

	// Auth routes
	auth := group.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/profile", middleware.AuthMiddleware(), h.GetProfile)
	}

	// Protected routes
	protected := group.Group("/")
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

// Ping godoc
// @Summary Ping endpoint
// @Description Health check endpoint
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ping [get]
func (h *UserHandler) Ping(c *gin.Context) {
    httpresp.JSON(c, http.StatusOK, "pong", gin.H{"message": "pong"})
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Registration request"
// @Success 201 {object} httpresp.Response{data=domain.AuthResponse}
// @Failure 400 {object} httpresp.Response
// @Failure 500 {object} httpresp.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpresp.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    result, err := h.UserUsecase.Register(c.Request.Context(), &req)
    if err != nil {
        httpresp.JSON(c, http.StatusBadRequest, err.Error(), nil)
        return
    }

    httpresp.JSON(c, http.StatusCreated, "User registered successfully", result)
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login request"
// @Success 200 {object} httpresp.Response{data=domain.AuthResponse}
// @Failure 400 {object} httpresp.Response
// @Failure 401 {object} httpresp.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        httpresp.JSON(c, http.StatusBadRequest, "Invalid request body", nil)
        return
    }

    result, err := h.UserUsecase.Login(c.Request.Context(), &req)
    if err != nil {
        httpresp.JSON(c, http.StatusUnauthorized, err.Error(), nil)
        return
    }

    httpresp.JSON(c, http.StatusOK, "Login successful", result)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpresp.Response{data=domain.User}
// @Failure 401 {object} httpresp.Response
// @Failure 404 {object} httpresp.Response
// @Router /auth/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        httpresp.JSON(c, http.StatusUnauthorized, "User not authenticated", nil)
        return
    }

    user, err := h.UserUsecase.GetProfile(c.Request.Context(), userID.(uint))
    if err != nil {
        httpresp.JSON(c, http.StatusNotFound, "User not found", nil)
        return
    }

    httpresp.JSON(c, http.StatusOK, "Profile retrieved successfully", user)
}

// AdminOnly godoc
// @Summary Admin only endpoint ??
// @Description Example endpoint that requires admin access
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpresp.Response
// @Failure 401 {object} httpresp.Response
// @Failure 403 {object} httpresp.Response
// @Router /admin/test [get]
func (h *UserHandler) AdminOnly(c *gin.Context) {
    httpresp.JSON(c, http.StatusOK, "Admin access granted", gin.H{
        "message": "This is an admin-only endpoint",
        "user":    c.GetString("username"),
    })
}