package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	IsActive   bool      `json:"is_active"`
	IsAdmin    bool      `json:"is_admin"`
	CreatedAt  time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
