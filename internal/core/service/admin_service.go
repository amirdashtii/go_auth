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

func (s *AdminService) GetUsers(status, role, sort, order string) ([]*entities.User, error) {
	statusInt := entities.ParseStatusType(status)
	roleInt := entities.ParseRoleType(role)
	return s.db.FindUsers(statusInt, roleInt, sort, order)
}

func (s *AdminService) AdminGetUserByID(userID uuid.UUID) (*entities.User, error) {
	// TODO: Implement get user by ID logic
	return nil, nil
}

func (s *AdminService) AdminUpdateUser(userID uuid.UUID, user *entities.User) error {
	// TODO: Implement update user logic
	return s.db.AdminUpdateUser(user)
}

func (s *AdminService) PromoteToAdmin(userID uuid.UUID) error {
	// TODO: Implement promote to admin logic
	return nil
}

func (s *AdminService) DeactivateUser(userID uuid.UUID) error {
	// TODO: Implement deactivate user logic
	return nil
}

func (s *AdminService) ActivateUser(userID uuid.UUID) error {
	// TODO: Implement activate user logic
	return nil
}

func (s *AdminService) AdminDeleteUser(userID uuid.UUID) error {
	// TODO: Implement delete user logic
	return s.db.AdminDeleteUser(userID)
}

func (s *AdminService) FindActiveUsers() ([]*entities.User, error) {
	// TODO: Implement find active users logic
	return nil, nil
}

func (s *AdminService) FindAdmins() ([]*entities.User, error) {
	// TODO: Implement find admins logic
	return nil, nil
}
