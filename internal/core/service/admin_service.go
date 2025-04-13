package service

import (
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type AdminService struct {
	db ports.UserRepository
}

func NewAdminService() *AdminService {
	db, err := repository.NewPGRepository()
	if err != nil {
		// Handle the error appropriately, e.g., log it or return it
		panic(err)
	}

	return &AdminService{
		db: db,
	}
}

func (s *AdminService) GetUsers() ([]*entities.User, error) {
	// TODO: Implement get users logic
	return nil, nil
}

func (s *AdminService) GetUserByID(userID uuid.UUID) (*entities.User, error) {
	// TODO: Implement get user by ID logic
	return nil, nil
}

func (s *AdminService) UpdateUser(userID uuid.UUID, user *entities.User) error {
	// TODO: Implement update user logic
	return nil
}

func (s *AdminService) PromoteToAdmin(userID uuid.UUID) error {
	// TODO: Implement promote to admin logic
	return nil
}

func (s *AdminService) DeactivateUser(userID uuid.UUID) error {
	// TODO: Implement deactivate user logic
	return nil
}

func (s *AdminService) ActivateUser(userID uuid.UUID) error {
	// TODO: Implement activate user logic
	return nil
}

func (s *AdminService) DeleteUser(userID uuid.UUID) error {
	// TODO: Implement delete user logic
	return nil
}

func (s *AdminService) FindActiveUsers() ([]*entities.User, error) {
	// TODO: Implement find active users logic
	return nil, nil
}

func (s *AdminService) FindAdmins() ([]*entities.User, error) {
	// TODO: Implement find admins logic
	return nil, nil
}