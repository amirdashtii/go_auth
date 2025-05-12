package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGUserRepository struct {
	db     *sql.DB
	logger ports.Logger
}

func NewPGUserRepository(db *sql.DB, logger ports.Logger) ports.UserRepository {
	return &PGUserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PGUserRepository) FindUserByID(ctx context.Context, id *uuid.UUID) (*entities.User, error) {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while finding user by ID",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
		)
		return nil, errors.ErrContextCancelled
	}
	query := `
	SELECT id, phone_number, first_name, last_name, email, password, status, role, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1
	`

	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.PhoneNumber,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		r.logger.Error("Database error in FindUserByID",
			ports.F("error", err),
			ports.F("id", id),
		)
		if err == sql.ErrNoRows {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrGetUser
	}

	return &user, nil
}

func (r *PGUserRepository) Update(ctx context.Context, user *entities.User) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while updating user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", user.ID),
		)
		return errors.ErrContextCancelled
	}
	query := "UPDATE users SET "
	args := []interface{}{}
	i := 1

	if user.PhoneNumber != "" {
		query += "phone_number = $" + fmt.Sprint(i) + ", "
		args = append(args, user.PhoneNumber)
		i++
	}
	if user.FirstName != "" {
		query += "first_name = $" + fmt.Sprint(i) + ", "
		args = append(args, user.FirstName)
		i++
	}
	if user.LastName != "" {
		query += "last_name = $" + fmt.Sprint(i) + ", "
		args = append(args, user.LastName)
		i++
	}
	if user.Email != "" {
		query += "email = $" + fmt.Sprint(i) + ", "
		args = append(args, user.Email)
		i++
	}
	if user.Password != "" {
		query += "password = $" + fmt.Sprint(i) + ", "
		args = append(args, user.Password)
		i++
	}

	query += "updated_at = NOW()"

	query += " WHERE id = $" + fmt.Sprint(i)
	args = append(args, user.ID)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Database error in Update",
			ports.F("error", err),
			ports.F("user", user),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		if err.Error() == "duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return errors.ErrDuplicatePhoneNumber
		}
		if err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.ErrDuplicateEmail
		}
		return errors.ErrUpdateUser
	}
	return nil
}

func (r *PGUserRepository) Delete(ctx context.Context, id *uuid.UUID) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while deleting user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
		)
		return errors.ErrContextCancelled
	}
	query := `UPDATE users SET deleted_at = NOW(), status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, entities.Deleted)
	if err != nil {
		r.logger.Error("Database error in Delete",
			ports.F("error", err),
			ports.F("id", id),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrDeleteUser
	}
	return nil
}

func (r *PGUserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while creating user",
			ports.F("error", ctx.Err()),
			ports.F("user", user),
		)
		return errors.ErrContextCancelled
	}
	query := `
		INSERT INTO users (id, phone_number, first_name, last_name, email, password, status, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.PhoneNumber,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Status,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Database error in CreateUser",
			ports.F("error", err),
			ports.F("user", user),
		)
		if err.Error() == "duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return errors.ErrDuplicatePhoneNumber
		}
		if err.Error() == "duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.ErrDuplicateEmail
		}
		return errors.ErrCreateUser
	}
	return nil
}

func (r *PGUserRepository) UpdatePassword(ctx context.Context, userID *uuid.UUID, hashedPassword string) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while updating user password",
			ports.F("error", ctx.Err()),
			ports.F("user_id", userID),
		)
		return errors.ErrContextCancelled
	}
	query := "UPDATE users SET password = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	if err != nil {
		r.logger.Error("Database error in UpdatePassword",
			ports.F("error", err),
			ports.F("user_id", userID),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrUpdateUser
	}
	return nil
}
