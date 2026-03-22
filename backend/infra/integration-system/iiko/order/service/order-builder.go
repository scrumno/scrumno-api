package service

import (
	"github.com/scrumno/scrumno-api/infra/integration-system/iiko/config"
	"github.com/scrumno/scrumno-api/infra/integration-system/iiko/order/model"
	"github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"
)

type orderBuilder struct{}

func NewOrderBuilder() interfaces.OrderBuilder {
	return &orderBuilder{}
}

func (b *orderBuilder) BuildBody(data any) any {
	order, ok := data.(*model.DeliveryOrder)
	if !ok {
		return nil
	}

	cfg := config.Load()

	return model.CreateOrderRequest{
		OrganizationID:  cfg.OrganizationID,
		TerminalGroupID: cfg.TerminalGroupID,
		Order:           *order,
		CreateOrderSettings: &model.CreateOrderSettings{
			TransportToFrontTimeout: 10,
			CheckStopList:           false,
		},
	}
}
