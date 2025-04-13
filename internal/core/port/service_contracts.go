package port

import "github.com/amirdashtii/go_auth/internal/core/entities"

type AuthService interface {
	Register(user *entities.User) error
}