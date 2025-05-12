package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
)

type AuthService interface {
	Register(ctx context.Context, registerReq *dto.RegisterRequest) error
	Login(ctx context.Context, loginReq *dto.LoginRequest) (*entities.TokenPair, error)
	Logout(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*entities.TokenPair, error)
	ValidateToken(ctx context.Context, userID, token string) error
}
