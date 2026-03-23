package allergen_group

import (
	"time"

	"gorm.io/gorm"
)

type AllergenGroup struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	Code       string `gorm:"size:64;index"`
	Name       string `gorm:"size:255;not null"`
	IsDeleted  bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

