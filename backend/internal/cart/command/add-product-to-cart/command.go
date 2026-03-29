package add_product_to_cart

import (
	"github.com/google/uuid"
)

type Command struct {
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  float64
	BasePrice float64
}
