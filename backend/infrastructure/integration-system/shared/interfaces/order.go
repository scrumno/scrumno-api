package interfaces

import (
	"context"

	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/order"
)

type OrderProvider interface {
	SetOrder(ctx context.Context, builderBody any) (any, error)
}

type OrderBodyBuilder interface {
	BuildSetFromOrder(ctx context.Context, input *order.BuildInput) any
}

type OrderSyncService interface {
	Sync(ctx context.Context, input *order.BuildInput) (*order.ResponseSet, error)
}
