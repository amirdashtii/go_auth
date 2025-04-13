package repository

import (
	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (r *PGRepository) Create(user *entities.User) error {
	// TODO: Implement create user logic
	return nil
}

func (r *PGRepository) FindByEmail(email string) (*entities.User, error) {
	// TODO: Implement find user by email logic
	return nil, nil
}

func (r *PGRepository) FindByID(id uuid.UUID) (*entities.User, error) {
	// TODO: Implement find user by ID logic
	return nil, nil
}

func (r *PGRepository) Update(user *entities.User) error {
	// TODO: Implement update user logic
	return nil
}

func (r *PGRepository) Delete(id uuid.UUID) error {
	// TODO: Implement delete user logic
	return nil
}

func (r *PGRepository) FindAll() ([]*entities.User, error) {
	// TODO: Implement find all users logic
	return nil, nil
}

func (r *PGRepository) FindActiveUsers() ([]*entities.User, error) {
	// TODO: Implement find active users logic
	return nil, nil
}

func (r *PGRepository) FindAdmins() ([]*entities.User, error) {
	// TODO: Implement find admins logic
	return nil, nil
}
