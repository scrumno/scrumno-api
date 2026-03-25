package combo_group

import (
	"time"

	"gorm.io/gorm"
)

type ComboGroup struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	ComboID    uint   `gorm:"index;not null"`

	Name        string `gorm:"size:255;not null"`
	IsMainGroup bool   `gorm:"default:false"`
	SkipStep    bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

