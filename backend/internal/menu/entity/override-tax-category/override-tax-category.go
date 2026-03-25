package override_tax_category

import (
	"time"

	"gorm.io/gorm"
)

type OverrideTaxCategory struct {
	ID uint `gorm:"primaryKey"`

	OrderTypeExternalID       string `gorm:"size:128;index;not null"`
	BaseTaxCategoryExternalID string `gorm:"size:128;index;not null"`
	NewTaxCategoryExternalID  string `gorm:"size:128;index;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

