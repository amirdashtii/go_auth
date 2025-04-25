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

func (s *UserService) GetProfile(userID uuid.UUID) (*entities.User, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, user *entities.User) error {
	user.ID = userID

	return s.db.Update(user)
}

func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {
	// TODO: Implement change password logic
	return nil
}
