package ports

import (
	"context"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindUserByPhoneNumber(ctx context.Context, phoneNumber *string) (*entities.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}
