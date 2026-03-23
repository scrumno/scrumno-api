package product_allergen

import "time"

type ProductAllergen struct {
	ID uint `gorm:"primaryKey"`

	ProductID       uint `gorm:"index:idx_product_allergen,unique;not null"`
	AllergenGroupID uint `gorm:"index:idx_product_allergen,unique;not null"`

	CreatedAt time.Time
}

