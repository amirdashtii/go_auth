package controller

import (
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
	adminGroup.GET("/users", h.GetUsersHandler)
	adminGroup.GET("/users/:user_id", h.GetUserByIdHandler)
	adminGroup.PATCH("/users/:user_id", h.UpdateUserByIdHandler)
	adminGroup.PATCH("/users/:user_id/promote", h.PromoteUserToAdminHandler)
	adminGroup.PATCH("/users/:user_id/deactivate", h.DeactivateUserHandler)
	adminGroup.PATCH("/users/:user_id/activate", h.ActivateUserHandler)
	adminGroup.DELETE("/users/:user_id", h.DeleteUserHandler)
}




func (h *AdminHTTPHandler) GetUsersHandler(c *gin.Context){
	// TODO: Implement user list retrieval with pagination
}
func (h *AdminHTTPHandler) GetUserByIdHandler(c *gin.Context){
	// TODO: Implement user retrieval by ID
}
func (h *AdminHTTPHandler) UpdateUserByIdHandler(c *gin.Context){
	// TODO: Implement user update by ID
}
func (h *AdminHTTPHandler) PromoteUserToAdminHandler(c *gin.Context){
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