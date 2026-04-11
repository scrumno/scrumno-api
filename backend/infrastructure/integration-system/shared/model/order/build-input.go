package order

import (
	"github.com/google/uuid"
)

type CustomerType string

const (
	OneTime CustomerType = "one-time"
	Regular CustomerType = "regular"
)

type BuildInput struct {
	Comment      string
	Customer     *Customer
	Items        []BuildItem
	Combos       *[]Combos
	Payment      *[]Payment
	DiscountInfo *[]DiscountsInfo
	SourceKey    *string
}

type Customer struct {
	ID           *uuid.UUID
	Name         string
	Surnmae      *string
	CustomerType CustomerType
	Phone        string
}

type Combos struct {
	Id        uuid.UUID
	Name      string
	Amount    int32
	Price     float64
	SourceID  uuid.UUID
	ProgramID *uuid.UUID
	SizeID    *uuid.UUID
}

type PaymentAdditionalData struct {
	Type string `json:"type"`
}

type Payment struct {
	PaymentTypeKind        string
	Sum                    float64
	PaymentTypeID          string
	IsProcessedExternally  bool
	PaymentAdditionalData  *PaymentAdditionalData
	IsFiscalizedExternally bool
	IsPrepay               bool
}

type BuildItem struct {
	ProductID string
	Quantity  float64
	Comment   string
	Price     float64
	Type      string
}

type DiscountsInfo struct {
	Card                  *Card
	Discounts             []DiscountEntry
	FixedLoyaltyDiscounts bool
}

type Card struct {
	Track string
}

type DiscountEntry struct {
	DiscountType       string
	Sum                float64
	SelectivePositions []string
	Type               string
}
