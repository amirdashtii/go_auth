package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(ctx context.Context, userID *uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID *uuid.UUID, updateReq *dto.UserUpdateRequest) error
	ChangePassword(ctx context.Context, userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error
	DeleteProfile(ctx context.Context, userID *uuid.UUID) error
}
