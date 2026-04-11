package order_status_changed

import (
	"context"
	"log/slog"
	"strings"

	"github.com/scrumno/scrumno-api/internal/orders/service"
	queueEntity "github.com/scrumno/scrumno-api/internal/queue/entity"
)

type Listener struct {
	queueRepo queueEntity.QueueRepository
}

func NewListener(queueRepo queueEntity.QueueRepository) *Listener {
	return &Listener{queueRepo: queueRepo}
}

func (l *Listener) Listen(payload any) {
	p, ok := payload.(service.OrderStatusChangedPayload)
	if !ok {
		return
	}

	status := strings.ToLower(strings.TrimSpace(p.Status))
	if status != "closed" && status != "cancelled" {
		return
	}

	if err := l.queueRepo.RemoveOrderFromQueue(context.Background(), p.ProviderOrderID); err != nil {
		slog.Error("queue.remove.order", "error", err)
	}
}
