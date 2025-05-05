package dto

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Password    string `json:"password" binding:"required" validate:"password,min=8"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" validate:"phone"`
	Password    string `json:"password" binding:"required" validate:"password,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" validate:"required"`
}
