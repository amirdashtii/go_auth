package port

import "github.com/amirdashtii/go_auth/internal/core/entities"

type UserRepository interface {
	Create(user *entities.User) error
}
