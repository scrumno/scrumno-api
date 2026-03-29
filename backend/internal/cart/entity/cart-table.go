package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ProductType string

const (
	ProductTypeProduct  ProductType = "Product"
	ProductTypeCompound ProductType = "Compound"
)

type Cart struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_carts_user" json:"user_id"`
	Items       []CartItem `gorm:"foreignKey:CartID" json:"items"`
	TotalAmount float64    `gorm:"type:numeric(12,2);not null" json:"total_amount"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type CartItem struct {
	ID               uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CartID           uuid.UUID       `gorm:"type:uuid;not null;index:idx_cart_items_cart" json:"cart_id"`
	ProductID        uuid.UUID       `gorm:"type:uuid;not null;index:idx_cart_items_product" json:"product_id"`
	Type             ProductType     `gorm:"type:varchar(20);not null;default:'Product'" json:"type"`
	Quantity         float64         `gorm:"type:numeric(10,3);not null;check:quantity >= 0" json:"quantity"` // amount (iiko: 0..999.999)
	Modifiers        *datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"modifiers"`
	CreatedAt        time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	Components       datatypes.JSON  `gorm:"type:jsonb;default:'{}'" json:"components"`
	CommonModifiers  *datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"common_modifiers,omitempty"`
	ComboInformation *datatypes.JSON `gorm:"type:jsonb" json:"combo_information,omitempty"`
	Comment          string          `gorm:"type:varchar(255)" json:"comment,omitempty"`
	SizeID           *uuid.UUID      `gorm:"type:uuid" json:"size_id,omitempty"`
	BasePrice        float64         `gorm:"type:numeric(12,2);not null;column:unit_price" json:"base_price"`
	TotalPrice       float64         `gorm:"type:numeric(12,2);not null" json:"total_price"`
	MeasureUnit      string          `gorm:"type:varchar(10);default:'шт'" json:"measure_unit"`
	IsWeighted       bool            `gorm:"default:false" json:"is_weighted"`
	Weight           float64         `gorm:"type:numeric(10,3)" json:"weight,omitempty"`
}

func NewCart(userID uuid.UUID) *Cart {
	return &Cart{
		ID:          uuid.New(),
		UserID:      userID,
		TotalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Items:       []CartItem{},
	}
}
