package service

import (
	"context"

	"github.com/scrumno/scrumno-api/infrastructure/payment-system/shared/interfaces"
)

type paymentProvider struct {
	interfaces.PaymentProvider
}

func NewPaymentProvider(provider interfaces.PaymentProvider) *paymentProvider {
	return &paymentProvider{provider}
}

func (p *paymentProvider) CreatePayment(ctx context.Context, payment *any) (*any, error) {
	return p.PaymentProvider.CreatePayment(ctx, payment)
}

func (p *paymentProvider) GetPayment(ctx context.Context, paymentID string) (*any, error) {
	return p.PaymentProvider.GetPayment(ctx, paymentID)
}

func (p *paymentProvider) Pay() bool {
	// TODO: implement pay logic
	return true
}
