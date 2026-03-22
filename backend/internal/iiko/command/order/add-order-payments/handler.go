package add_order_payments

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderPaymentsAdder interface {
	AddPayments(ctx context.Context, request order.AddOrderPaymentsRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo OrderPaymentsAdder
}

func NewHandler(repo OrderPaymentsAdder) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	orderID := strings.TrimSpace(command.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("не передан orderId для добавления оплат в заказ iiko")
	}

	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для добавления оплат в заказ iiko")
	}

	if len(command.Payments) == 0 {
		return nil, fmt.Errorf("не переданы payments для добавления оплат в заказ iiko")
	}

	payments := make([]order.Payment, 0, len(command.Payments))
	for _, p := range command.Payments {
		paymentTypeKind := strings.TrimSpace(p.PaymentTypeKind)
		if paymentTypeKind == "" {
			return nil, fmt.Errorf("у одного из payments не передан paymentTypeKind")
		}
		paymentTypeID := strings.TrimSpace(p.PaymentTypeID)
		if paymentTypeID == "" {
			return nil, fmt.Errorf("у одного из payments не передан paymentTypeId")
		}
		if p.Sum < 0 {
			return nil, fmt.Errorf("у одного из payments некорректная сумма")
		}

		payments = append(payments, order.Payment{
			PaymentTypeKind:        paymentTypeKind,
			Sum:                    p.Sum,
			PaymentTypeID:          paymentTypeID,
			IsProcessedExternally:  p.IsProcessedExternally,
			IsFiscalizedExternally: p.IsFiscalizedExternally,
		})
	}

	tips := make([]order.TipsPayment, 0, len(command.Tips))
	for _, t := range command.Tips {
		paymentTypeKind := strings.TrimSpace(t.PaymentTypeKind)
		if paymentTypeKind == "" {
			return nil, fmt.Errorf("у одного из tips не передан paymentTypeKind")
		}
		tipsTypeID := strings.TrimSpace(t.TipsTypeID)
		if tipsTypeID == "" {
			return nil, fmt.Errorf("у одного из tips не передан tipsTypeId")
		}
		paymentTypeID := strings.TrimSpace(t.PaymentTypeID)
		if paymentTypeID == "" {
			return nil, fmt.Errorf("у одного из tips не передан paymentTypeId")
		}
		if t.Sum < 0 {
			return nil, fmt.Errorf("у одного из tips некорректная сумма")
		}

		tips = append(tips, order.TipsPayment{
			PaymentTypeKind:        paymentTypeKind,
			TipsTypeID:             tipsTypeID,
			Sum:                    t.Sum,
			PaymentTypeID:          paymentTypeID,
			IsProcessedExternally:  t.IsProcessedExternally,
			IsFiscalizedExternally: t.IsFiscalizedExternally,
		})
	}

	return h.repo.AddPayments(ctx, order.AddOrderPaymentsRequest{
		OrderID:        orderID,
		OrganizationID: organizationID,
		Payments:       payments,
		Tips:           tips,
	})
}
