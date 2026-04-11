package get_queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/queue/entity"
)

type Fetcher struct {
	queueRepo entity.QueueRepository
}

func NewFetcher(queueRepo entity.QueueRepository) *Fetcher {
	return &Fetcher{
		queueRepo: queueRepo,
	}
}

type Command struct {
	QueueID        uuid.UUID
	ExcludeOrderID uuid.UUID
}

func (h *Fetcher) Fetch(ctx context.Context, cmd Command) ([]entity.OrdersQueueTable, error) {
	orders, err := h.queueRepo.ListOrdersAhead(ctx, cmd.QueueID, cmd.ExcludeOrderID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
