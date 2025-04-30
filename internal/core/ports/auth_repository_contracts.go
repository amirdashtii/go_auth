package ports

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Create(user *entities.User) error
	FindUserByPhoneNumber(phoneNumber *string) (*entities.User, error)
	FindUserByID(id uuid.UUID) (*entities.User, error)
}
