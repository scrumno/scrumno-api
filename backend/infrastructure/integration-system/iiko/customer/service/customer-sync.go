package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/customer/model"
	interfaces "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/customer"
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

func (s *CustomerSyncService) Sync(ctx context.Context, u *user.User) (*customer.ResponseSet, error) {
	body := s.builder.BuildSetFromUser(ctx, u)
	resp, err := s.provider.SetCustomer(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("Создание пользователя в IIKO, : %w", err)
	}
	if resp == nil {
		return &customer.ResponseSet{}, nil
	}

	cBody, ok := resp.(*model.CustomerSetResponse)
	if !ok || cBody == nil {
		return nil, fmt.Errorf("Создание пользователя в IIKO, не тот формат данных")
	}

	return &customer.ResponseSet{
		ID: cBody.ID,
	}, nil
}

func (s *CustomerSyncService) SyncGet(ctx context.Context, phone *string, id *uuid.UUID) (*customer.ResponseGet, error) {
	body := s.builder.BuildGet(ctx, phone, id)
	resp, err := s.provider.GetCustomer(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("Получение пользователя с IIKO, : %w", err)
	}
	if resp == nil {
		return nil, nil
	}

	cr, ok := resp.(*model.CustomerResponse)
	if !ok || cr == nil {
		return nil, fmt.Errorf("Получение пользователя в IIKO, не тот формат данных")
	}

	wallets := make([]customer.WalletBalance, 0, len(cr.WalletBalances))
	for _, w := range cr.WalletBalances {
		wallets = append(wallets, customer.WalletBalance{
			ID:      w.ID,
			Name:    w.Name,
			Type:    w.Type,
			Balance: w.Balance,
		})
	}

	return &customer.ResponseGet{
		ID:             cr.ID,
		Name:           cr.Name,
		MiddleName:     cr.MiddleName,
		Surname:        cr.Surname,
		Sex:            cr.Sex,
		Birthday:       cr.Birthday,
		Email:          cr.Email,
		IsDeleted:      cr.IsDeleted,
		WalletBalances: wallets,
	}, nil
}
