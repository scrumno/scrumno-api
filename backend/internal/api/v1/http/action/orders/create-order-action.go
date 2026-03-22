package orders

import (
	"net/http"

	"github.com/scrumno/scrumno-api/infra/integration-system/shared/interfaces"
)

type CreateOrderAction struct {
	OrderProvider interfaces.OrderProvider
	OrderBuilder  interfaces.OrderBuilder
}

func NewCreateOrderAction(
	orderProvider interfaces.OrderProvider,
	orderBuilder interfaces.OrderBuilder,
) *CreateOrderAction {
	return &CreateOrderAction{
		OrderProvider: orderProvider,
		OrderBuilder:  orderBuilder,
	}
}

func (a *CreateOrderAction) Action(w http.ResponseWriter, r *http.Request) {
	return
}
