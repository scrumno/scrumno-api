package create_cart

import (
	"github.com/google/uuid"
)

type Command struct {
	UserID uuid.UUID
}
