package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"

	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/errors"
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
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateRegisterRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	err := h.svc.Register(&req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHTTPHandler) LoginHandler(c *gin.Context) {
	var req *dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateLoginRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	tokens, err := h.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (h *AuthHTTPHandler) LogoutHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrInvalidUserIDType,
		})
		return
	}

	err := h.svc.Logout(userIDStr)

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
	var req dto.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateRefreshTokenRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})	 
		return
	}

	tokens, err := h.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}
