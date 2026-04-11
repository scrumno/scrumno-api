package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderSubscribersTable struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID      uuid.UUID `gorm:"type:uuid;not null;index" json:"order_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ConnectionID string    `gorm:"type:varchar(128);not null;index" json:"connection_id"`
	IsActive     bool      `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt    time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

func (OrderSubscribersTable) TableName() string {
	return "order-subscribers-table"
}
