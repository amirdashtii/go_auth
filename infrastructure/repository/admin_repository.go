package repository

import (
	"database/sql"
	"fmt"

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

func (r *PGAdminRepository) FindUsers(status entities.StatusType, role entities.RoleType, sort, order string) ([]*entities.User, error) {

	query := fmt.Sprintf(`
	SELECT * FROM users
	WHERE status = $1 AND role = $2
	ORDER BY %s %s
	`, sort, order)
	rows, err := r.db.Query(query, status, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*entities.User{}
	for rows.Next() {
		fmt.Println(rows)
		var user entities.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Status, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}
		fmt.Println(user)
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
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
