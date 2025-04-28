package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	adminGroup.PUT("/users/:id/role", h.ChangeUserRoleHandler)
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

	status := entities.ParseStatusType(c.DefaultQuery("status", "active"))
	roleFilter := entities.ParseRoleType(c.DefaultQuery("role", "user"))
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	users, err := h.svc.GetUsers(&status, &roleFilter, &sort, &order)
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

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.svc.AdminGetUserByID(&userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	resp := AdminUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    user.Status.String(),
		Role:      user.Role.String(),
	}
	c.JSON(200, resp)
}

func (h *AdminHTTPHandler) UpdateUserHandler(c *gin.Context) {
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

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	user := &entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	err = h.svc.AdminUpdateUser(&userID, user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

func (h *AdminHTTPHandler) 	ChangeUserRoleHandler(c *gin.Context) {
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

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AdminUserUpdateRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	reqRole := entities.ParseRoleType(req.Role)
	err = h.svc.ChangeUserRole(&userID, &reqRole)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to change user role", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User role changed successfully"})
}

func (h *AdminHTTPHandler) ChangeUserStatusHandler(c *gin.Context) {
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

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AdminUserUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	status := entities.ParseStatusType(req.Status)
	err = h.svc.ChangeUserStatus(&userID, &status)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to change user status", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User status changed successfully"})
}

func (h *AdminHTTPHandler) DeleteUserHandler(c *gin.Context) {
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

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.svc.AdminDeleteUser(&userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}
