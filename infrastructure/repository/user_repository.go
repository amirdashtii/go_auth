package repository

import (
	"database/sql"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGUserRepository struct {
	db *sql.DB
}

func NewPGUserRepository(db *sql.DB) ports.UserRepository {
	return &PGUserRepository{db: db}
}

func (r *PGUserRepository) FindUserOwnByID(id uuid.UUID) (*entities.User, error) {
	query := `
	SELECT id, first_name, last_name, email, password, is_active, is_admin, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	var user entities.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PGUserRepository) Update(user *entities.User) error {
	return nil
}

func (r *PGUserRepository) Delete(id uuid.UUID) error {
	return nil
}
