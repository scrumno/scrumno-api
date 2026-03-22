package add_order_payments

type Command struct {
	OrderID        string
	OrganizationID string
	Payments       []Payment
	Tips           []TipsPayment
}

type Payment struct {
	PaymentTypeKind        string
	Sum                    float64
	PaymentTypeID          string
	IsProcessedExternally  *bool
	IsFiscalizedExternally *bool
}

type TipsPayment struct {
	PaymentTypeKind        string
	TipsTypeID             string
	Sum                    float64
	PaymentTypeID          string
	IsProcessedExternally  *bool
	IsFiscalizedExternally *bool
}
