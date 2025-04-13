package service

import (
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type UserService struct {
	db ports.UserRepository
}

func NewUserService() *UserService {
	db, err := repository.NewPGRepository()
	if err != nil {
		// Handle the error appropriately, e.g., log it or return it
		panic(err)
	}

	return &UserService{
		db: db,
	}
}

func (s *UserService) GetProfile(userID uuid.UUID) (*entities.User, error) {
	// TODO: Implement get profile logic
	return nil, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, user *entities.User) error {
	// TODO: Implement update profile logic
	return nil
}

func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	// TODO: Implement change password logic
	return nil
}

 