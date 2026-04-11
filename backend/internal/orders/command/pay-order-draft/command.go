package pay_order_draft

import "github.com/google/uuid"

type Command struct {
	DraftID uuid.UUID
}
