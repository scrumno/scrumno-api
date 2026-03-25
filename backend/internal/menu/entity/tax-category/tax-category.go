package tax_category

import (
	"time"

	"gorm.io/gorm"
)

type TaxCategory struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string  `gorm:"size:128;uniqueIndex;not null"`
	Name       string  `gorm:"size:255;not null"`
	Percentage float64 `gorm:"type:numeric(5,2);default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

