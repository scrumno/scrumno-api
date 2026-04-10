package entity

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

type OrderServiceType string

const (
	OrderServiceDeliveryByCourier OrderServiceType = "DeliveryByCourier"
	OrderServiceDeliveryByClient  OrderServiceType = "DeliveryByClient"
)

type ExternalDataEntry struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsPublic bool   `json:"isPublic,omitempty"`
}
