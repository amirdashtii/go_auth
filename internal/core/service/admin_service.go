package service

import (
	"context"
	"os"
	"time"

	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type AdminService struct {
	db     ports.AdminRepository
	logger ports.Logger
}

func NewAdminService() *AdminService {
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
		Level:       "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output:      logfile,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)

	adminRepo := repository.NewPGAdminRepository(db, appLogger)
	return &AdminService{
		db:     adminRepo,
		logger: appLogger,
	}
}

func (s *AdminService) GetUsers(ctx context.Context, status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error) {	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while getting users",
			ports.F("error", ctx.Err()),
			ports.F("status", status),
			ports.F("role", role),
		)
		return nil, errors.ErrContextCancelled
	}
	users, err := s.db.FindUsers(ctx, status, role, sort, order)
	if err != nil {
		s.logger.Error("Error getting users",
			ports.F("error", err),
			ports.F("status", status),
			ports.F("role", role),
		)
		return nil, errors.ErrGetUsers
	}

	var response []dto.AdminUserResponse
	for _, user := range users {
		response = append(response, dto.AdminUserResponse{
			ID:          user.ID.String(),
			PhoneNumber: user.PhoneNumber,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Status:      user.Status.String(),
			Role:        user.Role.String(),
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}
	return response, nil
}

func (s *AdminService) AdminGetUserByID(ctx context.Context, userID *uuid.UUID) (*dto.AdminUserResponse, error) {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while getting user by ID",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return nil, errors.ErrContextCancelled
	}
	user, err := s.db.AdminGetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.AdminUserResponse{
		ID:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Status:      user.Status.String(),
		Role:        user.Role.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *AdminService) AdminUpdateUser(ctx context.Context, userID *uuid.UUID, updateReq *dto.AdminUserUpdateRequest) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while updating user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}
	user := &entities.User{
		ID:          *userID,
		PhoneNumber: updateReq.PhoneNumber,
		FirstName:   updateReq.FirstName,
		LastName:    updateReq.LastName,
		Email:       updateReq.Email,
		UpdatedAt:   time.Now(),
	}
	if err := s.db.AdminUpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) ChangeUserRole(ctx context.Context, userID *uuid.UUID, updateRole *entities.RoleType) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while changing user role",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
			ports.F("new_role", updateRole),
		)
		return errors.ErrContextCancelled
	}
	if err := s.db.AdminChangeUserRole(ctx, userID, updateRole); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) ChangeUserStatus(ctx context.Context, userID *uuid.UUID, updateStatus *entities.StatusType) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while changing user status",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
			ports.F("new_status", updateStatus),
		)
		return errors.ErrContextCancelled
	}
	if err := s.db.AdminChangeUserStatus(ctx, userID, updateStatus); err != nil {
		return err
	}

	return nil
}

func (s *AdminService) AdminDeleteUser(ctx context.Context, userID *uuid.UUID) error {
	if ctx.Err() != nil {
		s.logger.Error("Context cancelled while deleting user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}
	if err := s.db.AdminDeleteUser(ctx, userID); err != nil {
		return err
	}

	return nil
}
