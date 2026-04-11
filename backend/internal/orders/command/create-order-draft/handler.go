package create_order_draft

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	createOrder "github.com/scrumno/scrumno-api/internal/orders/command/create-order"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
)

type Handler struct {
	orderRepo ordersEntity.OrderRepository
}

type Result struct {
	IsSuccess bool      `json:"is_success"`
	DraftID   uuid.UUID `json:"draft_id,omitempty"`
	Amount    float64   `json:"amount,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func NewHandler(orderRepo ordersEntity.OrderRepository) *Handler {
	return &Handler{orderRepo: orderRepo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) Result {
	if cmd.UserID == uuid.Nil || cmd.VenueID == uuid.Nil {
		return Result{IsSuccess: false, Error: "Пользователь и точка обязательны"}
	}
	if len(cmd.CartItems) == 0 {
		return Result{IsSuccess: false, Error: "Корзина пуста"}
	}

	sourceKey := strings.TrimSpace(cmd.SourceKey)
	if sourceKey == "" {
		sourceKey = "app"
	}

	amount := 0.0
	for _, item := range cmd.CartItems {
		amount += item.Price
	}

	snapshot, err := json.Marshal(map[string]any{
		"source_key": sourceKey,
		"phone":      cmd.CustomerPhone,
		"name":       cmd.CustomerName,
		"comment":    cmd.Comment,
		"items":      cmd.CartItems,
	})
	if err != nil {
		return Result{IsSuccess: false, Error: err.Error()}
	}

	draftID := uuid.New()
	draft := &ordersEntity.OrderDraftTable{
		ID:               draftID,
		UserID:           cmd.UserID,
		VenueID:          cmd.VenueID,
		CartSnapshotJSON: snapshot,
		Amount:           amount,
		PaymentStatus:    false,
	}
	if err := h.orderRepo.CreateDraft(ctx, draft); err != nil {
		return Result{IsSuccess: false, Error: err.Error()}
	}

	return Result{
		IsSuccess: true,
		DraftID:   draftID,
		Amount:    amount,
	}
}

func BuildProviderCommandFromDraft(draft *ordersEntity.OrderDraftTable) (createOrder.Command, error) {
	var payload struct {
		Phone     string                     `json:"phone"`
		Name      string                     `json:"name"`
		SourceKey string                     `json:"source_key"`
		Comment   *string                    `json:"comment"`
		Items     []createOrder.CartLineItem `json:"items"`
	}
	if err := json.Unmarshal(draft.CartSnapshotJSON, &payload); err != nil {
		return createOrder.Command{}, err
	}

	return createOrder.Command{
		CustomerPhone:    payload.Phone,
		CustomerFullName: payload.Name,
		SourceKey:        payload.SourceKey,
		OrderComment:     payload.Comment,
		CartItems:        payload.Items,
	}, nil
}
