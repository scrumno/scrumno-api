package product

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	CategoryID uint   `gorm:"index;not null"`

	SKU                string  `gorm:"size:128;index"`
	Name               string  `gorm:"size:255;not null"`
	Description        string  `gorm:"type:text"`
	ModifierSchemaID   string  `gorm:"size:128;index"`
	ModifierSchemaName string  `gorm:"size:255"`
	Type               string  `gorm:"size:32;default:DISH"`
	OrderItemType      string  `gorm:"size:32;default:Product"`
	MeasureUnit        string  `gorm:"size:32"`
	ProductCategoryID  string  `gorm:"size:128;index"`
	TaxCategoryID      string  `gorm:"size:128;index"`
	PaymentSubject     string  `gorm:"size:64"`
	PaymentSubjectCode string  `gorm:"size:32"`
	OuterEANCode       string  `gorm:"size:64"`
	Price              float64 `gorm:"type:numeric(12,2);default:0"`
	Sort               int     `gorm:"default:0"`
	Splittable         bool    `gorm:"default:false"`
	CanSetOpenPrice    bool    `gorm:"default:false"`
	UseBalanceForSell  bool    `gorm:"default:false"`
	IsMarked           bool    `gorm:"default:false"`
	IsHidden           bool    `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewBuild() *Product {
	return &Product{}
}
