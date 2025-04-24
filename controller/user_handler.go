package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"

	"github.com/gin-gonic/gin"
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
	userGroup.POST("/change-password", h.ChangePasswordHandler)
}

func (h *UserHTTPHandler) GetUserProfileHandler(c *gin.Context) {
	// TODO: Implement user profile retrieval logic
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

	userProfile, err := h.svc.GetOwnProfile(userIDStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user profile",
		})
		return
	}
	c.JSON(http.StatusOK, userProfile)
}
func (h *UserHTTPHandler) UpdateUserProfileHandler(c *gin.Context) {
	// TODO: Implement user profile update logic
}
func (h *UserHTTPHandler) ChangePasswordHandler(c *gin.Context) {
	// TODO: Implement password change logic with old password validation
}
