package queue

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	addInQueue "github.com/scrumno/scrumno-api/internal/queue/command/add-in-queue"
)

type AddInQueueAction struct {
	Handler *addInQueue.Handler
}

type AddInQueueRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	QueueID uuid.UUID `json:"queue_id"`
}

type AddInQueueResponse struct {
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
}

func NewAddInQueueAction(handler *addInQueue.Handler) *AddInQueueAction {
	return &AddInQueueAction{
		Handler: handler,
	}
}

func (a *AddInQueueAction) Action(w http.ResponseWriter, r *http.Request) {
	var req AddInQueueRequest
	if err := utils.DecodeJSONBody(r, &req); err != nil {
		utils.JSONResponse(w, AddInQueueResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusBadRequest)
		return
	}

	if req.OrderID == uuid.Nil || req.QueueID == uuid.Nil {
		utils.JSONResponse(w, AddInQueueResponse{
			IsSuccess: false,
			Error:     "Поля order_id и queue_id обязательны",
		}, http.StatusBadRequest)
		return
	}

	if err := a.Handler.Handle(r.Context(), addInQueue.Command{
		OrderID: req.OrderID,
		QueueID: req.QueueID,
	}); err != nil {
		utils.JSONResponse(w, AddInQueueResponse{
			IsSuccess: false,
			Error:     err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, AddInQueueResponse{
		IsSuccess: true,
	}, http.StatusOK)
}

func (a *AddInQueueAction) GetInputType() reflect.Type {
	return reflect.TypeOf(AddInQueueRequest{})
}
