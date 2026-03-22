package add_order_items

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderItemsAdder interface {
	AddItems(ctx context.Context, request order.AddOrderItemsRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo OrderItemsAdder
}

func NewHandler(repo OrderItemsAdder) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	orderID := strings.TrimSpace(command.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("не передан orderId для добавления позиций заказа iiko")
	}

	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для добавления позиций заказа iiko")
	}

	if len(command.Items) == 0 {
		return nil, fmt.Errorf("не переданы позиции для добавления в заказ iiko")
	}

	items := make([]order.ProductOrderItem, 0, len(command.Items))
	for _, item := range command.Items {
		productID := strings.TrimSpace(item.ProductID)
		if productID == "" {
			return nil, fmt.Errorf("одна из добавляемых позиций не содержит productId")
		}
		if item.Amount <= 0 {
			return nil, fmt.Errorf("некорректное количество для продукта %s", productID)
		}
		if item.Price < 0 {
			return nil, fmt.Errorf("некорректная цена для продукта %s", productID)
		}

		var comment *string
		if trimmedComment := strings.TrimSpace(item.Comment); trimmedComment != "" {
			comment = &trimmedComment
		}

		items = append(items, order.ProductOrderItem{
			Type:      "Product",
			ProductID: productID,
			Amount:    item.Amount,
			Price:     item.Price,
			Comment:   comment,
		})
	}

	return h.repo.AddItems(ctx, order.AddOrderItemsRequest{
		OrderID:        orderID,
		OrganizationID: organizationID,
		Items:          items,
	})
}
