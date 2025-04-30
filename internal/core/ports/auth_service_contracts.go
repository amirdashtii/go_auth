package ports

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
)

type AuthService interface {
	Register(registerReq *dto.RegisterRequest) error
	Login(loginReq *dto.LoginRequest) (*entities.TokenPair, error)
	Logout(userID string) error
	RefreshToken(refreshToken string) (*entities.TokenPair, error)
	ValidateToken(userID, token string) error
}
