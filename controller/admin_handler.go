package controller

import (
	"github.com/amirdashtii/go_auth/controller/middleware"
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
	adminGroup.Use(middleware.AdminMiddleware())
	
	adminGroup.GET("/users", h.GetUsersHandler)
	adminGroup.GET("/users/:id", h.GetUserByIDHandler)
	adminGroup.PUT("/users/:id", h.UpdateUserHandler)
	adminGroup.POST("/users/:id/promote", h.PromoteToAdminHandler)
	adminGroup.POST("/users/:id/deactivate", h.DeactivateUserHandler)
	adminGroup.POST("/users/:id/activate", h.ActivateUserHandler)
	adminGroup.DELETE("/users/:id", h.DeleteUserHandler)
	adminGroup.GET("/users/active", h.FindActiveUsersHandler)
	adminGroup.GET("/users/admins", h.FindAdminsHandler)
}




func (h *AdminHTTPHandler) GetUsersHandler(c *gin.Context){
	// TODO: Implement user list retrieval with pagination
}
func (h *AdminHTTPHandler) GetUserByIDHandler(c *gin.Context){
	// TODO: Implement user retrieval by ID
}
func (h *AdminHTTPHandler) UpdateUserHandler(c *gin.Context){
	// TODO: Implement user update by ID
}
func (h *AdminHTTPHandler) PromoteToAdminHandler(c *gin.Context){
	// TODO: Implement user promotion to admin role
}
func (h *AdminHTTPHandler) DeactivateUserHandler(c *gin.Context){
	// TODO: Implement user deactivation
}
func (h *AdminHTTPHandler) ActivateUserHandler(c *gin.Context){
	// TODO: Implement user activation
}
func (h *AdminHTTPHandler) DeleteUserHandler(c *gin.Context){
	// TODO: Implement user deletion (soft delete)
}
func (h *AdminHTTPHandler) FindActiveUsersHandler(c *gin.Context){
	// TODO: Implement finding active users
}
func (h *AdminHTTPHandler) FindAdminsHandler(c *gin.Context){
	// TODO: Implement finding admins
}