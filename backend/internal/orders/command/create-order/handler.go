package create_order

import (
	"context"
	"strings"

	"github.com/google/uuid"
	iikoOrderModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	sharedOrder "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/order"
)

type Handler struct {
	orderProvider interfaces.OrderProvider
	orderBuilder  interfaces.OrderBodyBuilder
}

type Provider interface {
	Handle(ctx context.Context, cmd Command) OrderDTO
}

func NewHandler(orderProvider interfaces.OrderProvider, orderBuilder interfaces.OrderBodyBuilder) *Handler {
	return &Handler{
		orderProvider: orderProvider,
		orderBuilder:  orderBuilder,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) OrderDTO {
	phone := strings.TrimSpace(cmd.CustomerPhone)
	if phone == "" {
		return OrderDTO{IsSuccess: false, Error: "укажите номер телефона"}
	}
	if len(cmd.CartItems) == 0 {
		return OrderDTO{IsSuccess: false, Error: "корзина пользователя пуста"}
	}

	customerName := strings.TrimSpace(cmd.CustomerFullName)
	if customerName == "" {
		customerName = normalizePhone(phone)
	}

	comment := ""
	if cmd.OrderComment != nil {
		comment = *cmd.OrderComment
	}

	sourceKey := strings.TrimSpace(cmd.SourceKey)
	if sourceKey == "" {
		sourceKey = "app"
	}

	var discountInfo []sharedOrder.DiscountsInfo
	discountInfo = make([]sharedOrder.DiscountsInfo, 0, 10)

	payments := []sharedOrder.Payment{}
	combos := []sharedOrder.Combos{}

	input := &sharedOrder.BuildInput{
		Customer: &sharedOrder.Customer{
			CustomerType: sharedOrder.Regular,
			Phone:        normalizePhone(phone),
			Name:         customerName,
			ID:           nil,
		},
		Items:        make([]sharedOrder.BuildItem, 0, len(cmd.CartItems)),
		Combos:       &combos,
		Payment:      &payments,
		DiscountInfo: &discountInfo,
		SourceKey:    &sourceKey,
		Comment:      comment,
	}

	for _, item := range cmd.CartItems {
		input.Items = append(input.Items, sharedOrder.BuildItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Comment:   item.Comment,
		})
	}

	request := h.orderBuilder.BuildSetFromOrder(ctx, input)

	resp, err := h.orderProvider.SetOrder(ctx, request)
	if err != nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     err.Error(),
		}
	}
	if resp == nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     "iiko вернул пустой ответ при создании заказа",
		}
	}

	created, ok := resp.(*iikoOrderModel.OrderSetResponse)
	if !ok || created == nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     "iiko вернул некорректный orderId",
		}
	}

	orderID := created.OrderID
	if orderID == uuid.Nil && created.OrderInfo != nil {
		orderID = created.OrderInfo.ID
	}
	if orderID == uuid.Nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     "iiko вернул некорректный orderId",
		}
	}

	creationStatus := ""
	if created.OrderInfo != nil {
		creationStatus = created.OrderInfo.CreationStatus
	}

	return OrderDTO{
		IsSuccess:      true,
		OrderID:        orderID.String(),
		CorrelationID:  created.CorrelationID.String(),
		CreationStatus: creationStatus,
		Response:       resp,
		Error:          "",
	}
}

func normalizePhone(phone string) string {
	p := strings.TrimSpace(phone)
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "+") {
		return p
	}
	return "+" + p
}

type OrderDTO struct {
	IsSuccess      bool   `json:"isSuccess"`
	OrderID        string `json:"orderId"`
	CorrelationID  string `json:"correlationId,omitempty"`
	CreationStatus string `json:"creationStatus,omitempty"`
	Response       any    `json:"response,omitempty"`
	Error          string `json:"error,omitempty"`
}
