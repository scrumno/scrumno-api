package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/customer"
	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type CustomerProvider interface {
	GetCustomer(ctx context.Context, builderBody any) (any, error)
	SetCustomer(ctx context.Context, builderBody any) (any, error)
}

type CustomerBodyBuilder interface {
	BuildGet(ctx context.Context, phone *string, id *uuid.UUID) any
	BuildSetFromUser(ctx context.Context, u *user.User) any
}

type CustomerSyncService interface {
	Sync(ctx context.Context, u *user.User) (*customer.ResponseSet, error)
	SyncGet(ctx context.Context, phone *string, id *uuid.UUID) (*customer.ResponseGet, error)
	// SyncWalletBalance(ctx context.Context, walletId uuid.UUID, balance float64)
}
