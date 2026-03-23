package product_size

import (
	"time"

	"gorm.io/gorm"
)

type ProductSize struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	ProductID  uint   `gorm:"index;not null"`

	SKU             string  `gorm:"size:128;index"`
	SizeCode        string  `gorm:"size:64;index"`
	SizeName        string  `gorm:"size:255"`
	MeasureUnitType string  `gorm:"size:32;default:GRAM"`
	ButtonImageURL  string  `gorm:"size:1024"`
	Weight          float64 `gorm:"type:numeric(10,3);default:0"`
	Price           float64 `gorm:"type:numeric(12,2);default:0"`
	IsDefault       bool    `gorm:"default:false"`
	IsHidden        bool    `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

