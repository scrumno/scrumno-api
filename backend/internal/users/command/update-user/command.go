package update_user

import (
	"time"

	"github.com/google/uuid"
)

type Command struct {
	ID        string
	Phone     *string `json:"phone"`
	FullName  *string `json:"full_name"`
	BirthDate *string `json:"birth_date"`
}

type UserDTO struct {
	ID        uuid.UUID  `json:"id"`
	Phone     string     `json:"phone"`
	FullName  *string    `json:"full_name,omitempty"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
}
