package entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	getWorkingTime "github.com/scrumno/scrumno-api/internal/app/query/get-working-time"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"gorm.io/gorm"
)

type QueueRepository interface {
	GetQueueTime(ctx context.Context, workingHours getWorkingTime.WorkingHours, ordersQueueOrder []OrdersQueueOrder) (time.Duration, error)
}

type queueRepository struct {
	*factory.GormRepository[OrdersQueueTable]
}

func NewQueueRepository(db *gorm.DB) QueueRepository {
	return &queueRepository{
		GormRepository: factory.NewGormRepository[OrdersQueueTable](db),
	}
}

func (r *queueRepository) GetQueueTime(ctx context.Context, workingHours getWorkingTime.WorkingHours, ordersQueueOrder []OrdersQueueOrder) (time.Duration, error) {
	var queueTime time.Duration
	err := r.DB.WithContext(ctx).Model(&OrdersQueueTable{}).
		Where("working_hours = ?", workingHours).
		Where("orders_queue_order_id = ?", ordersQueueOrder[0].ID).
		First(&queueTime).Error
	return queueTime, err
}

func (r *queueRepository) GetQueueTimeOrder(ctx context.Context, orderID uuid.UUID) (time.Duration, error) {
	var queueTime time.Duration
	err := r.DB.WithContext(ctx).Model(&OrdersQueueTable{}).
		Where("order_id = ?", orderID).
		First(&queueTime).Error
	return queueTime, err
}

func (r *queueRepository) RemoveQueueTimeOrder(ctx context.Context, orderID uuid.UUID) error {
	err := r.DB.WithContext(ctx).Model(&OrdersQueueTable{}).
		Where("order_id = ?", orderID).
		Delete(&OrdersQueueTable{}).Error
	return err
}
