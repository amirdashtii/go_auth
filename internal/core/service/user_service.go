package service

import (
	"fmt"

	"github.com/amirdashtii/go_auth/controller/dto"
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

func (s *UserService) GetProfile(userID *uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.db.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user.Status == entities.Deleted {
		return nil, fmt.Errorf("user not found")
	}
	if user.Status == entities.Deactivated {
		return nil, fmt.Errorf("user is deactivated")
	}

	resp := dto.UserProfileResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return &resp, nil
}

func (s *UserService) UpdateProfile(userID *uuid.UUID, req *dto.UserUpdateRequest) error {

	user := &entities.User{
		ID:        *userID,
		PhoneNumber: req.PhoneNumber,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	return s.db.Update(user)
}

func (s *UserService) ChangePassword(userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error {

	currentUser, err := s.db.FindUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if currentUser.Status == entities.Deleted {
		return fmt.Errorf("user not found")
	}
	if currentUser.Status == entities.Deactivated {
		return fmt.Errorf("user is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(changePasswordReq.OldPassword)); err != nil {
		return fmt.Errorf("current password is incorrect: %w", err)
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entities.User{}
	user.ID = *userID
	user.Password = string(hashedNewPassword)
	return s.db.Update(user)
}

func (s *UserService) DeleteProfile(userID *uuid.UUID) error {
	return s.db.Delete(userID)
}
