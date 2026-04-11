package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrdersQueueOrder мок-модель заказа для расчета.
type OrdersQueueOrder struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ExternalID   string      `gorm:"not null;index:idx_queue_estimation_order_external_id" json:"external_id"`
	QueueID      uuid.UUID   `gorm:"-" json:"queue_id"`
	SetupMinutes int         `gorm:"not null;default:0" json:"setup_minutes"`
	Items        []OrderItem `gorm:"-" json:"items"`
	CreatedAt    time.Time   `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

// OrderItem мок-модель позиции заказа.
// GrowthFactor описывает нелинейное изменение времени для 2+, 3+, ... штук одного блюда.
type OrderItem struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID          uuid.UUID `gorm:"type:uuid;not null;index:idx_queue_estimation_item_order_id" json:"order_id"`
	ProductID        string    `gorm:"not null;index:idx_queue_estimation_item_product_id" json:"product_id"`
	Quantity         int       `gorm:"not null;default:1" json:"quantity"`
	BaseCookMinutes  int       `gorm:"not null;default:0" json:"base_cook_minutes"`
	GrowthFactor     float64   `gorm:"not null;default:0.0" json:"growth_factor"`
	ComplexityFactor float64   `gorm:"not null;default:1.0" json:"complexity_factor"`
	CreatedAt        time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

func NewOrderItem(orderID uuid.UUID, productID string, quantity int, baseCookMinutes int, growthFactor float64, complexityFactor float64) *OrderItem {
	return &OrderItem{
		ID:               uuid.New(),
		OrderID:          orderID,
		ProductID:        productID,
		Quantity:         quantity,
		BaseCookMinutes:  baseCookMinutes,
		GrowthFactor:     growthFactor,
		ComplexityFactor: complexityFactor,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
