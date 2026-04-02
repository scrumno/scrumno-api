package cart

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"
	clearCart "github.com/scrumno/scrumno-api/internal/cart/command/clear-cart"
)

type ClearAction struct {
	Handler *clearCart.Handler
}

func NewClearAction(handler *clearCart.Handler) *ClearAction {
	return &ClearAction{
		Handler: handler,
	}
}

type ClearRequest struct {
	UserID uuid.UUID `json:"UserID"`
}

func (a *ClearAction) GetInputType() reflect.Type {
	return reflect.TypeOf(ClearRequest{})
}

func (a *ClearAction) Action(w http.ResponseWriter, r *http.Request) {
	var req ClearRequest
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

	cmd := clearCart.Command{
		UserID: req.UserID,
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
