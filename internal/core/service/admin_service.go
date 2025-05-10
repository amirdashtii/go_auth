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
)

type AdminService struct {
	db ports.AdminRepository
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
		Level: "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output: logfile,
	}
	appLogger := logger.NewZerologLogger(loggerConfig)
	
	adminRepo := repository.NewPGAdminRepository(db, appLogger)
	return &AdminService{
		db: adminRepo,
		logger: appLogger,
	}
}

func (s *AdminService) GetUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error) {
	users, err := s.db.FindUsers(status, role, sort, order)
	if err != nil {
		s.logger.Error("Error getting users",
			ports.F("error", err),
		)
		return nil, errors.ErrGetUsers
	}

	var resp []dto.AdminUserResponse
	for _, u := range users {
		resp = append(resp, dto.AdminUserResponse{
			ID:          u.ID.String(),
			PhoneNumber: u.PhoneNumber,
			FirstName:   u.FirstName,
			LastName:    u.LastName,
			Email:       u.Email,
			Status:      u.Status.String(),
			Role:        u.Role.String(),
		})
	}

	return resp, nil
}

func (s *AdminService) AdminGetUserByID(userID *uuid.UUID) (*dto.AdminUserResponse, error) {
	user, err := s.db.AdminGetUserByID(userID)
	if err != nil {
		return nil, err
	}
	resp := &dto.AdminUserResponse{
		ID:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Status:      user.Status.String(),
		Role:        user.Role.String(),
	}
	return resp, nil
}

func (s *AdminService) AdminUpdateUser(userID *uuid.UUID, req *dto.AdminUserUpdateRequest) error {
	user := &entities.User{
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
	}
	user.ID = *userID
	if err := s.db.AdminUpdateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *AdminService) ChangeUserRole(userID *uuid.UUID, role *entities.RoleType) error {
	if err := s.db.AdminChangeUserRole(userID, role); err != nil {
		return err
	}
	return nil
}

func (s *AdminService) ChangeUserStatus(userID *uuid.UUID, status *entities.StatusType) error {
	if err := s.db.AdminChangeUserStatus(userID, status); err != nil {
		return err
	}
	return nil
}

func (s *AdminService) AdminDeleteUser(userID *uuid.UUID) error {
	if err := s.db.AdminDeleteUser(userID); err != nil {
		return err
	}
	return nil
}
