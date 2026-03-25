package combo_category

import (
	"time"

	"gorm.io/gorm"
)

type ComboCategory struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	Name       string `gorm:"size:255;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

