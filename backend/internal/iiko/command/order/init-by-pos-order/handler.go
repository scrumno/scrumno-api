package init_by_pos_order

import (
	"context"
	"fmt"
	"strings"

	"github.com/scrumno/scrumno-api/internal/iiko/entity/order"
)

type POSOrdersInitializer interface {
	InitByPosOrder(ctx context.Context, request order.InitByPosOrderRequest) (*order.CorrelationIDResponse, error)
}

type Handler struct {
	repo POSOrdersInitializer
}

func NewHandler(repo POSOrdersInitializer) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, command Command) (*order.CorrelationIDResponse, error) {
	organizationID := strings.TrimSpace(command.OrganizationID)
	if organizationID == "" {
		return nil, fmt.Errorf("не передан organizationId для init_by_posOrder iiko")
	}

	terminalGroupID := strings.TrimSpace(command.TerminalGroupID)
	if terminalGroupID == "" {
		return nil, fmt.Errorf("не передан terminalGroupId для init_by_posOrder iiko")
	}

	posOrderIDs := make([]string, 0, len(command.PosOrderIDs))
	for _, id := range command.PosOrderIDs {
		trimmed := strings.TrimSpace(id)
		if trimmed == "" {
			continue
		}
		posOrderIDs = append(posOrderIDs, trimmed)
	}
	if len(posOrderIDs) == 0 {
		return nil, fmt.Errorf("не переданы posOrderIds для init_by_posOrder iiko")
	}

	return h.repo.InitByPosOrder(ctx, order.InitByPosOrderRequest{
		OrganizationID:  organizationID,
		TerminalGroupID: terminalGroupID,
		PosOrderIDs:     posOrderIDs,
	})
}
