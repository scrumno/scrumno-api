package modifier_item

import (
	"time"

	"gorm.io/gorm"
)

type ModifierItem struct {
	ID uint `gorm:"primaryKey"`

	ExternalID      string `gorm:"size:128;uniqueIndex;not null"`
	ModifierGroupID uint   `gorm:"index;not null"`

	SKU                string  `gorm:"size:128;index"`
	Name               string  `gorm:"size:255;not null"`
	Description        string  `gorm:"type:text"`
	TaxCategoryID      string  `gorm:"size:128;index"`
	ProductCategoryID  string  `gorm:"size:128;index"`
	PaymentSubject     string  `gorm:"size:64"`
	PaymentSubjectCode string  `gorm:"size:32"`
	OuterEANCode       string  `gorm:"size:64"`
	MeasureUnitType    string  `gorm:"size:32;default:GRAM"`
	ButtonImageURL     string  `gorm:"size:1024"`
	Weight             float64 `gorm:"type:numeric(10,3);default:0"`
	Price              float64 `gorm:"type:numeric(12,2);default:0"`
	Position           int     `gorm:"default:0"`
	IsHidden           bool    `gorm:"default:false"`
	IsMarked           bool    `gorm:"default:false"`
	IndependentQty     bool    `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

