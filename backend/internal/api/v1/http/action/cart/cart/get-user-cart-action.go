package cart

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/api/v1/http/action/cart/model"

	getCart "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
)

type GetCartAction struct {
	Fetcher *getCart.Fetcher
}

func NewGetCartAction(fetcher *getCart.Fetcher) *GetCartAction {
	return &GetCartAction{
		Fetcher: fetcher,
	}
}

type GetCartRequest struct {
	UserID uuid.UUID `json:"UserID"`
}

func (a *GetCartAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetCartRequest{})
}

func (a *GetCartAction) Action(w http.ResponseWriter, r *http.Request) {
	var req GetCartRequest
	raw := r.URL.Query().Get("userId")

	if raw == "" {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     "userId обязателен",
		}, http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(raw)
	if err != nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     "некорректный userId",
		}, http.StatusBadRequest)
		return
	}
	req.UserID = uid

	if req.UserID == uuid.Nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     "ID пользователя обязательно должно быть передано",
		}, http.StatusBadRequest)
		return
	}

	cmd := getCart.Query{
		UserID: req.UserID,
	}

	cart, err := a.Fetcher.Fetch(r.Context(), cmd)
	if err != nil {
		utils.JSONResponse(w, model.BaseErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, model.CartSuccessResponse{
		IsSuccess: true,
		Cart:      cart,
	}, http.StatusOK)
}
