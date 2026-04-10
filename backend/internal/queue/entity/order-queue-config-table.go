package entity

import (
	"time"

	"github.com/google/uuid"
)

// OrdersQueueConfigTable хранит параметры расчета времени по очереди.
// Поля в минутах, кроме коэффициентов.
type OrdersQueueConfigTable struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	KitchenParallelSlots  int       `gorm:"not null;default:1" json:"kitchen_parallel_slots"`
	QueueGrowthFactor     float64   `gorm:"not null;default:0.15" json:"queue_growth_factor"`
	OrderReserveMinutes   int       `gorm:"not null;default:2" json:"order_reserve_minutes"`
	RestaurantOpenAt      string    `gorm:"not null;default:'10:00'" json:"restaurant_open_at"`
	RestaurantCloseAt     string    `gorm:"not null;default:'22:00'" json:"restaurant_close_at"`
	EmptyQueueWaitMinMins int       `gorm:"not null;default:10" json:"empty_queue_wait_min_mins"`
	EmptyQueueWaitMaxMins int       `gorm:"not null;default:10" json:"empty_queue_wait_max_mins"`
	QueueTimeMinFactor    float64   `gorm:"not null;default:0.90" json:"queue_time_min_factor"`
	QueueTimeMaxFactor    float64   `gorm:"not null;default:1.25" json:"queue_time_max_factor"`
	CreatedAt             time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt             time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}

func NewOrdersQueueConfigTable(kitchenParallelSlots int, queueGrowthFactor float64, orderReserveMinutes int, restaurantOpenAt string, restaurantCloseAt string, emptyQueueWaitMinMins int, emptyQueueWaitMaxMins int, queueTimeMinFactor float64, queueTimeMaxFactor float64) *OrdersQueueConfigTable {
	return &OrdersQueueConfigTable{
		ID:                    uuid.New(),
		KitchenParallelSlots:  kitchenParallelSlots,
		QueueGrowthFactor:     queueGrowthFactor,
		OrderReserveMinutes:   orderReserveMinutes,
		RestaurantOpenAt:      restaurantOpenAt,
		RestaurantCloseAt:     restaurantCloseAt,
		EmptyQueueWaitMinMins: emptyQueueWaitMinMins,
		EmptyQueueWaitMaxMins: emptyQueueWaitMaxMins,
		QueueTimeMinFactor:    queueTimeMinFactor,
		QueueTimeMaxFactor:    queueTimeMaxFactor,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}
