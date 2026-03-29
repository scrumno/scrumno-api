package cart

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"
	removeProduct "github.com/scrumno/scrumno-api/internal/cart/command/remove-product"
)

type RemoveProductAction struct {
	Handler *removeProduct.Handler
}

func NewRemoveProductAction(handler *removeProduct.Handler) *RemoveProductAction {
	return &RemoveProductAction{
		Handler: handler,
	}
}

type RemoveProductRequest struct {
	UserID    uuid.UUID `json:"UserID"`
	ProductID uuid.UUID `json:"ProductID"`
}

func (a *RemoveProductAction) GetInputType() reflect.Type {
	return reflect.TypeOf(RemoveProductRequest{})
}

func (a *RemoveProductAction) Action(w http.ResponseWriter, r *http.Request) {
	var req RemoveProductRequest
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

	cmd := removeProduct.Command{
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}

	del, err := a.Handler.Handle(r.Context(), cmd)
	if err != nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if !del {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     "не удалось удалить товар из корзины",
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, model.BaseSuccessResponse{
		IsSuccess: true,
	}, http.StatusOK)
}
