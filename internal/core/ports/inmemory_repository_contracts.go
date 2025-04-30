package ports

import "time"

type InMemoryRespositoryContracts interface {
	AddToken(userID, token string, expiration time.Duration) error
	RemoveToken(userID string) error
	FindToken(userID string) (string, error)
}
