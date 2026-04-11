package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	iikoOrderService "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/service"
)

type IikoOrdersByRevisionAdapter struct {
	provider *iikoOrderService.OrdersByRevisionProvider
}

func NewIikoOrdersByRevisionAdapter(provider *iikoOrderService.OrdersByRevisionProvider) *IikoOrdersByRevisionAdapter {
	return &IikoOrdersByRevisionAdapter{provider: provider}
}

func (a *IikoOrdersByRevisionAdapter) GetOrdersByRevision(ctx context.Context, startRevision int64, sourceKeys []string) (*OrdersByRevisionPayload, error) {
	response, err := a.provider.GetOrdersByRevision(ctx, startRevision, sourceKeys)
	if err != nil {
		return nil, err
	}

	payload := &OrdersByRevisionPayload{
		MaxRevision: response.MaxRevision,
		Orders:      make([]OrdersByRevisionOrder, 0),
	}
	for _, orgOrders := range response.OrdersByOrganizations {
		for _, order := range orgOrders.Orders {
			orderID, err := uuid.Parse(order.ID)
			if err != nil {
				continue
			}

			isDeleted := order.Order.IsDeleted != nil && *order.Order.IsDeleted
			estimated := 0
			for _, item := range order.Order.Items {
				if item.Amount <= 0 {
					continue
				}
				estimated += int(item.Amount * 5)
			}
			if estimated == 0 {
				estimated = 5
			}

			var completeBeforeAt *time.Time
			if order.Order.CompleteBefore != "" {
				if parsed, err := time.Parse("2006-01-02 15:04:05.000", order.Order.CompleteBefore); err == nil {
					completeBeforeAt = &parsed
				}
			}

			payload.Orders = append(payload.Orders, OrdersByRevisionOrder{
				OrderID:          orderID,
				Status:           order.Order.Status,
				IsDeleted:        isDeleted,
				EstimatedCookMins: estimated,
				CompleteBeforeAt: completeBeforeAt,
			})
		}
	}

	return payload, nil
}
