package combo

import (
	"time"

	"gorm.io/gorm"
)

type Combo struct {
	ID uint `gorm:"primaryKey"`

	ExternalID       string `gorm:"size:128;uniqueIndex;not null"`
	ComboCategoryID  uint   `gorm:"index;not null"`
	Name             string `gorm:"size:255;not null"`
	Description      string `gorm:"type:text"`
	PriceStrategy    string `gorm:"size:32;default:BY_COMPONENT"`
	ButtonImageURL   string `gorm:"size:1024"`
	ButtonImageHash  string `gorm:"size:255"`
	Price            float64 `gorm:"type:numeric(12,2);default:0"`
	StartDateISO8601 string  `gorm:"size:64"`
	EndDateISO8601   string  `gorm:"size:64"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

