package modifier

type ProductModifier struct {
	ID                  string `json:"id" gorm:"primaryKey;size:128"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}

type ProductChildModifier struct {
	ID                  string `json:"id" gorm:"primaryKey;size:128"`
	DefaultAmount       *int32 `json:"defaultAmount"`
	MinAmount           int32  `json:"minAmount"`
	MaxAmount           int32  `json:"maxAmount"`
	Required            *bool  `json:"required"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount"`
	Splittable          *bool  `json:"splittable"`
	FreeOfChargeAmount  *int32 `json:"freeOfChargeAmount"`
}
