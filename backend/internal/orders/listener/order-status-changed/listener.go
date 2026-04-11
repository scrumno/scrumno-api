package order_status_changed

import (
	"context"
	"log/slog"
	"strings"

	ordersEntity "github.com/scrumno/scrumno-api/internal/orders/entity"
	ordersService "github.com/scrumno/scrumno-api/internal/orders/service"
)

type Listener struct {
	orderRepo ordersEntity.OrderRepository
	hub       ordersService.OrdersWebSocketHub
}

func NewListener(orderRepo ordersEntity.OrderRepository, hub ordersService.OrdersWebSocketHub) *Listener {
	return &Listener{
		orderRepo: orderRepo,
		hub:       hub,
	}
}

func (l *Listener) Listen(payload any) {
	p, ok := payload.(ordersService.OrderStatusChangedPayload)
	if !ok {
		return
	}

	if err := l.orderRepo.UpdateHistoryStatus(context.Background(), p.ProviderOrderID, p.Status); err != nil {
		slog.Error("order.history.update.status", "error", err)
	}

	subscribers, err := l.orderRepo.ListActiveSubscribersByOrder(context.Background(), p.ProviderOrderID)
	if err != nil {
		slog.Error("order.subscribers.list", "error", err)
		return
	}

	connectionIDs := make([]string, 0, len(subscribers))
	for _, subscriber := range subscribers {
		connectionIDs = append(connectionIDs, subscriber.ConnectionID)
	}
	l.hub.Notify(connectionIDs, p.ProviderOrderID, p.Status)

	status := strings.ToLower(strings.TrimSpace(p.Status))
	if status == "closed" || status == "cancelled" {
		for _, subscriber := range subscribers {
			orderID := p.ProviderOrderID
			if err := l.orderRepo.DeactivateSubscriber(context.Background(), subscriber.ConnectionID, &orderID); err != nil {
				slog.Error("order.subscriber.deactivate", "error", err)
			}
		}
	}
}
