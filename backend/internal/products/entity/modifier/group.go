package modifier

type ProductModifierGroup struct {
	ID                                   string                 `json:"id" gorm:"primaryKey;size:128"`
	MinAmount                            int32                  `json:"minAmount"`
	MaxAmount                            int32                  `json:"maxAmount"`
	Required                             *bool                  `json:"required"`
	ChildModifiersHaveMinMaxRestrictions *bool                  `json:"childModifiersHaveMinMaxRestrictions"`
	ChildModifiers                       []ProductChildModifier `json:"childModifiers" gorm:"-"`
}
