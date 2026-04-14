package service

import (
	"context"

	"github.com/google/uuid"
	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	sharedOrder "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/order"
)

type SearchableType string

const (
	Id SearchableType = "id"
)

type OrderBodyBuilder struct {
	config *iikoConfig.Config
}

func NewOrderBodyBuilder(config *iikoConfig.Config) *OrderBodyBuilder {
	return &OrderBodyBuilder{
		config: config,
	}
}

type BuilderCommand struct {
	IdValue        *uuid.UUID     `json:"id,omitempty"`
	Type           SearchableType `json:"type"`
	OrganizationId uuid.UUID      `json:"organizationId"`
}

type BuilderSetCommand struct {
	OrganizationID      uuid.UUID                  `json:"organizationId"`
	TerminalGroupID     uuid.UUID                  `json:"terminalGroupId,omitempty"`
	CreateOrderSettings *model.CreateOrderSettings `json:"createOrderSettings,omitempty"`
	Order               model.DeliveryOrder        `json:"order"`
}

func (h *OrderBodyBuilder) BuildSetFromOrder(ctx context.Context, input *sharedOrder.BuildInput) any {
	if input == nil {
		return &BuilderSetCommand{}
	}

	items := make([]model.OrderItem, 0, len(input.Items))
	for _, it := range input.Items {
		items = append(items, model.OrderItem{
			Type:      "Product",
			Amount:    it.Quantity,
			ProductID: it.ProductID,
			Comment:   it.Comment,
			Price:     it.Price,
		})
	}

	var customer *model.Customer
	if input.Customer != nil {
		customerID := ""
		if input.Customer.ID != nil {
			customerID = input.Customer.ID.String()
		}
		customer = &model.Customer{
			Type: string(input.Customer.CustomerType),
			Name: input.Customer.Name,
			ID:   customerID,
		}
	}

	var combosFinal *[]model.OrderCombo
	if input.Combos != nil {
		combos := make([]model.OrderCombo, 0, len(*input.Combos))
		for _, it := range *input.Combos {
			var programID string
			if it.ProgramID != nil {
				programID = it.ProgramID.String()
			}

			var sizeID string
			if it.SizeID != nil {
				sizeID = it.SizeID.String()
			}

			combos = append(combos, model.OrderCombo{
				ID:        it.Id.String(),
				Name:      it.Name,
				Amount:    float64(it.Amount),
				Price:     it.Price,
				SourceID:  it.SourceID.String(),
				ProgramID: programID,
				SizeID:    sizeID,
			})
		}

		combosFinal = &combos
	}

	var paymmentsFinal *[]model.Payment
	if input.Payment != nil {
		paymments := make([]model.Payment, 0, len(*input.Payment))
		for _, payment := range *input.Payment {
			var paymentAdditionalData *model.PaymentAdditionalData
			if payment.PaymentAdditionalData != nil {
				paymentAdditionalData = &model.PaymentAdditionalData{
					Type: payment.PaymentAdditionalData.Type,
				}
			}

			paymments = append(paymments, model.Payment{
				PaymentTypeKind:        payment.PaymentTypeID,
				Sum:                    payment.Sum,
				PaymentTypeID:          payment.PaymentTypeID,
				IsProcessedExternally:  payment.IsProcessedExternally,
				PaymentAdditionalData:  paymentAdditionalData,
				IsFiscalizedExternally: payment.IsFiscalizedExternally,
				IsPrepay:               payment.IsPrepay,
			})
		}

		paymmentsFinal = &paymments
	}

	var discountFinal *model.DiscountsInfo
	if input.DiscountInfo != nil && len(*input.DiscountInfo) > 0 {
		discount := (*input.DiscountInfo)[0]
		discountData := model.DiscountsInfo{
			FixedLoyaltyDiscounts: discount.FixedLoyaltyDiscounts,
		}
		if discount.Card != nil {
			discountData.Card = &model.Card{
				Track: discount.Card.Track,
			}
		}

		discountInfo := make([]model.DiscountEntry, 0, len(discount.Discounts))
		for _, discountInfoItem := range discount.Discounts {
			discountInfo = append(discountInfo, model.DiscountEntry{
				Type: discountInfoItem.Type,
			})
		}
		discountData.Discounts = discountInfo

		discountFinal = &discountData
	}

	orderPhone := ""
	if input.Customer != nil {
		orderPhone = input.Customer.Phone
	}

	cmd := &BuilderSetCommand{
		OrganizationID:  h.config.OrganizationID,
		TerminalGroupID: h.config.TerminalGroupID,
		CreateOrderSettings: &model.CreateOrderSettings{
			CheckStopList: true,
		},
		Order: model.DeliveryOrder{
			OrderServiceType: model.OrderServiceDeliveryByClient,
			Phone:            orderPhone,
			Items:            items,
			Customer:         customer,
			Comment:          input.Comment,
		},
	}
	if input.SourceKey != nil {
		cmd.Order.SourceKey = *input.SourceKey
	}

	if combosFinal != nil {
		cmd.Order.Combos = combosFinal
	}

	if paymmentsFinal != nil {
		cmd.Order.Payments = paymmentsFinal
	}

	if discountFinal != nil {
		cmd.Order.DiscountsInfo = discountFinal
	}

	return cmd
}
