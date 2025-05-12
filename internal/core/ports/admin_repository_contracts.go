package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AdminRepository interface {
	FindUsers(ctx context.Context, status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]entities.User, error)
	AdminGetUserByID(ctx context.Context, id *uuid.UUID) (*entities.User, error)
	AdminUpdateUser(ctx context.Context, user *entities.User) error
	AdminChangeUserRole(ctx context.Context, id *uuid.UUID, role *entities.RoleType) error
	AdminChangeUserStatus(ctx context.Context, id *uuid.UUID, status *entities.StatusType) error
	AdminDeleteUser(ctx context.Context, id *uuid.UUID) error
}
