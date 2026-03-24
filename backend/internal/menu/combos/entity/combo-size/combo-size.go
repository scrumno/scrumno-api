package combo_size

import (
	"time"

	"gorm.io/gorm"
)

type ComboSize struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	ComboID    uint   `gorm:"index;not null"`

	Name            string `gorm:"size:255;not null"`
	ShortName       string `gorm:"size:64"`
	ButtonImageURL  string `gorm:"size:1024"`
	ButtonImageHash string `gorm:"size:255"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

