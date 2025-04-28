package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
)

type AdminHTTPHandler struct {
	svc ports.AdminService
}

func NewAdminHTTPHandler() *AdminHTTPHandler {
	svc := service.NewAdminService()
	return &AdminHTTPHandler{
		svc: svc,
	}
}

func NewAdminRoutes(r *gin.Engine) {
	h := NewAdminHTTPHandler()

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware())

	adminGroup.GET("/users", h.GetUsersHandler)
	adminGroup.GET("/users/:id", h.GetUserByIDHandler)
	adminGroup.PUT("/users/:id", h.UpdateUserHandler)
	adminGroup.POST("/users/:id/role", h.ChangeUserRoleHandler)
	adminGroup.PUT("/users/:id/status", h.ChangeUserStatusHandler)
	adminGroup.DELETE("/users/:id", h.DeleteUserHandler)
}

func (h *AdminHTTPHandler) GetUsersHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found",
		})
		return
	}

	roleStr := role.(string)

	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	status := c.DefaultQuery("status", "active")
	roleFilter := c.DefaultQuery("role", "user")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	users, err := h.svc.GetUsers(status, roleFilter, sort, order)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve users", "details": err.Error()})
		return
	}

	var resp []AdminUserResponse
	for _, u := range users {
		resp = append(resp, AdminUserResponse{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Status:    u.Status.String(),
			Role:      u.Role.String(),
		})
	}

	c.JSON(200, AdminUserListResponse{Users: resp})
}

func (h *AdminHTTPHandler) GetUserByIDHandler(c *gin.Context) {
	// TODO: Implement user retrieval by ID
}

func (h *AdminHTTPHandler) UpdateUserHandler(c *gin.Context) {
	// TODO: Implement user update by ID
}

func (h *AdminHTTPHandler) ChangeUserRoleHandler(c *gin.Context) {
	// TODO: Implement user promotion to admin role
}

func (h *AdminHTTPHandler) ChangeUserStatusHandler(c *gin.Context) {
	// TODO: Implement user status change
}

func (h *AdminHTTPHandler) DeleteUserHandler(c *gin.Context) {
	// TODO: Implement user deletion (soft delete)
}
