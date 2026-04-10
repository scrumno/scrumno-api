package product

import (
	"time"

	"github.com/google/uuid"
)

type CookingTimeProductTable struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ProductID   uint      `gorm:"not null;index:idx_cooking_time_product_table_product_id" json:"product_id"`
	CookingTime int       `gorm:"not null;default:0" json:"cooking_time"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
