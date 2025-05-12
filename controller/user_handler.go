package controller

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/controller/middleware"
	"github.com/amirdashtii/go_auth/controller/validators"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/amirdashtii/go_auth/internal/core/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHTTPHandler struct {
	svc    ports.UserService
	logger ports.Logger
}

func NewUserHTTPHandler() *UserHTTPHandler {
	svc := service.NewUserService()

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

	return &UserHTTPHandler{
		svc:    svc,
		logger: appLogger,
	}
}

func NewUserRoutes(r *gin.Engine) {
	h := NewUserHTTPHandler()

	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.GET("/me", h.GetUserProfileHandler)
	userGroup.PUT("/me", h.UpdateUserProfileHandler)
	userGroup.PUT("/me/change-password", h.ChangePasswordHandler)
	userGroup.DELETE("/me", h.DeleteUserProfileHandler)
}

func (h *UserHTTPHandler) GetUserProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling get profile request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
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

	profile, err := h.svc.GetProfile(ctx, &userIDUUID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *UserHTTPHandler) UpdateUserProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling update profile request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
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

	var req dto.UserUpdateRequest
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

	err = validators.ValidateUserUpdateRequest(&req, h.logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := h.svc.UpdateProfile(ctx, &userIDUUID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
	})
}

func (h *UserHTTPHandler) ChangePasswordHandler(c *gin.Context) {
	ctx := c.Request.Context()
	if ctx.Err() != nil {
		h.logger.Error("Context cancelled while handling change password request",
			ports.F("error", ctx.Err()),
			ports.F("path", c.Request.URL.Path),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.ErrContextCancelled.ErrorPersian(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
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

	var req dto.ChangePasswordRequest
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
	err = validators.ValidateChangePasswordRequest(&req, h.logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := h.svc.ChangePassword(ctx, &userIDUUID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

func (h *UserHTTPHandler) DeleteUserProfileHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User not authenticated",
			ports.F("error", errors.ErrUserNotAuthenticated.Message.English),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": errors.ErrUserNotAuthenticated,
		})
		return
	}

	userIDUUID, err := uuid.Parse(userID.(string))
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

	if err := h.svc.DeleteProfile(ctx, &userIDUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted successfully",
	})
}
