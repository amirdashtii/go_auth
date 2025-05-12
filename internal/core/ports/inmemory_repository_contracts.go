package ports

import (
	"context"
	"time"
)

type InMemoryRespositoryContracts interface {
	AddToken(ctx context.Context, userID, token string, expiration time.Duration) error
	RemoveToken(ctx context.Context, userID string) error
	FindToken(ctx context.Context, userID string) (string, error)
}
