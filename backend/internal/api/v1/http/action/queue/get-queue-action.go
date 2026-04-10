package queue

import (
	"net/http"
	"reflect"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	getWorkingTime "github.com/scrumno/scrumno-api/internal/app/query/get-working-time"
	getQueue "github.com/scrumno/scrumno-api/internal/queue/query/get-queue"
	"github.com/scrumno/scrumno-api/internal/queue/service"
)

type GetQueueAction struct {
	GetWorkingHoursFetcher       *getWorkingTime.Fetcher
	GetCurrentOrdersQueueFetcher *getQueue.Fetcher

	// TODO:getOrderDetailByOrderIDFetcher for get order detail by order id

	QueueCalculator *service.OrdersQueueService
}

func NewGetQueueAction(
	getWorkingHoursFetcher *getWorkingTime.Fetcher,
	getCurrentOrdersQueueFetcher *getQueue.Fetcher,
	queueCalculator *service.OrdersQueueService,

	// TODO:getOrderDetailByOrderIDFetcher for get order detail by order id
) *GetQueueAction {
	return &GetQueueAction{
		GetWorkingHoursFetcher:       getWorkingHoursFetcher,
		GetCurrentOrdersQueueFetcher: getCurrentOrdersQueueFetcher,
		QueueCalculator:              queueCalculator,
	}
}

func (a *GetQueueAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetQueueRequest{})
}

type GetQueueRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}

type GetQueueErrorResponse struct {
	Error string `json:"error"`
}

func (a *GetQueueAction) Action(w http.ResponseWriter, r *http.Request) {
	var req GetQueueRequest
	err := utils.DecodeJSONBody(r, &req)
	if err != nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	workingHours, err := a.GetWorkingHoursFetcher.Fetch(r.Context())
	if err != nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// TODO: get order detail by order id logic here
	if workingHours.OpenAt == "" || workingHours.CloseAt == "" {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: "Working hours are not set",
		}, http.StatusInternalServerError)
		return
	}
}
