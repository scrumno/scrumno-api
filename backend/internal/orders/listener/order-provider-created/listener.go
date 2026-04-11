package order_provider_created

import (
	"context"
	"log/slog"

	cartEntity "github.com/scrumno/scrumno-api/internal/cart/entity"
	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
)

type Listener struct {
	orderRepo ordersEntity.OrderRepository
	cartRepo  cartEntity.CartRepository
}

func NewListener(orderRepo ordersEntity.OrderRepository, cartRepo cartEntity.CartRepository) *Listener {
	return &Listener{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (l *Listener) Listen(payload any) {
	p, ok := payload.(ordersService.OrderProviderCreatedPayload)
	if !ok {
		return
	}

	if err := l.orderRepo.SaveHistory(context.Background(), &ordersEntity.OrderHistoryTable{
		DraftID:         &p.DraftID,
		UserID:          p.UserID,
		VenueID:         p.VenueID,
		ProviderOrderID: p.ProviderOrderID,
		Status:          "Created",
	}); err != nil {
		slog.Error("order.history.save", "error", err)
	}

	if err := l.cartRepo.ClearCart(context.Background(), p.UserID); err != nil {
		slog.Error("cart.clear.after.order", "error", err)
	}

	if err := l.orderRepo.DeleteDraft(context.Background(), p.DraftID); err != nil {
		slog.Error("order.draft.delete", "error", err)
	}
}
