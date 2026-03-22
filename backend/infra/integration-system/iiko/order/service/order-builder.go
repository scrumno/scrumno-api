package service

import "github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"

type orderBuilder struct{}

func NewOrderBuilder() interfaces.OrderBuilder {
	return &orderBuilder{}
}

func (b *orderBuilder) BuildBody(data *any) *any {
	return nil
}
