package close_order

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderCloser interface {
	Close(ctx context.Context, request order.CloseOrderRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo OrderCloser
}

func NewHandler(repo OrderCloser) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	orderID := strings.TrimSpace(command.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("не передан orderId для закрытия заказа iiko")
	}

	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для закрытия заказа iiko")
	}

	return h.repo.Close(ctx, order.CloseOrderRequest{
		OrderID:        orderID,
		OrganizationID: organizationID,
	})
}
