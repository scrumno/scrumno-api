package model

import "github.com/google/uuid"

// OrderServiceType — Deliveries.Request.CreateOrder.OrderServiceType (только один из orderTypeId / orderServiceType).
type OrderServiceType string

const (
	OrderServiceDeliveryByCourier OrderServiceType = "DeliveryByCourier"
	OrderServiceDeliveryByClient  OrderServiceType = "DeliveryByClient" // самовывоз клиентом
)

type CreateOrderRequest struct {
	OrganizationID      string               `json:"organizationId"`
	TerminalGroupID     string               `json:"terminalGroupId,omitempty"`
	CreateOrderSettings *CreateOrderSettings `json:"createOrderSettings,omitempty"`
	Order               DeliveryOrder        `json:"order"`
}

type CreateOrderSettings struct {
	TransportToFrontTimeout int  `json:"transportToFrontTimeout,omitempty"`
	CheckStopList           bool `json:"checkStopList,omitempty"`
}

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
	Combos               *[]OrderCombo         `json:"combos,omitempty"`
	Payments             *[]Payment            `json:"payments,omitempty"`
	Tips                 *[]Tip                `json:"tips,omitempty"`
	SourceKey            string                `json:"sourceKey,omitempty"`
	DiscountsInfo        *DiscountsInfo        `json:"discountsInfo,omitempty"`
	LoyaltyInfo          *LoyaltyInfo          `json:"loyaltyInfo,omitempty"`
	ChequeAdditionalInfo *ChequeAdditionalInfo `json:"chequeAdditionalInfo,omitempty"`
	ExternalData         []ExternalDataEntry   `json:"externalData,omitempty"`
}

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

type Guests struct {
	Count               int   `json:"count"`
	SplitBetweenPersons *bool `json:"splitBetweenPersons,omitempty"`
}

type OrderItem struct {
	Type             string            `json:"type"`
	Amount           float64           `json:"amount"`
	ProductID        string            `json:"productId"`
	ProductSizeID    string            `json:"productSizeId,omitempty"`
	ComboInformation *ComboInformation `json:"comboInformation,omitempty"`
	Comment          string            `json:"comment,omitempty"`
	Price            float64           `json:"price"`
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

type OrderSetResponse struct {
	OrderID   uuid.UUID     `json:"orderId"`
	OrderInfo *OrderInfoDTO `json:"orderInfo,omitempty"`
}

type OrderInfoDTO struct {
	ID             uuid.UUID `json:"id"`
	CreationStatus string    `json:"creationStatus"`
}

type OrderResponse struct {
	ID     uuid.UUID `json:"id"`
	Status string    `json:"status"`
}
