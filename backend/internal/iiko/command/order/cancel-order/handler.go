package cancel_order

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderCanceler interface {
	Cancel(ctx context.Context, request order.CancelOrderRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo OrderCanceler
}

func NewHandler(repo OrderCanceler) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	orderID := strings.TrimSpace(command.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("не передан orderId для отмены заказа iiko")
	}

	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для отмены заказа iiko")
	}

	var removalTypeID *string
	if value := strings.TrimSpace(command.RemovalTypeID); value != "" {
		removalTypeID = &value
	}

	var removalComment *string
	if value := strings.TrimSpace(command.RemovalComment); value != "" {
		removalComment = &value
	}

	var userIDForWriteoff *string
	if value := strings.TrimSpace(command.UserIDForWriteoff); value != "" {
		userIDForWriteoff = &value
	}

	return h.repo.Cancel(ctx, order.CancelOrderRequest{
		OrderID:           orderID,
		OrganizationID:    organizationID,
		RemovalTypeID:     removalTypeID,
		RemovalComment:    removalComment,
		UserIDForWriteoff: userIDForWriteoff,
	})
}
