package dto

import "time"

type AdminGetUsersRequest struct {
	Status string `json:"status" validate:"omitempty,status"`
	Role   string `json:"role" validate:"omitempty,role"`
	Sort   string `json:"sort" validate:"omitempty,sort"`
	Order  string `json:"order" validate:"omitempty,order"`
}

type AdminUserListResponse struct {
	Users []AdminUserResponse `json:"users"`
}

type AdminUserResponse struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AdminUserUpdateRequest is used for admin updating user profile
// swagger:model
type AdminUserUpdateRequest struct {
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone"`
	FirstName   string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName    string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email       string `json:"email" validate:"omitempty,email"`
}

// AdminUserUpdateRoleRequest is used for admin changing user role
// swagger:model
type AdminUserUpdateRoleRequest struct {
	Role string `json:"role" binding:"required,role"`
}

// AdminUserUpdateStatusRequest is used for admin changing user status
// swagger:model
type AdminUserUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,status"`
}
