package ports

import (
	"time"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Create(user *entities.User) error
	FindByEmail(email *string) (*entities.User, error)
	FindAuthUserByID(id uuid.UUID) (*entities.User, error)
}
type UserRepository interface {
	FindUserByID(id *uuid.UUID) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id *uuid.UUID) error
}

type AdminRepository interface {
	FindUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]*entities.User, error)
	AdminGetUserByID(id *uuid.UUID) (*entities.User, error)
	AdminUpdateUser(user *entities.User) error
	AdminChangeUserRole(id *uuid.UUID, role *entities.RoleType) error
	AdminChangeUserStatus(id *uuid.UUID, status *entities.StatusType) error
	AdminDeleteUser(id *uuid.UUID) error
}
type InMemoryRespositoryContracts interface {
	AddToken(userID, token string, expiration time.Duration) error
	RemoveToken(userID string) error
	FindToken(userID string) (string, error)
}
