package get_queue

import (
	"context"
	"time"

	getWorkingTime "github.com/scrumno/scrumno-api/internal/app/query/get-working-time"
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
	WorkingHours     getWorkingTime.WorkingHours
	OrdersQueueOrder []entity.OrdersQueueOrder
}

func (h *Fetcher) Fetch(ctx context.Context, cmd Command) (time.Duration, error) {
	queueTime, err := h.queueRepo.GetQueueTime(ctx, cmd.WorkingHours, cmd.OrdersQueueOrder)
	if err != nil {
		return 0, err
	}
	return queueTime, nil
}
