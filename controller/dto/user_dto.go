package dto

type UserProfileResponse struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
}

// UserUpdateRequest is used for updating user profile
type UserUpdateRequest struct {
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone"`
	FirstName   string `json:"first_name" validate:"omitempty,name"`
	LastName    string `json:"last_name" validate:"omitempty,name"`
	Email       string `json:"email" validate:"omitempty,email"`
}

// ChangePasswordRequest is used for changing user password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
}
