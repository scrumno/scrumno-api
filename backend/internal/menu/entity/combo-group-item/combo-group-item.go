package combo_group_item

import (
	"time"

	"gorm.io/gorm"
)

type ComboGroupItem struct {
	ID uint `gorm:"primaryKey"`

	ComboGroupID uint `gorm:"index;not null"`
	ProductID    uint `gorm:"index;not null"`
	ProductSizeID *uint `gorm:"index"`

	PriceModificationAmount float64 `gorm:"type:numeric(12,2);default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

