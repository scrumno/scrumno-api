package service

import "context"

type PaymentService interface {
	Pay(ctx context.Context, draftID string, amount float64) (bool, error)
}

type paymentStubService struct{}

func NewPaymentStubService() PaymentService {
	return &paymentStubService{}
}

func (s *paymentStubService) Pay(ctx context.Context, draftID string, amount float64) (bool, error) {
	return true, nil
}
