package create_order

import (
	"context"
	"strings"

	"github.com/google/uuid"
	iikoOrderModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	sharedOrder "github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/model/order"
	userEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	cartEntity "github.com/scrumno/scrumno-api/internal/cart/entity"
)

type Handler struct {
	orderProvider interfaces.OrderProvider
	orderBuilder  interfaces.OrderBodyBuilder
	userReader    UserReader
	cartReader    CartReader
}

type UserReader interface {
	FindByPhone(ctx context.Context, phone string) (*userEntity.User, error)
}

type CartReader interface {
	GetCartByUserId(ctx context.Context, userID uuid.UUID) (*cartEntity.Cart, error)
}

func NewHandler(orderProvider interfaces.OrderProvider, orderBuilder interfaces.OrderBodyBuilder, userReader UserReader, cartReader CartReader) *Handler {
	return &Handler{
		orderProvider: orderProvider,
		orderBuilder:  orderBuilder,
		userReader:    userReader,
		cartReader:    cartReader,
	}
}

func (h *Handler) Handle(ctx context.Context, phone string, comment *string) OrderDTO {
	u, err := h.userReader.FindByPhone(ctx, phone)
	if err != nil || u == nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     "Пользователь не найден",
		}
	}

	cart, err := h.cartReader.GetCartByUserId(ctx, u.ID)
	if err != nil || cart == nil {
		return OrderDTO{
			IsSuccess: false,
			Error:     "Корзина не найдена",
		}
	}

	if len(cart.Items) == 0 {
		return OrderDTO{
			IsSuccess: false,
			Error:     "Корзина пуста",
		}
	}

	var sourceKey string
	sourceKey = "app"

	var discountInfo []sharedOrder.DiscountsInfo
	discountInfo = make([]sharedOrder.DiscountsInfo, 0, 10)

	payments := []sharedOrder.Payment{}
	combos := []sharedOrder.Combos{}

	input := &sharedOrder.BuildInput{
		Customer: &sharedOrder.Customer{
			CustomerType: sharedOrder.Regular,
			Phone:        normalizePhone(u.Phone),
			Name:         "ss",
			ID:           uuid.Max,
		},
		Items:        make([]sharedOrder.BuildItem, 0, len(cart.Items)),
		Combos:       &combos,
		Payment:      &payments,
		DiscountInfo: &discountInfo,
		SourceKey:    &sourceKey,
		Comment:      *comment,
	}

	for _, item := range cart.Items {
		input.Items = append(input.Items, sharedOrder.BuildItem{
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			Price:     item.TotalPrice,
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

	return OrderDTO{
		IsSuccess: true,
		OrderID:   orderID.String(),
		Response:  resp,
		Error:     "",
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
	IsSuccess bool   `json:"isSuccess"`
	OrderID   string `json:"orderId"`
	Response  any    `json:"response,omitempty"`
	Error     string `json:"error,omitempty"`
}
