package controller

import (
	"net/http"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
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
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	req := dto.AdminGetUsersRequest{
		Status: c.DefaultQuery("status", "active"),
		Role:   c.DefaultQuery("role", "user"),
		Sort:   c.DefaultQuery("sort", "created_at"),
		Order:  c.DefaultQuery("order", "desc"),
	}

	if err := validators.ValidateGetUsersRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	statuesType := entities.ParseStatusType(req.Status)
	roleType := entities.ParseRoleType(req.Role)

	resp, err := h.svc.GetUsers(&statuesType, &roleType, &req.Sort, &req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto.AdminUserListResponse{Users: resp}})
}

func (h *AdminHTTPHandler) GetUserByIDHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", nil)})
		return
	}

	resp, err := h.svc.AdminGetUserByID(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func (h *AdminHTTPHandler) UpdateUserHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", err)})
		return
	}

	var updateReq dto.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid request", "درخواست نامعتبر است", err)})
		return
	}

	if err := validators.ValidateUpdateUserRequest(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	err = h.svc.AdminUpdateUser(&userID, &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *AdminHTTPHandler) ChangeUserRoleHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", err)})
		return
	}

	var updateRoleReq dto.AdminUserUpdateRoleRequest
	if err := c.ShouldBindJSON(&updateRoleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid request", "درخواست نامعتبر است", err)})
		return
	}

	if err := validators.ValidateChangeRoleRequest(&updateRoleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	updateRole := entities.ParseRoleType(updateRoleReq.Role)
	err = h.svc.ChangeUserRole(&userID, &updateRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role changed successfully"})
}

func (h *AdminHTTPHandler) ChangeUserStatusHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", err)})
		return
	}

	var updateStatusReq dto.AdminUserUpdateStatusRequest
	if err := c.ShouldBindJSON(&updateStatusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid request", "درخواست نامعتبر است", err)})
		return
	}

	if err := validators.ValidateChangeStatusRequest(&updateStatusReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	updateStatus := entities.ParseStatusType(updateStatusReq.Status)
	err = h.svc.ChangeUserStatus(&userID, &updateStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status changed successfully"})
}

func (h *AdminHTTPHandler) DeleteUserHandler(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.New(errors.AuthenticationError, "User ID not found", "شناسه کاربر یافت نشد", nil),
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		c.JSON(http.StatusForbidden, gin.H{"error": errors.New(errors.AuthorizationError, "Forbidden", "دسترسی غیرمجاز", nil)})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errors.ValidationError, "Invalid user ID", "شناسه کاربر نامعتبر است", err)})
		return
	}

	err = h.svc.AdminDeleteUser(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
