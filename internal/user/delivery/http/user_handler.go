package http

import (
	"net/http"

	api "github.com/azahir21/go-backend-boilerplate/internal/shared/http"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/middleware"
	"github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http/dto"
	"github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
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
func (h *UserHandler) RegisterRoutes(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	grp := api.NewAPIRouterGroup(v1)
	authGrp := api.NewAPIRouterGroup(v1.Group("/auth"))
	adminGrp := api.NewAPIRouterGroup(v1.Group("/admin"))

	// Public + mixed endpoints
	grp.Register(
		api.EndpointSpec{Method: http.MethodGet, Path: "/ping", Handler: h.Ping},
	)

	// Auth endpoints with typed handlers (auto-binding)
	authGrp.Register(
		api.EndpointSpec{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.Register,
		},
		api.EndpointSpec{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: h.Login,
		},
		api.EndpointSpec{
			Method:      http.MethodGet,
			Path:        "/profile",
			Handler:     h.GetProfile,
			Middlewares: []gin.HandlerFunc{middleware.AuthMiddleware()},
		},
	)

	// Admin endpoints
	adminGrp.Register(
		api.EndpointSpec{
			Method:      http.MethodGet,
			Path:        "/test",
			Handler:     h.AdminOnly,
			Middlewares: []gin.HandlerFunc{middleware.AuthMiddleware(), middleware.AdminMiddleware()},
		},
	)
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
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} httpresp.Response{data=dto.AuthResponse}
// @Failure 400 {object} httpresp.Response
// @Failure 500 {object} httpresp.Response
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context, req *dto.RegisterRequest) error {
	result, err := h.UserUsecase.Register(c.Request.Context(), req)
	if err != nil {
		httpresp.JSON(c, http.StatusBadRequest, err.Error(), nil)
		return nil
	}

	httpresp.JSON(c, http.StatusCreated, "User registered successfully", result)
	return nil
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} httpresp.Response{data=dto.AuthResponse}
// @Failure 400 {object} httpresp.Response
// @Failure 401 {object} httpresp.Response
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context, req *dto.LoginRequest) error {
	result, err := h.UserUsecase.Login(c.Request.Context(), req)
	if err != nil {
		httpresp.JSON(c, http.StatusUnauthorized, err.Error(), nil)
		return nil
	}

	httpresp.JSON(c, http.StatusOK, "Login successful", result)
	return nil
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpresp.Response{data=entity.User}
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
