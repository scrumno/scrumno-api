package customer_tag_group

import (
	"time"

	"gorm.io/gorm"
)

type CustomerTagGroup struct {
	ID uint `gorm:"primaryKey"`

	ExternalID        string `gorm:"size:128;uniqueIndex;not null"`
	Name              string `gorm:"size:255;not null"`
	SelectSeveralTags bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

