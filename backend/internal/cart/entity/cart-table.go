package entity

import (
	"time"
	
	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_carts_user" json:"user_id"`
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
	TotalAmount float64 `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CartID    uuid.UUID `gorm:"type:uuid;not null;index:idx_cart_items_cart" json:"cart_id"`
	ProductID string `gorm:"type:text;not null;index:idx_cart_items_product" json:"product_id"`
	Quantity  int `gorm:"type:int;not null" json:"quantity"`
	UnitPrice float64 `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func NewCart(userID uuid.UUID) *Cart {
	return &Cart{
		ID:        uuid.New(),
		UserID:    userID,
		TotalAmount: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Items:     []CartItem{},
	}
}
