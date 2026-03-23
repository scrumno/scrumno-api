package model

// OrderServiceType — Deliveries.Request.CreateOrder.OrderServiceType (только один из orderTypeId / orderServiceType).
type OrderServiceType string

const (
	OrderServiceDeliveryByCourier OrderServiceType = "DeliveryByCourier"
	OrderServiceDeliveryByClient  OrderServiceType = "DeliveryByClient" // самовывоз клиентом
)

// CreateOrderRequest — Deliveries.Request.CreateOrderRequest (POST /api/1/deliveries/create).
// Обязательны organizationId и order; terminalGroupId и createOrderSettings — опциональны.
type CreateOrderRequest struct {
	OrganizationID      string               `json:"organizationId"`
	TerminalGroupID     string               `json:"terminalGroupId,omitempty"`
	Order               DeliveryOrder        `json:"order"`
	CreateOrderSettings *CreateOrderSettings `json:"createOrderSettings,omitempty"`
}

// CreateOrderSettings — Orders.Common.CreateOrderSettings (для доставки без servicePrint).
type CreateOrderSettings struct {
	TransportToFrontTimeout int  `json:"transportToFrontTimeout,omitempty"`
	CheckStopList           bool `json:"checkStopList,omitempty"`
}

// DeliveryOrder — Deliveries.Request.CreateOrder.DeliveryOrder.
// Для самовывоза: orderServiceType = DeliveryByClient, deliveryPoint не передаётся.
// Обязательны phone и items.
type DeliveryOrder struct {
	MenuID               *string               `json:"menuId,omitempty"`
	ID                   string                `json:"id,omitempty"`
	ExternalNumber       string                `json:"externalNumber,omitempty"`
	CompleteBefore       string                `json:"completeBefore,omitempty"`
	Phone                string                `json:"phone"`
	PhoneExtension       string                `json:"phoneExtension,omitempty"`
	OrderTypeID          string                `json:"orderTypeId,omitempty"`
	OrderServiceType     OrderServiceType      `json:"orderServiceType,omitempty"`
	DeliveryPoint        *DeliveryPoint        `json:"deliveryPoint,omitempty"`
	Comment              string                `json:"comment,omitempty"`
	Customer             *Customer             `json:"customer,omitempty"`
	Guests               *Guests               `json:"guests,omitempty"`
	MarketingSourceID    string                `json:"marketingSourceId,omitempty"`
	OperatorID           string                `json:"operatorId,omitempty"`
	DeliveryDuration     int                   `json:"deliveryDuration,omitempty"`
	DeliveryZone         string                `json:"deliveryZone,omitempty"`
	PriceCategoryID      string                `json:"priceCategoryId,omitempty"`
	Items                []OrderItem           `json:"items"`
	Combos               []OrderCombo          `json:"combos,omitempty"`
	Payments             []Payment             `json:"payments,omitempty"`
	Tips                 []Tip                 `json:"tips,omitempty"`
	SourceKey            string                `json:"sourceKey,omitempty"`
	DiscountsInfo        *DiscountsInfo        `json:"discountsInfo,omitempty"`
	LoyaltyInfo          *LoyaltyInfo          `json:"loyaltyInfo,omitempty"`
	ChequeAdditionalInfo *ChequeAdditionalInfo `json:"chequeAdditionalInfo,omitempty"`
	ExternalData         []ExternalDataEntry   `json:"externalData,omitempty"`
}

// DeliveryPoint — для курьерской доставки; для самовывоза обычно nil.
// address — полиморфный (legacy/city), при необходимости заполняйте через map/отдельные DTO.
type DeliveryPoint struct {
	Coordinates           *Coordinates `json:"coordinates,omitempty"`
	Address               any          `json:"address,omitempty"`
	ExternalCartographyID string       `json:"externalCartographyId,omitempty"`
	Comment               string       `json:"comment,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Customer — discriminator type: regular / one-time (см. iiko API).
type Customer struct {
	ID                                    string `json:"id,omitempty"`
	Name                                  string `json:"name,omitempty"`
	Surname                               string `json:"surname,omitempty"`
	Comment                               string `json:"comment,omitempty"`
	Birthdate                             string `json:"birthdate,omitempty"`
	Email                                 string `json:"email,omitempty"`
	ShouldReceivePromoActionsInfo         bool   `json:"shouldReceivePromoActionsInfo,omitempty"`
	ShouldReceiveOrderStatusNotifications bool   `json:"shouldReceiveOrderStatusNotifications,omitempty"`
	Gender                                string `json:"gender,omitempty"`
	Type                                  string `json:"type,omitempty"`
}

// Guests — count обязателен, если объект передан.
type Guests struct {
	Count               int   `json:"count"`
	SplitBetweenPersons *bool `json:"splitBetweenPersons,omitempty"`
}

type OrderItem struct {
	Type             string            `json:"type"`
	Amount           float64           `json:"amount"`
	ProductSizeID    string            `json:"productSizeId,omitempty"`
	ComboInformation *ComboInformation `json:"comboInformation,omitempty"`
	Comment          string            `json:"comment,omitempty"`
}

type ComboInformation struct {
	ComboID        string `json:"comboId"`
	ComboSourceID  string `json:"comboSourceId"`
	ComboGroupID   string `json:"comboGroupId"`
	ComboGroupName string `json:"comboGroupName,omitempty"`
}

type OrderCombo struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	SourceID  string  `json:"sourceId"`
	ProgramID string  `json:"programId,omitempty"`
	SizeID    string  `json:"sizeId,omitempty"`
}

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

type LoyaltyInfo struct {
	Coupon                     string            `json:"coupon,omitempty"`
	ApplicableManualConditions []string          `json:"applicableManualConditions,omitempty"`
	DynamicDiscounts           []DynamicDiscount `json:"dynamicDiscounts,omitempty"`
}

type DynamicDiscount struct {
	ManualConditionID string  `json:"manualConditionId"`
	Sum               float64 `json:"sum"`
}

type ChequeAdditionalInfo struct {
	NeedReceipt       bool   `json:"needReceipt"`
	Email             string `json:"email,omitempty"`
	SettlementPlace   string `json:"settlementPlace,omitempty"`
	Phone             string `json:"phone,omitempty"`
	RetailAddress     string `json:"retailAddress,omitempty"`
	IsInternetPayment bool   `json:"isInternetPayment,omitempty"`
}

type ExternalDataEntry struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsPublic bool   `json:"isPublic,omitempty"`
}
