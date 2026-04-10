package entity

type Payment struct {
	PaymentTypeKind        string                 `json:"paymentTypeKind"`
	Sum                    float64                `json:"sum"`
	PaymentTypeID          string                 `json:"paymentTypeId"`
	IsProcessedExternally  bool                   `json:"isProcessedExternally,omitempty"`
	PaymentAdditionalData  *PaymentAdditionalData `json:"paymentAdditionalData,omitempty"`
	IsFiscalizedExternally bool                   `json:"isFiscalizedExternally,omitempty"`
	IsPrepay               bool                   `json:"isPrepay,omitempty"`
}

type Tip struct {
	PaymentTypeKind        string                 `json:"paymentTypeKind"`
	TipsTypeID             string                 `json:"tipsTypeId,omitempty"`
	Sum                    float64                `json:"sum"`
	PaymentTypeID          string                 `json:"paymentTypeId"`
	IsProcessedExternally  bool                   `json:"isProcessedExternally,omitempty"`
	PaymentAdditionalData  *PaymentAdditionalData `json:"paymentAdditionalData,omitempty"`
	IsFiscalizedExternally bool                   `json:"isFiscalizedExternally,omitempty"`
	IsPrepay               bool                   `json:"isPrepay,omitempty"`
}

type PaymentAdditionalData struct {
	Type string `json:"type"`
}

type DiscountsInfo struct {
	Card                  *Card           `json:"card,omitempty"`
	Discounts             []DiscountEntry `json:"discounts,omitempty"`
	FixedLoyaltyDiscounts bool            `json:"fixedLoyaltyDiscounts,omitempty"`
}

type Card struct {
	Track string `json:"track"`
}

type DiscountEntry struct {
	Type string `json:"type"`
}

type ChequeAdditionalInfo struct {
	NeedReceipt       bool   `json:"needReceipt"`
	Email             string `json:"email,omitempty"`
	SettlementPlace   string `json:"settlementPlace,omitempty"`
	Phone             string `json:"phone,omitempty"`
	RetailAddress     string `json:"retailAddress,omitempty"`
	IsInternetPayment bool   `json:"isInternetPayment,omitempty"`
}
