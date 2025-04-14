package repository

import (
	"errors"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func (r *PGRepository) Create(user *entities.User) error {
	existingUser, err := r.FindByEmail(user.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return errors.New("user already exists")
	}

	query := `
		INSERT INTO users (id, email, password, is_active, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(query,
		user.ID,
		user.Email,
		user.Password,
		user.IsActive,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *PGRepository) FindByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, email, password, is_active, is_admin, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user entities.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
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

func (r *PGRepository) FindByID(id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT id, email, password, is_active, is_admin, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user entities.User	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
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
