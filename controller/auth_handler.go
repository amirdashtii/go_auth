package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"

	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
)

type AuthHTTPHandler struct {
	svc ports.AuthService
}

func NewAuthHTTPHandler() *AuthHTTPHandler {
	svc := service.NewAuthService()
	return &AuthHTTPHandler{
		svc: svc,
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
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateRegisterRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

	err := h.svc.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Registration failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHTTPHandler) LoginHandler(c *gin.Context) {
	var req *dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateLoginRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

	tokens, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Login failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"tokens": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}

func (h *AuthHTTPHandler) LogoutHandler(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID type",
		})
		return
	}

	err := h.svc.Logout(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Logout failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func (h *AuthHTTPHandler) RefreshTokenHandler(c *gin.Context) {
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateRefreshTokenRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

	tokens, err := h.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Token refresh failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens": gin.H{
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		},
	})
}
