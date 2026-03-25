package customer_tag

import (
	"time"

	"gorm.io/gorm"
)

type CustomerTag struct {
	ID uint `gorm:"primaryKey"`

	ExternalID         string `gorm:"size:128;uniqueIndex;not null"`
	CustomerTagGroupID uint   `gorm:"index;not null"`
	Name               string `gorm:"size:255;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

