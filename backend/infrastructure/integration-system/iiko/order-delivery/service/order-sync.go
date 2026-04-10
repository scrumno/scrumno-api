package service

import (
	"context"
	"fmt"

	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	interfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/order"
)

type OrderSyncService struct {
	builder  interfaces.OrderBodyBuilder
	provider interfaces.OrderProvider
}

func NewOrderSyncService(builder interfaces.OrderBodyBuilder, provider interfaces.OrderProvider) *OrderSyncService {
	return &OrderSyncService{
		builder:  builder,
		provider: provider,
	}
}

func (s *OrderSyncService) Sync(ctx context.Context, input *order.BuildInput) (*order.ResponseSet, error) {
	body := s.builder.BuildSetFromOrder(ctx, input)
	resp, err := s.provider.SetOrder(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("Создание заказа в IIKO, : %w", err)
	}
	if resp == nil {
		return &order.ResponseSet{}, nil
	}

	cBody, ok := resp.(*model.OrderSetResponse)
	if !ok || cBody == nil {
		return nil, fmt.Errorf("Создание заказа в IIKO, не тот формат данных")
	}

	return &order.ResponseSet{
		ID: cBody.OrderID,
	}, nil
}
