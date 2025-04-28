package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RoleType int

const (
	UserRole RoleType = iota
	SuperAdminRole
	AdminRole
)

func (r RoleType) String() string {
	switch r {
	case SuperAdminRole:
		return "SuperAdmin"
	case AdminRole:
		return "Admin"
	case UserRole:
		return "User"
	default:
		return "Unknown"
	}
}

func (r RoleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RoleType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "SuperAdmin":
		*r = SuperAdminRole
	case "Admin":
		*r = AdminRole
	case "User":
		*r = UserRole
	default:
		*r = UserRole
	}
	return nil
}

type StatusType int

const (
	Active StatusType = iota
	Deactivated
	Deleted
)

func (s StatusType) String() string {
	switch s {
	case Active:
		return "Active"
	case Deactivated:
		return "Deactivated"
	case Deleted:
		return "Deleted"
	default:
		return "Unknown"
	}
}

func (s StatusType) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *StatusType) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "Active":
		*s = Active
	case "Deactivated":
		*s = Deactivated
	case "Deleted":
		*s = Deleted
	default:
		*s = Active
	}
	return nil
}

type User struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Status    StatusType `json:"status"`
	Role      RoleType   `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt time.Time  `json:"deleted_at"`
}
