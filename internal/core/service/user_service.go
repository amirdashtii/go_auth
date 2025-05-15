package service

import (
	"context"
	"os"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db     ports.UserRepository
	logger ports.Logger
}

func NewUserService() *UserService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(errors.ErrDatabaseInit)
	}
	db := dbRepo.DB()

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      os.Stdout,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	userRepo := repository.NewPGUserRepository(db, appLogger)
	return &UserService{
		db:     userRepo,
		logger: appLogger,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID *uuid.UUID) (*dto.UserProfileResponse, error) {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while getting user profile",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return nil, errors.ErrContextCancelled
	}
	user, err := s.db.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user.Status == entities.Deleted {
		s.logger.Error("User is deleted",
			ports.F("user_id", userID),
		)
		return nil, errors.ErrInvalidCredentials
	}
	if user.Status == entities.Deactivated {
		s.logger.Error("User is deactivated",
			ports.F("user_id", userID),
		)
		return nil, errors.ErrAccountDeactivated
	}

	resp := dto.UserProfileResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return &resp, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID *uuid.UUID, req *dto.UserUpdateRequest) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while updating user profile",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}
	user := &entities.User{
		ID:          *userID,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
	}

	if err := s.db.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while changing user password",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}
	currentUser, err := s.db.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if currentUser.Status == entities.Deleted {
		s.logger.Error("User is deleted",
			ports.F("user_id", userID),
		)
		return errors.ErrInvalidCredentials
	}
	if currentUser.Status == entities.Deactivated {
		s.logger.Error("User is deactivated",
			ports.F("user_id", userID),
		)
		return errors.ErrAccountDeactivated
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(changePasswordReq.OldPassword)); err != nil {
		s.logger.Error("Old password is incorrect",
			ports.F("user_id", userID),
		)
		return errors.ErrInvalidCredentials
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error generating new password hash",
			ports.F("error", err),
			ports.F("user_id", userID),
		)
		return errors.ErrChangePassword
	}

	user := &entities.User{}
	user.ID = *userID
	user.Password = string(hashedNewPassword)
	if err := s.db.Update(ctx, user); err != nil {
		s.logger.Error("Error updating user password",
			ports.F("error", err),
			ports.F("user_id", userID),
		)
		return errors.ErrChangePassword
	}
	return nil
}

func (s *UserService) DeleteProfile(ctx context.Context, userID *uuid.UUID) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while deleting user profile",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}

	if err := s.db.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}
