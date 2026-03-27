package menu

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/interfaces"
	"github.com/scrumno/scrumno-api/internal/api/utils"
)

type RefreshMenuAction struct {
	Handler interfaces.GetMenuHandler
}

func NewRefreshMenuAction(handler interfaces.GetMenuHandler) RefreshMenuAction {
	return RefreshMenuAction{Handler: handler}
}

type RefreshMenuRequest struct{}

type RefreshMenuResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Error     string `json:"error,omitempty"`
	Menu      any    `json:"menu,omitempty"`
}

func (a *RefreshMenuAction) Action(w http.ResponseWriter, r *http.Request) {
	if a == nil || a.Handler == nil {
		utils.JSONResponse(w, RefreshMenuResponse{
			IsSuccess: false,
			Error:     "хэндлер не настроен",
		}, http.StatusInternalServerError)
		return
	}

	result := a.Handler.Handle()
	if result == nil {
		utils.JSONResponse(w, RefreshMenuResponse{
			IsSuccess: false,
			Error:     "пустой результат",
		}, http.StatusInternalServerError)
		return
	}

	if err, ok := result.(error); ok {
		utils.JSONResponse(w, RefreshMenuResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadGateway)
		return
	}

	utils.JSONResponse(w, RefreshMenuResponse{
		IsSuccess: true,
		Menu:      result,
	}, http.StatusOK)
}

func (a *RefreshMenuAction) GetInputType() reflect.Type {
	return reflect.TypeOf(RefreshMenuRequest{})
}
