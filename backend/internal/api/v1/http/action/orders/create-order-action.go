package orders

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	authEntity "github.com/scrumno/scrumno-api/internal/authorize/entity"
	userEntity "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
	cartModel "github.com/scrumno/scrumno-api/internal/cart/entity"
	cartEntity "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
)

const createOrderSourceApp = "app"

type CreateOrderAction struct {
	Handler     *createOrder.Handler
	UserFetcher *userEntity.Fetcher
	CartFetcher *cartEntity.Fetcher
}

func NewCreateOrderAction(
	handler *createOrder.Handler,
	userFetcher *userEntity.Fetcher,
	cartFetcher *cartEntity.Fetcher,
) *CreateOrderAction {
	return &CreateOrderAction{
		Handler:     handler,
		UserFetcher: userFetcher,
		CartFetcher: cartFetcher,
	}
}

func (a *CreateOrderAction) GetInputType() reflect.Type {
	return reflect.TypeOf(CreateOrderRequest{})
}

type CreateOrderRequest struct {
	CustomerPhone string `json:"phone" example:"79999009999"`
	Comment       string `json:"comment" example:"Комментарий"`
}

type CreateOrderResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	OrderID   string `json:"orderId,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (a *CreateOrderAction) Action(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	phone := strings.TrimSpace(req.CustomerPhone)
	if phone == "" {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     "Укажите номер телефона",
		}, http.StatusBadRequest)
		return
	}

	u, err := a.UserFetcher.Fetch(r.Context(), phone)
	if err != nil {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if u == nil {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     "Пользователь не найден",
		}, http.StatusBadRequest)
		return
	}

	cart, err := a.CartFetcher.Fetch(r.Context(), cartEntity.Query{UserID: u.ID})
	if err != nil {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if cart == nil {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     "Корзина пользователя не найдена",
		}, http.StatusBadRequest)
		return
	}

	if len(cart.Items) == 0 {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     "Корзина пользователя пуста",
		}, http.StatusBadRequest)
		return
	}

	cmd := createOrder.Command{
		CustomerPhone:    u.Phone,
		CustomerFullName: customerDisplayName(u),
		SourceKey:        createOrderSourceApp,
		OrderComment:     optionalTrimmedString(req.Comment),
		CartItems:        mapCartItems(cart),
	}

	res := a.Handler.Handle(r.Context(), cmd)
	if !res.IsSuccess {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     res.Error,
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, CreateOrderResponse{
		IsSuccess: true,
		OrderID:   res.OrderID,
	}, http.StatusOK)
}

func customerDisplayName(u *authEntity.User) string {
	if u.FullName != nil {
		if n := strings.TrimSpace(*u.FullName); n != "" {
			return n
		}
	}
	return normalizeCustomerPhone(u.Phone)
}

func normalizeCustomerPhone(phone string) string {
	p := strings.TrimSpace(phone)
	if p == "" {
		return p
	}
	if strings.HasPrefix(p, "+") {
		return p
	}
	return "+" + p
}

func mapCartItems(cart *cartModel.Cart) []createOrder.CartLineItem {
	out := make([]createOrder.CartLineItem, 0, len(cart.Items))
	for _, it := range cart.Items {
		out = append(out, createOrder.CartLineItem{
			ProductID: it.ProductID.String(),
			Quantity:  it.Quantity,
			Price:     it.TotalPrice,
			Comment:   it.Comment,
		})
	}
	return out
}

func optionalTrimmedString(s string) *string {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil
	}
	return &t
}

func writeCreateOrderError(w http.ResponseWriter, status int, msg string) {
	utils.JSONResponse(w, CreateOrderResponse{
		IsSuccess: false,
		Error:     msg,
	}, status)
}
