package create_order_draft

import (
	"github.com/google/uuid"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
)

type Command struct {
	UserID        uuid.UUID
	VenueID       uuid.UUID
	CustomerPhone string
	CustomerName  string
	Comment       *string
	SourceKey     string
	CartItems     []createOrder.CartLineItem
}
