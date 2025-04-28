package ports

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(user *entities.User) error
	Login(email, password string) (*entities.TokenPair, error)
	Logout(userID string) error
	RefreshToken(refreshToken string) (*entities.TokenPair, error)
	ValidateToken(userID, token string) error
}

type UserService interface {
	GetProfile(userID uuid.UUID) (*entities.User, error)
	UpdateProfile(userID uuid.UUID, user *entities.User) error
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
	DeleteProfile(userID uuid.UUID) error
}

type AdminService interface {
	GetUsers(status, role, sort, order string) ([]*entities.User, error)
	AdminGetUserByID(userID uuid.UUID) (*entities.User, error)
	AdminUpdateUser(userID uuid.UUID, user *entities.User) error
	PromoteToAdmin(userID uuid.UUID) error
	DeactivateUser(userID uuid.UUID) error
	ActivateUser(userID uuid.UUID) error
	AdminDeleteUser(userID uuid.UUID) error
	FindActiveUsers() ([]*entities.User, error)
	FindAdmins() ([]*entities.User, error)
}
