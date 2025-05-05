package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHTTPHandler struct {
	svc ports.UserService
}

func NewUserHTTPHandler() *UserHTTPHandler {
	svc := service.NewUserService()
	return &UserHTTPHandler{
		svc: svc,
	}
}

func NewUserRoutes(r *gin.Engine) {
	h := NewUserHTTPHandler()

	userGroup := r.Group("/profile")
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.GET("/", h.GetUserProfileHandler)
	userGroup.PUT("/", h.UpdateUserProfileHandler)
	userGroup.PUT("/change-password", h.ChangePasswordHandler)
	userGroup.DELETE("/", h.DeleteUserProfileHandler)
}

func (h *UserHTTPHandler) GetUserProfileHandler(c *gin.Context) {
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

	uuid, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	resp, err := h.svc.GetProfile(&uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHTTPHandler) UpdateUserProfileHandler(c *gin.Context) {
	var updateReq dto.UserUpdateRequest

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateUserUpdateRequest(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

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

	uuid, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.svc.UpdateProfile(&uuid, &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *UserHTTPHandler) ChangePasswordHandler(c *gin.Context) {
	var changePasswordReq dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&changePasswordReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := validators.ValidateChangePasswordRequest(&changePasswordReq); err != nil {
		c.JSON(http.StatusBadRequest, validators.HandleValidationError(err))
		return
	}

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

	uuid, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.svc.ChangePassword(&uuid, &changePasswordReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to change password",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *UserHTTPHandler) DeleteUserProfileHandler(c *gin.Context) {
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
	uuid, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	err = h.svc.DeleteProfile(&uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete profile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}
