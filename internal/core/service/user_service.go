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
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(err)
	}
	db := dbRepo.DB()
	userRepo := repository.NewPGUserRepository(db)
	return &UserService{
		db: userRepo,
	}
}

func (s *UserService) GetOwnProfile(userID string) (*entities.User, error) {
	uuid := uuid.Must(uuid.Parse(userID))

	user, err := s.db.FindUserOwnByID(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, user *entities.User) error {
	// TODO: Implement update profile logic
	return nil
}

func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	// TODO: Implement change password logic
	return nil
}
