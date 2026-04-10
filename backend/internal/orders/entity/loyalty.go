package entity

type LoyaltyInfo struct {
	Coupon                     string            `json:"coupon,omitempty"`
	ApplicableManualConditions []string          `json:"applicableManualConditions,omitempty"`
	DynamicDiscounts           []DynamicDiscount `json:"dynamicDiscounts,omitempty"`
}

type DynamicDiscount struct {
	ManualConditionID string  `json:"manualConditionId"`
	Sum               float64 `json:"sum"`
}
