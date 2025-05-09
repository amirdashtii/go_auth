package ports

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AdminRepository interface {
	FindUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]entities.User, error)
	AdminGetUserByID(id *uuid.UUID) (*entities.User, error)
	AdminUpdateUser(user *entities.User) error
	AdminChangeUserRole(id *uuid.UUID, role *entities.RoleType) error
	AdminChangeUserStatus(id *uuid.UUID, status *entities.StatusType) error
	AdminDeleteUser(id *uuid.UUID) error
}
