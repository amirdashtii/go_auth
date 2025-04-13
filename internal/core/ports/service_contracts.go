package ports

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(user *entities.User) error
	Login(email, password string) (string, error)
	Logout(token string) error
	RefreshToken(token string) (string, error)
	ValidateToken(token string) (*entities.User, error)
}

type UserService interface {
	GetProfile(userID uuid.UUID) (*entities.User, error)
	UpdateProfile(userID uuid.UUID, user *entities.User) error
	ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error
}

type AdminService interface {
	GetUsers() ([]*entities.User, error)
	GetUserByID(userID uuid.UUID) (*entities.User, error)
	UpdateUser(userID uuid.UUID, user *entities.User) error
	PromoteToAdmin(userID uuid.UUID) error
	DeactivateUser(userID uuid.UUID) error
	ActivateUser(userID uuid.UUID) error
	DeleteUser(userID uuid.UUID) error
	FindActiveUsers() ([]*entities.User, error)
	FindAdmins() ([]*entities.User, error)
}