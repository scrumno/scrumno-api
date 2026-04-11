package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderHistoryTable struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DraftID         *uuid.UUID `gorm:"type:uuid;index" json:"draft_id,omitempty"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	VenueID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"venue_id"`
	ProviderOrderID uuid.UUID  `gorm:"type:uuid;not null;index" json:"provider_order_id"`
	Status          string     `gorm:"type:varchar(64);not null;default:'Created'" json:"status"`
	CreatedAt       time.Time  `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

func (OrderHistoryTable) TableName() string {
	return "order-history-table"
}
