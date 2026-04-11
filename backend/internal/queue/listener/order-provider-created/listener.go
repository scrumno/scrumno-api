package order_provider_created

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
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
	p, ok := payload.(service.OrderProviderCreatedPayload)
	if !ok {
		return
	}

	queueID := p.VenueID
	if queueID == uuid.Nil {
		queueID = p.ProviderOrderID
	}

	if err := l.queueRepo.AddOrderToQueue(context.Background(), p.ProviderOrderID, queueID); err != nil {
		slog.Error("queue.add.order", "error", err)
	}
}
