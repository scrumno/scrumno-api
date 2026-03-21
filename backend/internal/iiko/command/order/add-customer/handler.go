package add_customer

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type OrderCustomerAdder interface {
	AddCustomer(ctx context.Context, request order.AddCustomerToOrderRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo OrderCustomerAdder
}

func NewHandler(repo OrderCustomerAdder) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для добавления клиента в заказ iiko")
	}

	orderID := strings.TrimSpace(command.OrderID)
	if orderID == "" {
		return nil, fmt.Errorf("не передан orderId для добавления клиента в заказ iiko")
	}

	customer := order.TableOrderCustomer{
		Birthdate:                             command.Customer.Birthdate,
		ShouldReceiveOrderStatusNotifications: command.Customer.ShouldReceiveOrderStatusNotifications,
	}

	if id := strings.TrimSpace(command.Customer.ID); id != "" {
		customer.ID = &id
	}
	if name := strings.TrimSpace(command.Customer.Name); name != "" {
		customer.Name = &name
	}
	if surname := strings.TrimSpace(command.Customer.Surname); surname != "" {
		customer.Surname = &surname
	}
	if comment := strings.TrimSpace(command.Customer.Comment); comment != "" {
		customer.Comment = &comment
	}
	if email := strings.TrimSpace(command.Customer.Email); email != "" {
		customer.Email = &email
	}
	if gender := strings.TrimSpace(command.Customer.Gender); gender != "" {
		customer.Gender = &gender
	}
	if phone := strings.TrimSpace(command.Customer.Phone); phone != "" {
		customer.Phone = &phone
	}

	if customer.ID == nil && customer.Name == nil && customer.Phone == nil {
		return nil, fmt.Errorf("для добавления клиента нужен customer.id или пара customer.name и customer.phone")
	}

	return h.repo.AddCustomer(ctx, order.AddCustomerToOrderRequest{
		OrganizationID: organizationID,
		OrderID:        orderID,
		Customer:       customer,
	})
}
