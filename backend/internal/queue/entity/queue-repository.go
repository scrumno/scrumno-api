package entity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"gorm.io/gorm"
)

type QueueRepository interface {
	AddOrderToQueue(ctx context.Context, orderID uuid.UUID, queueID uuid.UUID) error
	UpsertQueueOrder(ctx context.Context, orderID uuid.UUID, queueID uuid.UUID, status string, estimatedCookMins int, completeBeforeAt *time.Time) error
	ExistsInQueue(ctx context.Context, orderID uuid.UUID) (bool, error)
	ListOrdersAhead(ctx context.Context, queueID uuid.UUID, excludeOrderID uuid.UUID) ([]OrdersQueueTable, error)
	RemoveOrderFromQueue(ctx context.Context, orderID uuid.UUID) error
}

type queueRepository struct {
	*factory.GormRepository[OrdersQueueTable]
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{
		GormRepository: factory.NewGormRepository[OrdersQueueTable](db),
	}
}

func (r *queueRepository) AddOrderToQueue(ctx context.Context, orderID uuid.UUID, queueID uuid.UUID) error {
	return r.UpsertQueueOrder(ctx, orderID, queueID, "", 0, nil)
}

func (r *queueRepository) UpsertQueueOrder(ctx context.Context, orderID uuid.UUID, queueID uuid.UUID, status string, estimatedCookMins int, completeBeforeAt *time.Time) error {
	var existing OrdersQueueTable
	err := r.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return r.DB.WithContext(ctx).Create(&OrdersQueueTable{
			OrderID:           orderID,
			QueueID:           queueID,
			Status:            status,
			EstimatedCookMins: estimatedCookMins,
			CompleteBeforeAt:  completeBeforeAt,
		}).Error
	}

	return r.DB.WithContext(ctx).
		Model(&OrdersQueueTable{}).
		Where("order_id = ?", orderID).
		Updates(map[string]any{
			"queue_id":            queueID,
			"status":              status,
			"estimated_cook_mins": estimatedCookMins,
			"complete_before_at":  completeBeforeAt,
		}).Error
}

func (r *queueRepository) ExistsInQueue(ctx context.Context, orderID uuid.UUID) (bool, error) {
	var entity OrdersQueueTable
	err := r.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *queueRepository) ListOrdersAhead(ctx context.Context, queueID uuid.UUID, excludeOrderID uuid.UUID) ([]OrdersQueueTable, error) {
	query := r.DB.WithContext(ctx).
		Where("queue_id = ?", queueID).
		Order("created_at ASC")
	if excludeOrderID != uuid.Nil {
		query = query.Where("order_id <> ?", excludeOrderID)
	}

	var entities []OrdersQueueTable
	err := query.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *queueRepository) RemoveOrderFromQueue(ctx context.Context, orderID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		Delete(&OrdersQueueTable{}).Error
}
