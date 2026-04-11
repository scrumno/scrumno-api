package pay_order_draft

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	iikoModel "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/model"
	iikoService "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/order-delivery/service"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	createOrderDraft "github.com/scrumno/scrumno-api/internal/orders/command/create-order-draft"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

type CommandStatusChecker interface {
	GetStatus(ctx context.Context, correlationID uuid.UUID) (*iikoService.CommandStatusResponse, error)
}

type Handler struct {
	orderRepo      ordersEntity.OrderRepository
	paymentService ordersService.PaymentService
	provider       createOrder.Provider
	statusChecker  CommandStatusChecker
	eventManager   *eventManager.EventManager
}

type Result struct {
	IsSuccess     bool      `json:"is_success"`
	DraftID       uuid.UUID `json:"draft_id,omitempty"`
	OrderID       uuid.UUID `json:"order_id,omitempty"`
	CanSubscribe  bool      `json:"can_subscribe"`
	ProviderState string    `json:"provider_state,omitempty"`
	Error         string    `json:"error,omitempty"`
}

func NewHandler(
	orderRepo ordersEntity.OrderRepository,
	paymentService ordersService.PaymentService,
	provider createOrder.Provider,
	statusChecker CommandStatusChecker,
	eventManager *eventManager.EventManager,
) *Handler {
	return &Handler{
		orderRepo:      orderRepo,
		paymentService: paymentService,
		provider:       provider,
		statusChecker:  statusChecker,
		eventManager:   eventManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) Result {
	draft, err := h.orderRepo.GetDraftByID(ctx, cmd.DraftID)
	if err != nil {
		return Result{IsSuccess: false, Error: "Черновик заказа не найден"}
	}

	ok, err := h.paymentService.Pay(ctx, draft.ID.String(), draft.Amount)
	if err != nil || !ok {
		return Result{IsSuccess: false, Error: "Ошибка оплаты заказа"}
	}
	if err := h.orderRepo.MarkDraftPaymentSuccess(ctx, draft.ID); err != nil {
		return Result{IsSuccess: false, Error: err.Error()}
	}
	if h.eventManager != nil {
		h.eventManager.EmitEvent("order.payment.succeeded", draft)
	}

	providerCmd, err := createOrderDraft.BuildProviderCommandFromDraft(draft)
	if err != nil {
		return Result{IsSuccess: false, Error: err.Error()}
	}
	providerResult := h.provider.Handle(ctx, providerCmd)
	if !providerResult.IsSuccess {
		_ = h.orderRepo.MarkDraftProviderFailed(ctx, draft.ID, providerResult.Error)
		return Result{IsSuccess: false, Error: providerResult.Error}
	}

	correlationID := extractCorrelationID(providerResult.Response)
	orderID, hasOrder := parseOrderID(providerResult.OrderID)
	switch strings.ToLower(strings.TrimSpace(providerResult.CreationStatus)) {
	case "error":
		_ = h.orderRepo.MarkDraftProviderFailed(ctx, draft.ID, "Ошибка создания заказа в iiko")
		return Result{IsSuccess: false, Error: "Ошибка создания заказа в iiko"}
	case "inprogress":
		// продолжим в блок pending ниже
	default:
		if hasOrder {
			if err := h.orderRepo.MarkDraftProviderSuccess(ctx, draft.ID, orderID); err != nil {
				return Result{IsSuccess: false, Error: err.Error()}
			}
			h.emitProviderCreatedEvent(*draft, orderID)
			return Result{
				IsSuccess:     true,
				DraftID:       draft.ID,
				OrderID:       orderID,
				CanSubscribe:  true,
				ProviderState: "Success",
			}
		}
	}

	if hasOrder && correlationID == uuid.Nil {
		if err := h.orderRepo.MarkDraftProviderSuccess(ctx, draft.ID, orderID); err != nil {
			return Result{IsSuccess: false, Error: err.Error()}
		}
		h.emitProviderCreatedEvent(*draft, orderID)
		return Result{
			IsSuccess:     true,
			DraftID:       draft.ID,
			OrderID:       orderID,
			CanSubscribe:  true,
			ProviderState: "Success",
		}
	}

	if correlationID == uuid.Nil {
		_ = h.orderRepo.MarkDraftProviderFailed(ctx, draft.ID, "Не получен correlationId для проверки статуса")
		return Result{IsSuccess: false, Error: "Не удалось определить статус создания заказа в iiko"}
	}

	if err := h.orderRepo.MarkDraftProviderPending(ctx, draft.ID, correlationID); err != nil {
		return Result{IsSuccess: false, Error: err.Error()}
	}

	if h.statusChecker == nil {
		return Result{
			IsSuccess:     true,
			DraftID:       draft.ID,
			CanSubscribe:  false,
			ProviderState: "InProgress",
		}
	}

	checkCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	for {
		select {
		case <-checkCtx.Done():
			return Result{
				IsSuccess:     true,
				DraftID:       draft.ID,
				CanSubscribe:  false,
				ProviderState: "InProgress",
			}
		default:
			status, err := h.statusChecker.GetStatus(checkCtx, correlationID)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			switch strings.ToLower(status.State) {
			case "success":
				_ = h.orderRepo.MarkDraftProviderSuccess(ctx, draft.ID, draft.ID)
				h.emitProviderCreatedEvent(*draft, draft.ID)
				return Result{
					IsSuccess:     true,
					DraftID:       draft.ID,
					OrderID:       draft.ID,
					CanSubscribe:  true,
					ProviderState: "Success",
				}
			case "error":
				reason := "Ошибка создания заказа в iiko"
				if status.ErrorReason != nil {
					reason = *status.ErrorReason
				}
				_ = h.orderRepo.MarkDraftProviderFailed(ctx, draft.ID, reason)
				return Result{IsSuccess: false, Error: reason}
			default:
				time.Sleep(2 * time.Second)
			}
		}
	}
}

func (h *Handler) emitProviderCreatedEvent(draft ordersEntity.OrderDraftTable, orderID uuid.UUID) {
	if h.eventManager == nil {
		return
	}
	h.eventManager.EmitEvent("order.provider.created", ordersService.OrderProviderCreatedPayload{
		DraftID:         draft.ID,
		UserID:          draft.UserID,
		VenueID:         draft.VenueID,
		ProviderOrderID: orderID,
	})
}

func extractCorrelationID(response any) uuid.UUID {
	orderResponse, ok := response.(*iikoModel.OrderSetResponse)
	if !ok || orderResponse == nil {
		return uuid.Nil
	}
	if orderResponse.CorrelationID != uuid.Nil {
		return orderResponse.CorrelationID
	}
	if orderResponse.OrderInfo != nil && orderResponse.OrderInfo.ID != uuid.Nil {
		return orderResponse.OrderInfo.ID
	}
	return uuid.Nil
}

func parseOrderID(raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil || id == uuid.Nil {
		return uuid.Nil, false
	}
	return id, true
}
