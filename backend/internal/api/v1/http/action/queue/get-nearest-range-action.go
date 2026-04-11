package queue

import (
	"net/http"
	"reflect"

	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/queue/service"
)

type GetNearestRangeAction struct {
	GetQueueAction *GetQueueAction
}

func NewGetNearestRangeAction(getQueueAction *GetQueueAction) *GetNearestRangeAction {
	return &GetNearestRangeAction{
		GetQueueAction: getQueueAction,
	}
}

func (a *GetNearestRangeAction) Action(w http.ResponseWriter, r *http.Request) {
	var req GetQueueRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	result, err := a.GetQueueAction.estimateQueue(r.Context(), req)
	if err != nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, struct {
		Range service.DurationRange `json:"range"`
	}{
		Range: result.Total,
	}, http.StatusOK)
}

func (a *GetNearestRangeAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetQueueRequest{})
}
