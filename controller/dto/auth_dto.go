package dto

// RegisterRequest is used for user registration
// swagger:model
type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Password    string `json:"password" binding:"required" validate:"password,min=8"`
}

// LoginRequest is used for user login
// swagger:model
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Password    string `json:"password" binding:"required" validate:"password,min=8"`
}

// RefreshTokenRequest is used for refreshing JWT token
// swagger:model
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}
