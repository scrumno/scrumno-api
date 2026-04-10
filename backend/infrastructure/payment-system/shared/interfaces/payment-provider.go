package interfaces

import "context"

type PaymentProvider interface {
	CreatePayment(ctx context.Context, payment *any) (*any, error)
	GetPayment(ctx context.Context, paymentID string) (*any, error)
	Pay() bool
}

type PaymentBuilder interface {
	BuildPayment(ctx context.Context, payment *any) (*any, error)
}

type PaymentHandler interface {
	Handle(ctx context.Context, payment *any) (*any, error)
}
