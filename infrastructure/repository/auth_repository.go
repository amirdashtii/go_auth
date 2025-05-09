package repository

import (
	"database/sql"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrDuplicatePhoneNumber = errors.New(errors.ValidationError, "Phone number already exists", "این شماره تلفن قبلاً ثبت شده است", nil)
	ErrUserNotFound         = errors.New(errors.NotFoundError, "Invalid credentials", "نام کاربری یا رمز عبور اشتباه است", nil)
)

type PGAuthRepository struct {
	db *sql.DB
}

func NewPGAuthRepository(db *sql.DB) ports.AuthRepository {
	return &PGAuthRepository{db: db}
}

func (r *PGAuthRepository) Create(user *entities.User) error {
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

	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return errors.ErrDuplicatePhoneNumber
		}
		return errors.ErrCreateUser
	}

	return nil
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
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrGetUser
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
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
