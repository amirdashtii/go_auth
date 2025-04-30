package ports

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID *uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(userID *uuid.UUID, updateReq *dto.UserUpdateRequest) error
	ChangePassword(userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error
	DeleteProfile(userID *uuid.UUID) error
}
