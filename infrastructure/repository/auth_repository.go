package repository

import (
	"database/sql"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGAuthRepository struct {
	db *sql.DB
}

func NewPGAuthRepository(db *sql.DB) ports.AuthRepository {
	return &PGAuthRepository{db: db}
}

func (r *PGAuthRepository) Create(user *entities.User) error {
	// existingUser, err := r.FindByEmail(&user.Email)
	// if err != nil && err != sql.ErrNoRows {
	// 	return err
	// }

	// if existingUser != nil {
	// 	return errors.New("user already exists")
	// }

	query := `
		INSERT INTO users (id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.PhoneNumber,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Status,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *PGAuthRepository) FindUserByPhoneNumber(phoneNumber *string) (*entities.User, error) {
	query := `
	SELECT id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at, deleted_at
	FROM users
	WHERE phone_number = $1
`

	var user entities.User
	err := r.db.QueryRow(query, phoneNumber).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PGAuthRepository) FindUserByID(id uuid.UUID) (*entities.User, error) {
	query := `
	SELECT id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1
	`

	var user entities.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
