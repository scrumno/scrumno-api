package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrdersQueueTable struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index:idx_order_queue_order_id" json:"order_id"`
	QueueID   uuid.UUID `gorm:"type:uuid;not null;index:idx_order_queue_queue_id" json:"queue_id"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
