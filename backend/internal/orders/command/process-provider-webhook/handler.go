package process_provider_webhook

import (
	"context"
	"strings"

	"github.com/google/uuid"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
	eventManager "github.com/scrumno/scrumno-api/shared/services/event-manager"
)

type Handler struct {
	orderRepo    ordersEntity.OrderRepository
	eventManager *eventManager.EventManager
}

func NewHandler(orderRepo ordersEntity.OrderRepository, eventManager *eventManager.EventManager) *Handler {
	return &Handler{
		orderRepo:    orderRepo,
		eventManager: eventManager,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	if cmd.CorrelationID != nil && *cmd.CorrelationID != uuid.Nil {
		draft, err := h.orderRepo.GetDraftByCorrelationID(ctx, *cmd.CorrelationID)
		if err == nil && draft != nil {
			switch strings.ToLower(cmd.CreationStatus) {
			case "success":
				providerOrderID := draft.ID
				if cmd.ProviderOrderID != nil && *cmd.ProviderOrderID != uuid.Nil {
					providerOrderID = *cmd.ProviderOrderID
				}
				if err := h.orderRepo.MarkDraftProviderSuccess(ctx, draft.ID, providerOrderID); err != nil {
					return err
				}
				if h.eventManager != nil {
					h.eventManager.EmitEvent("order.provider.created", ordersService.OrderProviderCreatedPayload{
						DraftID:         draft.ID,
						UserID:          draft.UserID,
						VenueID:         draft.VenueID,
						ProviderOrderID: providerOrderID,
					})
				}
			case "error":
				reason := "Ошибка сохранения заказа в iiko"
				if cmd.Error != nil && *cmd.Error != "" {
					reason = *cmd.Error
				}
				if err := h.orderRepo.MarkDraftProviderFailed(ctx, draft.ID, reason); err != nil {
					return err
				}
			}
		}
	}

	if cmd.ProviderOrderID != nil && *cmd.ProviderOrderID != uuid.Nil && cmd.Status != "" {
		if h.eventManager != nil {
			h.eventManager.EmitEvent("order.status.changed", ordersService.OrderStatusChangedPayload{
				ProviderOrderID: *cmd.ProviderOrderID,
				Status:          cmd.Status,
			})
		}
	}

	return nil
}
