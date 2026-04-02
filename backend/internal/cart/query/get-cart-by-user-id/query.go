package get_cart_by_user_id

import (
	"github.com/google/uuid"
)

type Query struct {
	UserID uuid.UUID
}