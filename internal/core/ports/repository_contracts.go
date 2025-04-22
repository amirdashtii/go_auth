package ports

import (
	"time"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id uuid.UUID) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
	FindAll() ([]*entities.User, error)
	FindActiveUsers() ([]*entities.User, error)
	FindAdmins() ([]*entities.User, error)
}

type InMemoryRespositoryContracts interface {
	AddToken(userID, token string, expiration time.Duration) error
	RemoveToken(userID string) error
	FindToken(userID string) (string, error)
}
