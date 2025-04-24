package repository

import (
	"database/sql"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGAdminRepository struct {
	db *sql.DB
}

func NewPGAdminRepository(db *sql.DB) ports.AdminRepository {
	return &PGAdminRepository{db: db}
}

func (r *PGAdminRepository) FindAll() ([]*entities.User, error) {
	return nil, nil
}

func (r *PGAdminRepository) FindActiveUsers() ([]*entities.User, error) {
	return nil, nil
}

func (r *PGAdminRepository) FindAdmins() ([]*entities.User, error) {
	return nil, nil
}

func (r *PGAdminRepository) AdminUpdateUser(user *entities.User) error {
	return nil
}

func (r *PGAdminRepository) AdminDeleteUser(id uuid.UUID) error {
	return nil
}
