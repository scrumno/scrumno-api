package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderDraftTable struct {
	ID                    uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID                uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	VenueID               uuid.UUID      `gorm:"type:uuid;not null;index" json:"venue_id"`
	CartSnapshotJSON      datatypes.JSON `gorm:"type:jsonb;not null" json:"cart_snapshot_json"`
	Amount                float64        `gorm:"type:numeric(12,2);not null" json:"amount"`
	PaymentStatus         bool           `gorm:"not null;default:false" json:"payment_status"`
	ProviderCreateStatus  bool           `gorm:"not null;default:false" json:"provider_create_status"`
	ProviderPending       bool           `gorm:"not null;default:false" json:"provider_pending"`
	ProviderCorrelationID *uuid.UUID     `gorm:"type:uuid;index" json:"provider_correlation_id,omitempty"`
	ProviderOrderID       *uuid.UUID     `gorm:"type:uuid;index" json:"provider_order_id,omitempty"`
	ProviderError         *string        `gorm:"type:text" json:"provider_error,omitempty"`
	CreatedAt             time.Time      `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

func (OrderDraftTable) TableName() string {
	return "order-draft-table"
}
