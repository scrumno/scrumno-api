package create_order

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderCreator interface {
	Create(ctx context.Context, request order.CreateOrderRequest) (*order.OrderResponse, error)
}

type Handler struct {
	repo OrderCreator
}

func NewHandler(repo OrderCreator) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.OrderResponse, error) {
	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для создания заказа iiko")
	}

	terminalGroupID := strings.TrimSpace(command.TerminalGroupID)
	if terminalGroupID == "" {
		return nil, fmt.Errorf("не передан terminalGroupId для создания заказа iiko")
	}

	if len(command.Items) == 0 {
		return nil, fmt.Errorf("не переданы позиции заказа для создания заказа iiko")
	}

	items := make([]order.ProductOrderItem, 0, len(command.Items))
	for _, item := range command.Items {
		productID := strings.TrimSpace(item.ProductID)
		if productID == "" {
			return nil, fmt.Errorf("одна из позиций заказа не содержит productId")
		}
		if item.Amount <= 0 {
			return nil, fmt.Errorf("некорректное количество для продукта %s", productID)
		}
		if item.Price < 0 {
			return nil, fmt.Errorf("некорректная цена для продукта %s", productID)
		}

		var itemComment *string
		if trimmedComment := strings.TrimSpace(item.Comment); trimmedComment != "" {
			itemComment = &trimmedComment
		}

		items = append(items, order.ProductOrderItem{
			Type:      "Product",
			ProductID: productID,
			Amount:    item.Amount,
			Price:     item.Price,
			Comment:   itemComment,
		})
	}

	var customer *order.Customer
	phone := strings.TrimSpace(command.CustomerPhone)
	customerName := strings.TrimSpace(command.CustomerName)
	if phone != "" || customerName != "" {
		customer = &order.Customer{
			Type: "one-time",
		}
		if phone != "" {
			customer.Phone = &phone
		}
		if customerName != "" {
			customer.Name = &customerName
		}
	}

	var orderPhone *string
	if phone != "" {
		orderPhone = &phone
	}

	return h.repo.Create(ctx, order.CreateOrderRequest{
		OrganizationID:  organizationID,
		TerminalGroupID: &terminalGroupID,
		Order: order.TableOrder{
			Phone:    orderPhone,
			Customer: customer,
			Items:    items,
		},
	})
}
