package ports

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindUserByID(id *uuid.UUID) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id *uuid.UUID) error
}
