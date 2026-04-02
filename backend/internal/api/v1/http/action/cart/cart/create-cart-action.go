package cart

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"
	createCart "github.com/scrumno/scrumno-api/internal/cart/command/create-cart"
)

type CreateAction struct {
	Handler *createCart.Handler
}

func NewCreateAction(handler *createCart.Handler) *CreateAction {
	return &CreateAction{
		Handler: handler,
	}
}

type CreateRequest struct {
	UserID uuid.UUID `json:"UserID"`
}

func (a *CreateAction) GetInputType() reflect.Type {
	return reflect.TypeOf(CreateRequest{})
}

func (a *CreateAction) Action(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
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
			Error:     "user_id is required",
		}, http.StatusBadRequest)
		return
	}

	cmd := createCart.Command{
		UserID: req.UserID,
	}

	cart, err := a.Handler.Handle(r.Context(), cmd)
	if err != nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, model.SuccessResponse{
		IsSuccess: true,
		CartID:    cart.ID,
	}, http.StatusOK)
}
