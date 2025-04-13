package service

import (
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
)

type AuthService struct {
	db ports.UserRepository
}

func NewAuthService() *AuthService {
	db, err := repository.NewPGRepository()
	if err != nil {
		// Handle the error appropriately, e.g., log it or return it
		panic(err)
	}

	return &AuthService{
		db: db,
	}
}

func (s *AuthService) Register(user *entities.User) error {
	// TODO: Implement register logic
	return nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	// TODO: Implement login logic
	return "", nil
}

func (s *AuthService) Logout(token string) error {
	// TODO: Implement logout logic
	return nil
}

func (s *AuthService) RefreshToken(token string) (string, error) {
	// TODO: Implement refresh token logic
	return "", nil
}

func (s *AuthService) ValidateToken(token string) (*entities.User, error) {
	// TODO: Implement validate token logic
	return nil, nil
}
