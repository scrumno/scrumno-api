package interfaces

import (
	"context"

	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type CustomerProvider interface {
	GetCustomer(ctx context.Context, builderBody any) (any, error)
	SetCustomer(ctx context.Context, builderBody any) (any, error)
}

type CustomerBodyBuilder interface {
	BuildGet(ctx context.Context, phone string) any
	BuildSetFromUser(ctx context.Context, u *user.User) any
}

type CustomerSyncService interface {
	Sync(ctx context.Context, u *user.User) error
	SyncGet(ctx context.Context, phone string) (any, error)
}
