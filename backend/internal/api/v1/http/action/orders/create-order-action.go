package orders

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
)

type CreateOrderAction struct {
	Handler *createOrder.Handler
}

func NewCreateOrderAction(
	handler *createOrder.Handler,
) *CreateOrderAction {
	return &CreateOrderAction{
		Handler: handler,
	}
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

	res := a.Handler.Handle(r.Context(), req.CustomerPhone, &req.Comment)
	if !res.IsSuccess {
		utils.JSONResponse(w, CreateOrderResponse{
			IsSuccess: false,
			Error:     res.Error,
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, CreateOrderResponse{
		IsSuccess: res.IsSuccess,
		OrderID:   res.OrderID,
	}, http.StatusOK)
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

func (a *CreateOrderAction) GetInputType() reflect.Type {
	return reflect.TypeOf(CreateOrderRequest{})
}
