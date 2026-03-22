package order

type CreateOrderRequest struct {
	OrganizationID  string     `json:"organizationId"`
	TerminalGroupID *string    `json:"terminalGroupId,omitempty"`
	Order           TableOrder `json:"order"`
}

type AddOrderItemsRequest struct {
	OrderID        string             `json:"orderId"`
	OrganizationID string             `json:"organizationId"`
	Items          []ProductOrderItem `json:"items"`
}

type AddCustomerToOrderRequest struct {
	OrganizationID string             `json:"organizationId"`
	OrderID        string             `json:"orderId"`
	Customer       TableOrderCustomer `json:"customer"`
}

type AddOrderPaymentsRequest struct {
	OrderID        string        `json:"orderId"`
	OrganizationID string        `json:"organizationId"`
	Payments       []Payment     `json:"payments"`
	Tips           []TipsPayment `json:"tips,omitempty"`
}

type ChangeOrderPaymentsRequest struct {
	OrderID        string        `json:"orderId"`
	OrganizationID string        `json:"organizationId"`
	Payments       []Payment     `json:"payments"`
	Tips           []TipsPayment `json:"tips,omitempty"`
}

type CloseOrderRequest struct {
	OrderID        string `json:"orderId"`
	OrganizationID string `json:"organizationId"`
}

type CancelOrderRequest struct {
	OrderID           string  `json:"orderId"`
	OrganizationID    string  `json:"organizationId"`
	RemovalTypeID     *string `json:"removalTypeId,omitempty"`
	RemovalComment    *string `json:"removalComment,omitempty"`
	UserIDForWriteoff *string `json:"userIdForWriteoff,omitempty"`
}

type GetByIDRequest struct {
	OrganizationIDs []string `json:"organizationIds"`
	OrderIDs        []string `json:"orderIds,omitempty"`
	PosOrderIDs     []string `json:"posOrderIds,omitempty"`
}

type InitByPosOrderRequest struct {
	OrganizationID  string   `json:"organizationId"`
	TerminalGroupID string   `json:"terminalGroupId"`
	PosOrderIDs     []string `json:"posOrderIds"`
}

type TableOrder struct {
	ID              *string            `json:"id,omitempty"`
	ExternalNumber  *string            `json:"externalNumber,omitempty"`
	Phone           *string            `json:"phone,omitempty"`
	OrderTypeID     *string            `json:"orderTypeId,omitempty"`
	Customer        *Customer          `json:"customer,omitempty"`
	PriceCategoryID *string            `json:"priceCategoryId,omitempty"`
	Items           []ProductOrderItem `json:"items"`
}

type Customer struct {
	Type  string  `json:"type"`
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

type TableOrderCustomer struct {
	ID                                    *string `json:"id,omitempty"`
	Name                                  *string `json:"name,omitempty"`
	Surname                               *string `json:"surname,omitempty"`
	Comment                               *string `json:"comment,omitempty"`
	Birthdate                             *string `json:"birthdate,omitempty"`
	Email                                 *string `json:"email,omitempty"`
	ShouldReceiveOrderStatusNotifications *bool   `json:"shouldReceiveOrderStatusNotifications,omitempty"`
	Gender                                *string `json:"gender,omitempty"`
	Phone                                 *string `json:"phone,omitempty"`
}

type ProductOrderItem struct {
	Type      string  `json:"type"`
	ProductID string  `json:"productId"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	Comment   *string `json:"comment,omitempty"`
}

type Payment struct {
	PaymentTypeKind        string  `json:"paymentTypeKind"`
	Sum                    float64 `json:"sum"`
	PaymentTypeID          string  `json:"paymentTypeId"`
	IsProcessedExternally  *bool   `json:"isProcessedExternally,omitempty"`
	IsFiscalizedExternally *bool   `json:"isFiscalizedExternally,omitempty"`
}

type TipsPayment struct {
	PaymentTypeKind        string  `json:"paymentTypeKind"`
	TipsTypeID             string  `json:"tipsTypeId"`
	Sum                    float64 `json:"sum"`
	PaymentTypeID          string  `json:"paymentTypeId"`
	IsProcessedExternally  *bool   `json:"isProcessedExternally,omitempty"`
	IsFiscalizedExternally *bool   `json:"isFiscalizedExternally,omitempty"`
}

type OrderResponse struct {
	CorrelationID string    `json:"correlationId"`
	OrderInfo     OrderInfo `json:"orderInfo"`
}

type OrderInfo struct {
	ID             string     `json:"id"`
	PosID          *string    `json:"posId"`
	OrganizationID string     `json:"organizationId"`
	Timestamp      int64      `json:"timestamp"`
	CreationStatus string     `json:"creationStatus"`
	ErrorInfo      *ErrorInfo `json:"errorInfo"`
}

type ErrorInfo struct {
	Code        *string `json:"code"`
	Message     *string `json:"message"`
	Description *string `json:"description"`
}

type CorrelationIDResponse struct {
	CorrelationID string `json:"correlationId"`
}

type TableOrdersResponse struct {
	CorrelationID string           `json:"correlationId"`
	Orders        []map[string]any `json:"orders"`
}
