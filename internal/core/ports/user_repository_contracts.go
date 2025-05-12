package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type UserRepository interface {
	FindUserByID(ctx context.Context, id *uuid.UUID) (*entities.User, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id *uuid.UUID) error
}
