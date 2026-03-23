package modifier_group

import (
	"time"

	"gorm.io/gorm"
)

type ModifierGroup struct {
	ID uint `gorm:"primaryKey"`

	ExternalID    string `gorm:"size:128;uniqueIndex;not null"`
	ProductSizeID uint   `gorm:"index;not null"`

	Name                           string `gorm:"size:255;not null"`
	Description                    string `gorm:"type:text"`
	SKU                            string `gorm:"size:128;index"`
	MinQuantity                    int    `gorm:"default:0"`
	MaxQuantity                    int    `gorm:"default:1"`
	FreeQuantity                   int    `gorm:"default:0"`
	DefaultQuantity                int    `gorm:"default:0"`
	HideIfDefaultQuantity          bool   `gorm:"default:false"`
	ChildModifiersHaveRestrictions bool   `gorm:"default:false"`
	Splittable                     bool   `gorm:"default:false"`
	IsHidden                       bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

