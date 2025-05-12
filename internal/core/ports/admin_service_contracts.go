package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AdminService interface {
	GetUsers(ctx context.Context, status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error)
	AdminGetUserByID(ctx context.Context, userID *uuid.UUID) (*dto.AdminUserResponse, error)
	AdminUpdateUser(ctx context.Context, userID *uuid.UUID, updateReq *dto.AdminUserUpdateRequest) error
	ChangeUserRole(ctx context.Context, userID *uuid.UUID, updateRole *entities.RoleType) error
	ChangeUserStatus(ctx context.Context, userID *uuid.UUID, updateStatus *entities.StatusType) error
	AdminDeleteUser(ctx context.Context, userID *uuid.UUID) error
}
