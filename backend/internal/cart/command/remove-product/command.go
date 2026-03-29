package remmove_product

import (
	"github.com/google/uuid"
)

type Command struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
}
