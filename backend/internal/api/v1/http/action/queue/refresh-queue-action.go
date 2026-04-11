package queue

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	"github.com/scrumno/scrumno-api/internal/queue/service"
)

type RefreshQueueAction struct {
	QueueSync service.QueueSyncService
}

type RefreshQueueRequest struct {
	VenueID uuid.UUID `json:"venue_id"`
	QueueID uuid.UUID `json:"queue_id"`
}

type RefreshQueueResponse struct {
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
}

func NewRefreshQueueAction(queueSync service.QueueSyncService) *RefreshQueueAction {
	return &RefreshQueueAction{
		QueueSync: queueSync,
	}
}

func (a *RefreshQueueAction) Action(w http.ResponseWriter, r *http.Request) {
	if a.QueueSync == nil {
		utils.JSONResponse(w, RefreshQueueResponse{
			IsSuccess: false,
			Error:     "Сервис обновления очереди не настроен",
		}, http.StatusInternalServerError)
		return
	}

	var req RefreshQueueRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, RefreshQueueResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if req.VenueID == uuid.Nil || req.QueueID == uuid.Nil {
		utils.JSONResponse(w, RefreshQueueResponse{
			IsSuccess: false,
			Error:     "Поля venue_id и queue_id обязательны",
		}, http.StatusBadRequest)
		return
	}

	if err := a.QueueSync.RefreshQueue(r.Context(), req.VenueID, req.QueueID); err != nil {
		utils.JSONResponse(w, RefreshQueueResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, RefreshQueueResponse{
		IsSuccess: true,
	}, http.StatusOK)
}

func (a *RefreshQueueAction) GetInputType() reflect.Type {
	return reflect.TypeOf(RefreshQueueRequest{})
}
