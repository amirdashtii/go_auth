package repository

import (
	"context"
	"database/sql"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGAuthRepository struct {
	db     *sql.DB
	logger ports.Logger
}

func NewPGAuthRepository(db *sql.DB, logger ports.Logger) ports.AuthRepository {
	return &PGAuthRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PGAuthRepository) Create(ctx context.Context, user *entities.User) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while creating user",
			ports.F("error", ctx.Err()),
			ports.F("phone_number", user.PhoneNumber),
		)
		return errors.ErrContextCancelled
	}

	query := `
		INSERT INTO users (id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
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
		r.logger.Error("Database error in Create",
			ports.F("error", err),
			ports.F("phone_number", user.PhoneNumber),
		)

		if err.Error() == "duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return errors.ErrDuplicatePhoneNumber
		}
		return errors.ErrCreateUser
	}

	return nil
}

func (r *PGAuthRepository) FindUserByPhoneNumber(ctx context.Context, phoneNumber *string) (*entities.User, error) {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while finding user by phone number",
			ports.F("error", ctx.Err()),
			ports.F("phone_number", phoneNumber),
		)
		return nil, errors.ErrContextCancelled
	}

	query := `
	SELECT id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at, deleted_at
	FROM users
	WHERE phone_number = $1
`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, phoneNumber).Scan(
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
		r.logger.Error("Database error in FindUserByPhoneNumber",
			ports.F("error", err),
			ports.F("phone_number", phoneNumber),
		)

		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrGetUser
	}

	return &user, nil
}

func (r *PGAuthRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while finding user by ID",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
		)
		return nil, errors.ErrContextCancelled
	}

	query := `
	SELECT id, phone_number, password, first_name, last_name, email, status, role, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		r.logger.Error("Database error in FindUserByID",
			ports.F("error", err),
			ports.F("user_id", id),
		)

		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrGetUser
	}

	return &user, nil
}
