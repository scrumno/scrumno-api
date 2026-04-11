package queue

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/api/utils"
	getWorkingTime "github.com/scrumno/scrumno-api/internal/app/query/get-working-time"
	getCurrentCart "github.com/scrumno/scrumno-api/internal/cart/query/get-cart-by-user-id"
	queueEntity "github.com/scrumno/scrumno-api/internal/queue/entity"
	getQueue "github.com/scrumno/scrumno-api/internal/queue/query/get-queue"
	"github.com/scrumno/scrumno-api/internal/queue/service"
)

type GetQueueAction struct {
	GetWorkingHoursFetcher        *getWorkingTime.Fetcher
	GetCurrentOrdersQueueFetcher  *getQueue.Fetcher
	GetCurrentCartByUserIdFetcher *getCurrentCart.Fetcher
	QueueCalculator               service.OrdersQueueService
	QueueMapper                   service.QueueOrderMapper
	QueueSync                     service.QueueSyncService
}

func NewGetQueueAction(
	getWorkingHoursFetcher *getWorkingTime.Fetcher,
	getCurrentOrdersQueueFetcher *getQueue.Fetcher,
	queueCalculator service.OrdersQueueService,
	queueMapper service.QueueOrderMapper,
	queueSync service.QueueSyncService,
	getCurrentCartByUserIdFetcher *getCurrentCart.Fetcher,
) *GetQueueAction {
	return &GetQueueAction{
		GetWorkingHoursFetcher:        getWorkingHoursFetcher,
		GetCurrentOrdersQueueFetcher:  getCurrentOrdersQueueFetcher,
		QueueCalculator:               queueCalculator,
		QueueMapper:                   queueMapper,
		QueueSync:                     queueSync,
		GetCurrentCartByUserIdFetcher: getCurrentCartByUserIdFetcher,
	}
}

type GetQueueRequest struct {
	OrderID uuid.UUID `json:"order_id"`
	UserID  uuid.UUID `json:"user_id"`
	QueueID uuid.UUID `json:"queue_id"`
	VenueID uuid.UUID `json:"venue_id"`
}

type GetQueueErrorResponse struct {
	Error string `json:"error"`
}

type GetQueueResponse struct {
	CurrentOrderNoQueue service.DurationRange `json:"current_order_no_queue"`
	QueueWait           service.DurationRange `json:"queue_wait"`
	Total               service.DurationRange `json:"total"`
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
	if req.UserID == uuid.Nil || req.OrderID == uuid.Nil || req.QueueID == uuid.Nil || req.VenueID == uuid.Nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: "Поля order_id, user_id, queue_id и venue_id обязательны",
		}, http.StatusBadRequest)
		return
	}

	result, err := a.estimateQueue(r.Context(), req)
	if err != nil {
		utils.JSONResponse(w, GetQueueErrorResponse{
			Error: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, GetQueueResponse{
		CurrentOrderNoQueue: result.CurrentOrderNoQueue,
		QueueWait:           result.QueueWait,
		Total:               result.Total,
	}, http.StatusOK)
}

func (a *GetQueueAction) estimateQueue(ctx context.Context, req GetQueueRequest) (service.OrdersQueueResult, error) {
	ctx = context.WithValue(ctx, "venue_id", req.VenueID.String())
	if a.QueueSync != nil {
		if err := a.QueueSync.RefreshQueue(ctx, req.VenueID, req.QueueID); err != nil {
			return service.OrdersQueueResult{}, err
		}
	}

	workingHours, err := a.GetWorkingHoursFetcher.Fetch(ctx)
	if err != nil {
		return service.OrdersQueueResult{}, err
	}

	if workingHours.OpenAt == "" || workingHours.CloseAt == "" {
		return service.OrdersQueueResult{}, errors.New("рабочие часы не установлены")
	}

	currentCart, err := a.GetCurrentCartByUserIdFetcher.Fetch(ctx, getCurrentCart.Query{
		UserID: req.UserID,
	})
	if err != nil {
		return service.OrdersQueueResult{}, err
	}

	currentOrder, err := a.QueueMapper.MapCartToOrder(ctx, req.QueueID, req.OrderID, currentCart)
	if err != nil {
		return service.OrdersQueueResult{}, err
	}

	queueRows, err := a.GetCurrentOrdersQueueFetcher.Fetch(ctx, getQueue.Command{
		QueueID:        req.QueueID,
		ExcludeOrderID: req.OrderID,
	})
	if err != nil {
		return service.OrdersQueueResult{}, err
	}
	ordersAhead := buildOrdersAhead(queueRows)

	calculator := a.QueueCalculator
	if calculator == nil {
		cfg := &queueEntity.OrdersQueueConfigTable{
			KitchenParallelSlots:  1,
			QueueGrowthFactor:     0.15,
			OrderReserveMinutes:   2,
			RestaurantOpenAt:      workingHours.OpenAt,
			RestaurantCloseAt:     workingHours.CloseAt,
			EmptyQueueWaitMinMins: 10,
			EmptyQueueWaitMaxMins: 10,
			QueueTimeMinFactor:    0.90,
			QueueTimeMaxFactor:    1.25,
		}
		calculator = service.NewOrdersQueueService(cfg, nil)
	}

	return calculator.EstimateQueueTime(currentOrder, ordersAhead), nil
}

func (a *GetQueueAction) GetInputType() reflect.Type {
	return reflect.TypeOf(GetQueueRequest{})
}

func buildOrdersAhead(rows []queueEntity.OrdersQueueTable) []queueEntity.OrdersQueueOrder {
	result := make([]queueEntity.OrdersQueueOrder, 0, len(rows))
	now := time.Now()
	for _, row := range rows {
		estimated := row.EstimatedCookMins
		if estimated <= 0 {
			estimated = 5
		}

		result = append(result, queueEntity.OrdersQueueOrder{
			ID:           row.OrderID,
			ExternalID:   row.OrderID.String(),
			QueueID:      row.QueueID,
			SetupMinutes: 1,
			CreatedAt:    now,
			UpdatedAt:    now,
			Items: []queueEntity.OrderItem{
				{
					ID:               uuid.New(),
					OrderID:          row.OrderID,
					ProductID:        "offline",
					Quantity:         1,
					BaseCookMinutes:  estimated,
					GrowthFactor:     0,
					ComplexityFactor: 1,
				},
			},
		})
	}
	return result
}

