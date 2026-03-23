package create_order

import (
	"context"

	"github.com/scrumno/scrumno-api/infra/integration-system/iiko/order/model"
	"github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"
)

type Handler struct {
	orderProvider interfaces.OrderProvider
	orderBuilder  interfaces.OrderBuilder
}

func NewHandler(orderProvider interfaces.OrderProvider, orderBuilder interfaces.OrderBuilder) *Handler {
	return &Handler{
		orderProvider: orderProvider,
		orderBuilder:  orderBuilder,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) DTO {
	order := model.DeliveryOrder{}

	request := h.orderBuilder.BuildBody(order)

	_, err := h.orderProvider.Create(request)
	if err != nil {
		return DTO{
			IsSuccess: false,
			Error:     err.Error(),
		}
	}

	return DTO{
		IsSuccess: true,
		OrderID:   "123",
		Error:     "",
	}
}

type DTO struct {
	IsSuccess bool   `json:"isSuccess"`
	OrderID   string `json:"orderId"`
	Error     string `json:"error,omitempty"`
}

type Command struct {
}
