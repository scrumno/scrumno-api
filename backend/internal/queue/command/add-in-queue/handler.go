package add_in_queue

import (
	"context"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/internal/queue/entity"
)

type Handler struct {
	queueRepo entity.QueueRepository
}

func NewHandler(queueRepo entity.QueueRepository) *Handler {
	return &Handler{
		queueRepo: queueRepo,
	}
}

type Command struct {
	OrderID uuid.UUID
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	return nil
}
