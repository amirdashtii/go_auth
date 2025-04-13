package controller

import (
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

	userGroup := r.Group("/user")
	userGroup.GET("/profile", h.GetUserProfileHandler)
	userGroup.PATCH("/profile", h.UpdateUserProfileHandler)
	userGroup.PATCH("/change-password", h.ChangePasswordHandler)
}


func (h *UserHTTPHandler) GetUserProfileHandler(c *gin.Context){
	// TODO: Implement user profile retrieval logic
}
func (h *UserHTTPHandler) UpdateUserProfileHandler(c *gin.Context){
	// TODO: Implement user profile update logic
}
func (h *UserHTTPHandler) ChangePasswordHandler(c *gin.Context){
	// TODO: Implement password change logic with old password validation
}
