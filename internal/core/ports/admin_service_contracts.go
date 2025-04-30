package ports

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AdminService interface {
	GetUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error)
	AdminGetUserByID(userID *uuid.UUID) (*dto.AdminUserResponse, error)
	AdminUpdateUser(userID *uuid.UUID, updateReq *dto.AdminUserUpdateRequest) error
	ChangeUserRole(userID *uuid.UUID, updateRole *entities.RoleType) error
	ChangeUserStatus(userID *uuid.UUID, updateStatus *entities.StatusType) error
	AdminDeleteUser(userID *uuid.UUID) error
}
