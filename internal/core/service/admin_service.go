package service

import (
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

func (s *AdminService) GetUsers(status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]*entities.User, error) {
	return s.db.FindUsers(status, role, sort, order)
}

func (s *AdminService) AdminGetUserByID(userID *uuid.UUID) (*entities.User, error) {
	return s.db.AdminGetUserByID(userID)
}

func (s *AdminService) AdminUpdateUser(userID *uuid.UUID, user *entities.User) error {
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
