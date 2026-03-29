package cart

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"
	update "github.com/scrumno/scrumno-api/internal/cart/command/update-product-cart"
)

type UpdateAction struct {
	Handler *update.Handler
}

func NewUpdateAction(handler *update.Handler) *UpdateAction {
	return &UpdateAction{
		Handler: handler,
	}
}

type UpdateRequest struct {
	UserID    uuid.UUID `json:"UserID"`
	ProductID uuid.UUID `json:"ProductID"`
	Quantity  float64   `json:"Quantity"`
}

func (a *UpdateAction) GetInputType() reflect.Type {
	return reflect.TypeOf(UpdateRequest{})
}

func (a *UpdateAction) Action(w http.ResponseWriter, r *http.Request) {
	var req UpdateRequest
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
			Error:     "ID пользователя обязательно должно быть передано",
		}, http.StatusBadRequest)
		return
	}

	cmd := update.Command{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	err := a.Handler.Handle(r.Context(), cmd)
	if err != nil {
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
