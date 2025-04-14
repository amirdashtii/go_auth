package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/entities"
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
	authGroup.POST("/logout", h.LogoutHandler)
	authGroup.POST("/refresh-token", h.RefreshTokenHandler)
}

func (h *AuthHTTPHandler) RegisterHandler(c *gin.Context) {
	var user entities.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

	err := h.svc.Register(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Registration failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHTTPHandler) LoginHandler(c *gin.Context) {
	// TODO: Implement login logic with email and password validation
}

func (h *AuthHTTPHandler) LogoutHandler(c *gin.Context) {
	// TODO: Implement logout logic and token invalidation
}

func (h *AuthHTTPHandler) RefreshTokenHandler(c *gin.Context) {
	// TODO: Implement token refresh logic
}