package auth

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/authorize/command/logout"
	findUserByPhoneFetcher "github.com/scrumno/scrumno-api/internal/authorize/query/find-user-by-phone"
)

type LogoutAction struct {
	Handler *logout.Handler
	FindUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher
}

func NewLogoutAction(
	handler *logout.Handler, 
	findUserByPhoneFetcher *findUserByPhoneFetcher.Fetcher,
) *LogoutAction {
	return &LogoutAction{
		Handler: handler,
		FindUserByPhoneFetcher: findUserByPhoneFetcher,
	}
}

func (a *LogoutAction) GetInputType() reflect.Type {
    return reflect.TypeOf(LogoutRequest{})
}

type LogoutRequest struct {
	Phone string `json:"phone" example:"79099000000"`
}

func (a *LogoutAction) Action(w http.ResponseWriter, r *http.Request) {

	var req LogoutRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, LogoutErrorResponse{
			IsSuccess: false, 
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	user, err := a.FindUserByPhoneFetcher.Fetch(r.Context(), req.Phone)
	if err != nil {
		utils.JSONResponse(w, LogoutErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}
	if user == nil {
		utils.JSONResponse(w, LogoutResponse{IsSuccess: true}, http.StatusOK)
		return
	}

	err = a.Handler.Handle(r.Context(), user.ID)
	if err != nil {
		utils.JSONResponse(w, LogoutErrorResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	utils.JSONResponse(w, LogoutResponse{
		IsSuccess: true,
	}, http.StatusOK)
}

type LogoutResponse struct {
    IsSuccess    bool   `json:"isSuccess"`
}

type LogoutErrorResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error"`
}