package controller

import (
	"context"
	"net/http"
	"os"

	// "time" // This line should be removed or commented out

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHTTPHandler struct {
	svc    ports.AdminService
	logger ports.Logger
}

func NewAdminHTTPHandler() *AdminHTTPHandler {
	svc := service.NewAdminService()

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
	return &AdminHTTPHandler{
		svc:    svc,
		logger: appLogger,
	}
}

func NewAdminRoutes(r *gin.Engine) {
	h := NewAdminHTTPHandler()

	usersGroup := r.Group("/users")
	usersGroup.Use(middleware.AuthMiddleware())

	usersGroup.GET("", h.GetUsersHandler)
	usersGroup.GET("/:id", h.GetUserByIDHandler)
	usersGroup.PUT("/:id", h.UpdateUserHandler)
	usersGroup.PUT("/:id/role", h.ChangeUserRoleHandler)
	usersGroup.PUT("/:id/status", h.ChangeUserStatusHandler)
	usersGroup.DELETE("/:id", h.DeleteUserHandler)
}

// GetUsersHandler godoc
// @Summary Get all users
// @Description Get list of all users (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *AdminHTTPHandler) GetUsersHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling get users request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden.Message.English),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	req := dto.AdminGetUsersRequest{
		Status: c.DefaultQuery("status", "active"),
		Role:   c.DefaultQuery("role", "user"),
		Sort:   c.DefaultQuery("sort", "created_at"),
		Order:  c.DefaultQuery("order", "desc"),
	}

	if err := validators.ValidateGetUsersRequest(&req, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	statuesType := entities.ParseStatusType(req.Status)
	roleType := entities.ParseRoleType(req.Role)

	resp, err := h.svc.GetUsers(ctx, &statuesType, &roleType, &req.Sort, &req.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dto.AdminUserListResponse{Users: resp}})
}

// GetUserByIDHandler godoc
// @Summary Get user by ID
// @Description Get user information by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *AdminHTTPHandler) GetUserByIDHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden.Message.English),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Invalid user ID",
			ports.F("error", errors.ErrInvalidUserID.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidUserID,
		})
		return
	}

	resp, err := h.svc.AdminGetUserByID(ctx, &userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// UpdateUserHandler godoc
// @Summary Update user
// @Description Update user information by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.AdminUserUpdateRequest true "Update User Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *AdminHTTPHandler) UpdateUserHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Invalid user ID",
			ports.F("error", errors.ErrInvalidUserID.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidUserID,
		})
		return
	}

	var updateReq dto.AdminUserUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("request", updateReq),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateUpdateUserRequest(&updateReq, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	err = h.svc.AdminUpdateUser(ctx, &userID, &updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// ChangeUserRoleHandler godoc
// @Summary Change user role
// @Description Change user role by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.AdminUserUpdateStatusRequest true "Change Status Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/role [put]
func (h *AdminHTTPHandler) ChangeUserRoleHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling change user role request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden.Message.English),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Invalid user ID",
			ports.F("error", errors.ErrInvalidUserID.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidUserID,
		})
		return
	}

	var updateRoleReq dto.AdminUserUpdateRoleRequest
	if err := c.ShouldBindJSON(&updateRoleReq); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateChangeRoleRequest(&updateRoleReq, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	updateRole := entities.ParseRoleType(updateRoleReq.Role)
	err = h.svc.ChangeUserRole(ctx, &userID, &updateRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role changed successfully"})
}

// ChangeUserStatusHandler godoc
// @Summary Change user status
// @Description Change user status by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body dto.AdminUserUpdateStatusRequest true "Change Status Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/status [put]
func (h *AdminHTTPHandler) ChangeUserStatusHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling change user status request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden.Message.English),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Invalid user ID",
			ports.F("error", errors.ErrInvalidUserID.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidUserID,
		})
		return
	}

	var updateStatusReq dto.AdminUserUpdateStatusRequest
	if err := c.ShouldBindJSON(&updateStatusReq); err != nil {
		h.logger.Error("Invalid request",
			ports.F("error", errors.ErrInvalidRequest.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidRequest,
		})
		return
	}

	if err := validators.ValidateChangeStatusRequest(&updateStatusReq, h.logger); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	updateStatus := entities.ParseStatusType(updateStatusReq.Status)
	err = h.svc.ChangeUserStatus(ctx, &userID, &updateStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status changed successfully"})
}

// DeleteUserHandler godoc
// @Summary Delete user
// @Description Delete user by ID (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *AdminHTTPHandler) DeleteUserHandler(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling delete user request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	roleStr := role.(string)
	if roleStr != entities.SuperAdminRole.String() && roleStr != entities.AdminRole.String() {
		h.logger.Error("User not authorized",
			ports.F("error", errors.ErrForbidden.Message.English),
		)
		c.JSON(http.StatusForbidden, gin.H{
			"error": errors.ErrForbidden,
		})
		return
	}

	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("Invalid user ID",
			ports.F("error", errors.ErrInvalidUserID.Message.English),
			ports.F("user_id", userID),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.ErrInvalidUserID,
		})
		return
	}

	err = h.svc.AdminDeleteUser(ctx, &userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
