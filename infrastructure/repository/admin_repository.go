package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/amirdashtii/go_auth/internal/core/entities"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/google/uuid"
)

type PGAdminRepository struct {
	db     *sql.DB
	logger ports.Logger
}

func NewPGAdminRepository(db *sql.DB, logger ports.Logger) ports.AdminRepository {
	return &PGAdminRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PGAdminRepository) FindUsers(ctx context.Context, status *entities.StatusType, role *entities.RoleType, sort, order *string) ([]entities.User, error) {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while finding users",
			ports.F("error", ctx.Err()),
			ports.F("status", status),
			ports.F("role", role),
		)
		return nil, errors.ErrContextCancelled
	}
	query := fmt.Sprintf(`
	SELECT * FROM users
	WHERE status = $1 AND role = $2
	ORDER BY %s %s
	`, *sort, *order)

	rows, err := r.db.QueryContext(ctx, query, status, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(
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
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PGAdminRepository) AdminGetUserByID(ctx context.Context, id *uuid.UUID) (*entities.User, error) {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while getting user by ID",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
		)
		return nil, errors.ErrContextCancelled
	}
	query := `SELECT * FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var user entities.User
	err := row.Scan(
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
		r.logger.Error("Database error in AdminGetUserByID",
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

func (r *PGAdminRepository) AdminUpdateUser(ctx context.Context, user *entities.User) error {
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

	query += "updated_at = NOW()"

	query += " WHERE id = $" + fmt.Sprint(i)
	args = append(args, user.ID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		r.logger.Error("Database error in AdminUpdateUser",
			ports.F("error", err),
			ports.F("user_id", user.ID),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}

		if err.Error() == "duplicate key value violates unique constraint \"users_phone_number_key\"" {
			return errors.ErrDuplicatePhoneNumber
		}

		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return errors.ErrDuplicateEmail
		}
		return errors.ErrUpdateUser
	}
	return nil
}

func (r *PGAdminRepository) AdminChangeUserRole(ctx context.Context, id *uuid.UUID, role *entities.RoleType) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while changing user role",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
			ports.F("new_role", role),
		)
		return errors.ErrContextCancelled
	}
	query := `UPDATE users SET role = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, role, id)
	if err != nil {
		r.logger.Error("Database error in AdminChangeUserRole",
			ports.F("error", err),
			ports.F("user_id", id),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrChangeRole
	}
	return nil
}

func (r *PGAdminRepository) AdminChangeUserStatus(ctx context.Context, id *uuid.UUID, status *entities.StatusType) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while changing user status",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
			ports.F("new_status", status),
		)
		return errors.ErrContextCancelled
	}
	query := `UPDATE users SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		r.logger.Error("Database error in AdminChangeUserStatus",
			ports.F("error", err),
			ports.F("user_id", id),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrChangeStatus
	}
	return nil
}

func (r *PGAdminRepository) AdminDeleteUser(ctx context.Context, id *uuid.UUID) error {
	if ctx.Err() != nil {
		r.logger.Error("Context cancelled while deleting user",
			ports.F("error", ctx.Err()),
			ports.F("user_id", id),
		)
		return errors.ErrContextCancelled
	}
	query := `UPDATE users SET deleted_at = $1, status = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, time.Now(), entities.Deleted, id)
	if err != nil {
		r.logger.Error("Database error in AdminDeleteUser",
			ports.F("error", err),
			ports.F("user_id", id),
		)
		if err == sql.ErrNoRows {
			return errors.ErrUserNotFound
		}
		return errors.ErrDeleteUser
	}
	return nil
}
