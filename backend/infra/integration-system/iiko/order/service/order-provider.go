package service

import "github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"

type orderProvider struct{}

func NewOrderProvider() interfaces.OrderProvider {
	return &orderProvider{}
}

func (p *orderProvider) Create(order *any) error {
	return nil
}
