package entity

import (
	"time"
	
	"github.com/google/uuid"
)

type ProductType string

const (
	ProductTypeProduct  ProductType = "Product"
	ProductTypeCompound ProductType = "Compound"
)

type Cart struct {
	ID        uuid.UUID 	`gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID 	`gorm:"type:uuid;not null;index:idx_carts_user" json:"user_id"`
	Items     []CartItem 	`gorm:"foreignKey:CartID" json:"items"`
	TotalAmount float64 	`gorm:"type:numeric(12,2);not null" json:"total_amount"`
	CreatedAt time.Time 	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time 	`gorm:"autoUpdateTime" json:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID 	`gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CartID    uuid.UUID 	`gorm:"type:uuid;not null;index:idx_cart_items_cart" json:"cart_id"`
	ProductID string 		`gorm:"type:text;not null;index:idx_cart_items_product" json:"product_id"`
	Type	  ProductType 	`gorm:"type:varchar(20);not null;default:'Product'" json:"type"`
	Quantity  int 			`gorm:"type:int;not null;check:quantity >= 0" json:"quantity;"`
	UnitPrice float64 		`gorm:"type:numeric(12,2);not null;check:unit_price > 0" json:"unit_price"`
	Modifiers JSONB 		`gorm:"type:jsonb;default:'[]'" json:"modifiers"`
	CreatedAt time.Time 	`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time 	`gorm:"autoUpdateTime" json:"updated_at"`
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
