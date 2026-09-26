package handler

import (
	"time"

	"github.com/google/uuid"
)

type userResponse struct {
	ID        uuid.UUID  `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
