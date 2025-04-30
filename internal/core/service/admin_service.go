package service

import (
	"github.com/amirdashtii/go_auth/controller/dto"
	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type AdminService struct {
	db ports.AdminRepository
}

func NewAdminService() *AdminService {
	dbRepo, err := repository.NewPGRepository()
	if err != nil {
		panic(err)
	}
	db := dbRepo.DB()
	adminRepo := repository.NewPGAdminRepository(db)
	return &AdminService{
		db: adminRepo,
	}
}

func (s *AdminService) GetUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]dto.AdminUserResponse, error) {
	users, err := s.db.FindUsers(status, role, sort, order)
	if err != nil {
		return nil, err
	}

	var resp []dto.AdminUserResponse
	for _, u := range users {
		resp = append(resp, dto.AdminUserResponse{
			ID:        u.ID.String(),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Status:    u.Status.String(),
			Role:      u.Role.String(),
		})
	}

	return resp, err
}

func (s *AdminService) AdminGetUserByID(userID *uuid.UUID) (*dto.AdminUserResponse, error) {

	user, err := s.db.AdminGetUserByID(userID)
	if err != nil {
		return nil, err
	}
	resp := &dto.AdminUserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Status:    user.Status.String(),
		Role:      user.Role.String(),
	}
	return resp, nil
}

func (s *AdminService) AdminUpdateUser(userID *uuid.UUID, req *dto.AdminUserUpdateRequest) error {
	user := &entities.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	user.ID = *userID
	return s.db.AdminUpdateUser(user)
}

func (s *AdminService) ChangeUserRole(userID *uuid.UUID, role *entities.RoleType) error {

	return s.db.AdminChangeUserRole(userID, role)
}

func (s *AdminService) ChangeUserStatus(userID *uuid.UUID, status *entities.StatusType) error {

	return s.db.AdminChangeUserStatus(userID, status)
}

func (s *AdminService) AdminDeleteUser(userID *uuid.UUID) error {
	return s.db.AdminDeleteUser(userID)
}
