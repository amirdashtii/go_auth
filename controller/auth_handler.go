package controller

import (
	"context"
	"net/http"
	"os"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/infrastructure/logger"

	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
)

type AuthHTTPHandler struct {
	svc    ports.AuthService
	logger ports.Logger
}

func NewAuthHTTPHandler() *AuthHTTPHandler {
	svc := service.NewAuthService()

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      os.Stdout,
	}

	appLogger := logger.NewZerologLogger(loggerConfig)

	return &AuthHTTPHandler{
		svc:    svc,
		logger: appLogger,
	}
}

func NewAuthRoutes(r *gin.Engine) {
	h := NewAuthHTTPHandler()

	authGroup := r.Group("/auth")
	authGroup.POST("/register", h.RegisterHandler)
	authGroup.POST("/login", h.LoginHandler)
	authGroup.POST("/logout", middleware.AuthMiddleware(), h.LogoutHandler)
	authGroup.POST("/refresh-token", h.RefreshTokenHandler)
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Register a new user with phone number and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHTTPHandler) RegisterHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling register request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("request", req),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateRegisterRequest(&req, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	err := h.svc.Register(ctx, &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginHandler godoc
// @Summary Login user
// @Description Login user with phone number and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHTTPHandler) LoginHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling login request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	var req *dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("request", req),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateLoginRequest(req, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.svc.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

// LogoutHandler godoc
// @Summary Logout user
// @Description Logout user and invalidate tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHTTPHandler) LogoutHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling logout request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid user ID type",
			ports.F("error", errors.ErrInvalidUserIDType.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrInvalidUserIDType,
		})
		return
	}

	err := h.svc.Logout(ctx, userIDStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// RefreshTokenHandler godoc
// @Summary Refresh access token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/refresh-token [post]
func (h *AuthHTTPHandler) RefreshTokenHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling refresh token request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("request", req),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateRefreshTokenRequest(&req, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}
