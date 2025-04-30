package entities

import (
	"strings"
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

func ParseRoleType(s string) RoleType {
	switch strings.ToLower(s) {
	case "superadmin":
		return SuperAdminRole
	case "admin":
		return AdminRole
	case "user":
		return UserRole
	default:
		return UserRole
	}
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

func ParseStatusType(s string) StatusType {
	switch strings.ToLower(s) {
	case "active":
		return Active
	case "deactivated":
		return Deactivated
	case "deleted":
		return Deleted
	default:
		return Active
	}
}

type User struct {
	ID          uuid.UUID  `json:"id"`
	PhoneNumber string     `json:"phone_number"`
	Password    string     `json:"password"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Status      StatusType `json:"status"`
	Role        RoleType   `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
