package entity

import (
	"context"
	"errors"

	"github.com/google/uuid"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateDraft(ctx context.Context, draft *OrderDraftTable) error
	GetDraftByID(ctx context.Context, draftID uuid.UUID) (*OrderDraftTable, error)
	GetDraftByCorrelationID(ctx context.Context, correlationID uuid.UUID) (*OrderDraftTable, error)
	MarkDraftPaymentSuccess(ctx context.Context, draftID uuid.UUID) error
	MarkDraftProviderPending(ctx context.Context, draftID uuid.UUID, correlationID uuid.UUID) error
	MarkDraftProviderSuccess(ctx context.Context, draftID uuid.UUID, providerOrderID uuid.UUID) error
	MarkDraftProviderFailed(ctx context.Context, draftID uuid.UUID, reason string) error
	DeleteDraft(ctx context.Context, draftID uuid.UUID) error
	ListPendingDrafts(ctx context.Context, limit int) ([]OrderDraftTable, error)
	UpsertSubscriber(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, connectionID string) error
	DeactivateSubscriber(ctx context.Context, connectionID string, orderID *uuid.UUID) error
	ListActiveSubscribersByOrder(ctx context.Context, orderID uuid.UUID) ([]OrderSubscribersTable, error)
	SaveHistory(ctx context.Context, history *OrderHistoryTable) error
	UpdateHistoryStatus(ctx context.Context, providerOrderID uuid.UUID, status string) error
}

type orderRepository struct {
	draftRepo      *factory.GormRepository[OrderDraftTable]
	subscriberRepo *factory.GormRepository[OrderSubscribersTable]
	historyRepo    *factory.GormRepository[OrderHistoryTable]
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		draftRepo:      factory.NewGormRepository[OrderDraftTable](db),
		subscriberRepo: factory.NewGormRepository[OrderSubscribersTable](db),
		historyRepo:    factory.NewGormRepository[OrderHistoryTable](db),
	}
}

func (r *orderRepository) CreateDraft(ctx context.Context, draft *OrderDraftTable) error {
	_, err := r.draftRepo.Save(ctx, draft)
	return err
}

func (r *orderRepository) GetDraftByID(ctx context.Context, draftID uuid.UUID) (*OrderDraftTable, error) {
	var draft OrderDraftTable
	err := r.draftRepo.DB.WithContext(ctx).Where("id = ?", draftID).First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

func (r *orderRepository) GetDraftByCorrelationID(ctx context.Context, correlationID uuid.UUID) (*OrderDraftTable, error) {
	var draft OrderDraftTable
	err := r.draftRepo.DB.WithContext(ctx).
		Where("provider_correlation_id = ?", correlationID).
		First(&draft).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

func (r *orderRepository) MarkDraftPaymentSuccess(ctx context.Context, draftID uuid.UUID) error {
	return r.draftRepo.DB.WithContext(ctx).
		Model(&OrderDraftTable{}).
		Where("id = ?", draftID).
		Updates(map[string]any{
			"payment_status": true,
		}).Error
}

func (r *orderRepository) MarkDraftProviderPending(ctx context.Context, draftID uuid.UUID, correlationID uuid.UUID) error {
	return r.draftRepo.DB.WithContext(ctx).
		Model(&OrderDraftTable{}).
		Where("id = ?", draftID).
		Updates(map[string]any{
			"provider_pending":        true,
			"provider_correlation_id": correlationID,
			"provider_error":          nil,
		}).Error
}

func (r *orderRepository) MarkDraftProviderSuccess(ctx context.Context, draftID uuid.UUID, providerOrderID uuid.UUID) error {
	return r.draftRepo.DB.WithContext(ctx).
		Model(&OrderDraftTable{}).
		Where("id = ?", draftID).
		Updates(map[string]any{
			"provider_create_status": true,
			"provider_pending":       false,
			"provider_order_id":      providerOrderID,
			"provider_error":         nil,
		}).Error
}

func (r *orderRepository) MarkDraftProviderFailed(ctx context.Context, draftID uuid.UUID, reason string) error {
	return r.draftRepo.DB.WithContext(ctx).
		Model(&OrderDraftTable{}).
		Where("id = ?", draftID).
		Updates(map[string]any{
			"provider_create_status": false,
			"provider_pending":       false,
			"provider_error":         reason,
		}).Error
}

func (r *orderRepository) DeleteDraft(ctx context.Context, draftID uuid.UUID) error {
	return r.draftRepo.DB.WithContext(ctx).
		Where("id = ?", draftID).
		Delete(&OrderDraftTable{}).Error
}

func (r *orderRepository) ListPendingDrafts(ctx context.Context, limit int) ([]OrderDraftTable, error) {
	if limit <= 0 {
		limit = 50
	}
	var drafts []OrderDraftTable
	err := r.draftRepo.DB.WithContext(ctx).
		Where("payment_status = ? AND provider_pending = ? AND provider_create_status = ?", true, true, false).
		Order("updated_at ASC").
		Limit(limit).
		Find(&drafts).Error
	if err != nil {
		return nil, err
	}
	return drafts, nil
}

func (r *orderRepository) UpsertSubscriber(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, connectionID string) error {
	var existing OrderSubscribersTable
	err := r.subscriberRepo.DB.WithContext(ctx).
		Where("order_id = ? AND user_id = ? AND connection_id = ?", orderID, userID, connectionID).
		First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return r.subscriberRepo.DB.WithContext(ctx).Create(&OrderSubscribersTable{
			OrderID:      orderID,
			UserID:       userID,
			ConnectionID: connectionID,
			IsActive:     true,
		}).Error
	}

	return r.subscriberRepo.DB.WithContext(ctx).
		Model(&OrderSubscribersTable{}).
		Where("id = ?", existing.ID).
		Update("is_active", true).Error
}

func (r *orderRepository) DeactivateSubscriber(ctx context.Context, connectionID string, orderID *uuid.UUID) error {
	query := r.subscriberRepo.DB.WithContext(ctx).
		Model(&OrderSubscribersTable{}).
		Where("connection_id = ?", connectionID).
		Where("is_active = ?", true)
	if orderID != nil && *orderID != uuid.Nil {
		query = query.Where("order_id = ?", *orderID)
	}
	return query.Update("is_active", false).Error
}

func (r *orderRepository) ListActiveSubscribersByOrder(ctx context.Context, orderID uuid.UUID) ([]OrderSubscribersTable, error) {
	var subscribers []OrderSubscribersTable
	err := r.subscriberRepo.DB.WithContext(ctx).
		Where("order_id = ? AND is_active = ?", orderID, true).
		Find(&subscribers).Error
	if err != nil {
		return nil, err
	}
	return subscribers, nil
}

func (r *orderRepository) SaveHistory(ctx context.Context, history *OrderHistoryTable) error {
	_, err := r.historyRepo.Save(ctx, history)
	return err
}

func (r *orderRepository) UpdateHistoryStatus(ctx context.Context, providerOrderID uuid.UUID, status string) error {
	return r.historyRepo.DB.WithContext(ctx).
		Model(&OrderHistoryTable{}).
		Where("provider_order_id = ?", providerOrderID).
		Update("status", status).Error
}
