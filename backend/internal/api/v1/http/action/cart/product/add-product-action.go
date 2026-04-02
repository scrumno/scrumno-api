package cart

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"
	addProduct "github.com/scrumno/scrumno-api/internal/cart/command/add-product-to-cart"
)

type AddProductAction struct {
	Handler *addProduct.Handler
}

func NewAddProductAction(handler *addProduct.Handler) *AddProductAction {
	return &AddProductAction{
		Handler: handler,
	}
}

type AddProductRequest struct {
	UserID    uuid.UUID `json:"UserID"`
	ProductID uuid.UUID `json:"ProductID"`
	Quantity  float64   `json:"Quantity"`
	BasePrice float64   `json:"BasePrice"`
}

func (a *AddProductAction) GetInputType() reflect.Type {
	return reflect.TypeOf(AddProductRequest{})
}

func (a *AddProductAction) Action(w http.ResponseWriter, r *http.Request) {
	var req AddProductRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if req.UserID == uuid.Nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     "UserID обязателен",
		}, http.StatusBadRequest)
		return
	}

	cmd := addProduct.Command{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		BasePrice: req.BasePrice,
	}

	err := a.Handler.Handle(r.Context(), cmd)
	if err != nil {
		slog.Error("add product to cart", "err", err)
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, model.BaseSuccessResponse{
		IsSuccess: true,
	}, http.StatusOK)
}
