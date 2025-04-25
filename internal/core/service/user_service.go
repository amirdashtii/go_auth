package service

import (
	"fmt"

	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if !user.DeletedAt.IsZero() {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, user *entities.User) error {
	user.ID = userID

	return s.db.Update(user)
}

func (s *UserService) ChangePassword(userID uuid.UUID, oldPassword, newPassword string) error {

	currentUser, err := s.db.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if !currentUser.DeletedAt.IsZero() {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect: %w", err)
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entities.User{}
	user.ID = userID
	user.Password = string(hashedNewPassword)
	return s.db.Update(user)
}

func (s *UserService) DeleteProfile(userID uuid.UUID) error {
	return s.db.Delete(userID)
}
