package ports

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(email, password *string) (*entities.TokenPair, error)
	Logout(userID string) error
	RefreshToken(refreshToken string) (*entities.TokenPair, error)
	ValidateToken(userID, token string) error
}

type UserService interface {
	GetProfile(userID *uuid.UUID) (*dto.UserProfileResponse, error)
	UpdateProfile(userID *uuid.UUID, req *dto.UserUpdateRequest) error
	ChangePassword(userID *uuid.UUID, oldPassword, newPassword *string) error
	DeleteProfile(userID *uuid.UUID) error
}

type AdminService interface {
	GetUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error)
	AdminGetUserByID(userID *uuid.UUID) (*dto.AdminUserResponse, error)
	AdminUpdateUser(userID *uuid.UUID, req *dto.AdminUserUpdateRequest) error
	ChangeUserRole(userID *uuid.UUID, role *entities.RoleType) error
	ChangeUserStatus(userID *uuid.UUID, status *entities.StatusType) error
	AdminDeleteUser(userID *uuid.UUID) error
}
