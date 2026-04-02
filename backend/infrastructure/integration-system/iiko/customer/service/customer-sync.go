package service

import (
	"context"
	"fmt"

	interfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type CustomerSyncService struct {
	builder  interfaces.CustomerBodyBuilder
	provider interfaces.CustomerProvider
}

func NewCustomerSyncService(builder interfaces.CustomerBodyBuilder, provider interfaces.CustomerProvider) *CustomerSyncService {
	return &CustomerSyncService{
		builder:  builder,
		provider: provider,
	}
}

func (s *CustomerSyncService) Sync(ctx context.Context, u *user.User) error {
	body := s.builder.BuildSetFromUser(ctx, u)
	if _, err := s.provider.SetCustomer(ctx, body); err != nil {
		return fmt.Errorf("Создание пользователя в IIKO, : %w", err)
	}
	return nil
}

func (s *CustomerSyncService) SyncGet(ctx context.Context, u *user.User) error {
	body := s.builder.BuildGet(ctx, u)
	if _, err := s.provider.GetCustomer(ctx, body); err != nil {
		return fmt.Errorf("Получение пользователя с IIKO, : %w", err)
	}
	return nil
}
