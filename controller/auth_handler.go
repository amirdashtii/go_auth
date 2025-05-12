package controller

import (
	"context"
	"net/http"
	"os"
	"time"

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

	// Create log file
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      logFile,
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

func (h *AuthHTTPHandler) RegisterHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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

func (h *AuthHTTPHandler) LoginHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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

func (h *AuthHTTPHandler) LogoutHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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

func (h *AuthHTTPHandler) RefreshTokenHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

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
