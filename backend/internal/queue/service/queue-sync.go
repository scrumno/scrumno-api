package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	appConfig "github.com/scrumno/scrumno-api/internal/app/entity/app-config"
	"github.com/scrumno/scrumno-api/internal/queue/entity"
)

type OrdersByRevisionSyncProvider interface {
	GetOrdersByRevision(ctx context.Context, startRevision int64, sourceKeys []string) (*OrdersByRevisionPayload, error)
}

type OrdersByRevisionPayload struct {
	MaxRevision int64
	Orders      []OrdersByRevisionOrder
}

type OrdersByRevisionOrder struct {
	OrderID           uuid.UUID
	Status            string
	IsDeleted         bool
	EstimatedCookMins int
	CompleteBeforeAt  *time.Time
}

type QueueSyncService interface {
	RefreshQueue(ctx context.Context, venueID uuid.UUID, queueID uuid.UUID) error
}

type queueSyncService struct {
	appConfigRepo appConfig.AppConfigRepository
	queueRepo     entity.QueueRepository
	provider      OrdersByRevisionSyncProvider
}

func NewQueueSyncService(
	appConfigRepo appConfig.AppConfigRepository,
	queueRepo entity.QueueRepository,
	provider OrdersByRevisionSyncProvider,
) QueueSyncService {
	return &queueSyncService{
		appConfigRepo: appConfigRepo,
		queueRepo:     queueRepo,
		provider:      provider,
	}
}

func (s *queueSyncService) RefreshQueue(ctx context.Context, venueID uuid.UUID, queueID uuid.UUID) error {
	revision, err := s.appConfigRepo.GetQueueSyncState(ctx, venueID)
	if err != nil {
		revision = 0
	}

	payload, err := s.provider.GetOrdersByRevision(ctx, revision, []string{"app"})
	if err != nil {
		return err
	}

	for _, order := range payload.Orders {
		if order.IsDeleted || !isQueueActiveStatus(order.Status) {
			if err := s.queueRepo.RemoveOrderFromQueue(ctx, order.OrderID); err != nil {
				return err
			}
			continue
		}

		if err := s.queueRepo.UpsertQueueOrder(ctx, order.OrderID, queueID, order.Status, order.EstimatedCookMins, order.CompleteBeforeAt); err != nil {
			return err
		}
	}

	if payload.MaxRevision > revision {
		if err := s.appConfigRepo.UpdateQueueSyncState(ctx, venueID, payload.MaxRevision); err != nil {
			return err
		}
	}

	return nil
}

func isQueueActiveStatus(status string) bool {
	switch strings.ToLower(status) {
	case "unconfirmed", "waitcooking", "readyforcooking", "cookingstarted", "cookingcompleted", "waiting", "onway", "delivered":
		return true
	default:
		return false
	}
}
