package modifier_allergen

import "time"

type ModifierAllergen struct {
	ID uint `gorm:"primaryKey"`

	ModifierItemID  uint `gorm:"index:idx_modifier_allergen,unique;not null"`
	AllergenGroupID uint `gorm:"index:idx_modifier_allergen,unique;not null"`

	CreatedAt time.Time
}

