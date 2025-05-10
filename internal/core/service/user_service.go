package service

import (
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
	db ports.UserRepository
	logger ports.Logger
}

func NewUserService() *UserService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(errors.ErrDatabaseInit)
	}
	db := dbRepo.DB()

	// Create log file
	logfile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Initialize logger with both file and console output
	loggerConfig := ports.LoggerConfig{
		Level: "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output: logfile,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	userRepo := repository.NewPGUserRepository(db, appLogger)
	return &UserService{
		db: userRepo,
		logger: appLogger,
	}
}

func (s *UserService) GetProfile(userID *uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.db.FindUserByID(userID)
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

func (s *UserService) UpdateProfile(userID *uuid.UUID, req *dto.UserUpdateRequest) error {
	user := &entities.User{
		ID:          *userID,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
	}

	if err := s.db.Update(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) ChangePassword(userID *uuid.UUID, changePasswordReq *dto.ChangePasswordRequest) error {
	currentUser, err := s.db.FindUserByID(userID)
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
	if err := s.db.Update(user); err != nil {
		s.logger.Error("Error updating user password",
			ports.F("error", err),
			ports.F("user_id", userID),
		)
		return errors.ErrChangePassword
	}
	return nil
}

func (s *UserService) DeleteProfile(userID *uuid.UUID) error {
	if err := s.db.Delete(userID); err != nil {
		return err
	}
	return nil
}
