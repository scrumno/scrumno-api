package app_config

import (
	"context"
	"errors"

	"github.com/google/uuid"
	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type AppConfigRepository interface {
	base.BaseRepository[AppConfig]
	GetWorkingHours(ctx context.Context) (string, string, error)
	GetQueueSyncState(ctx context.Context, venueID uuid.UUID) (int64, error)
	UpdateQueueSyncState(ctx context.Context, venueID uuid.UUID, revision int64) error
}

type appConfigRepository struct {
	*factory.GormRepository[AppConfig]
}

func NewAppConfigRepository(db *gorm.DB) AppConfigRepository {
	return &appConfigRepository{
		GormRepository: factory.NewGormRepository[AppConfig](db),
	}
}

func (r *appConfigRepository) GetWorkingHours(ctx context.Context) (string, string, error) {
	var cfg AppConfig
	db := r.DB.WithContext(ctx)
	venueID, ok := ctx.Value("venue_id").(string)
	if ok && venueID != "" {
		db = db.Where("venue_id = ?", venueID)
	}

	err := db.Select("open_at", "close_at").First(&cfg).Error
	if err != nil {
		return "", "", err
	}
	return cfg.OpenAt, cfg.CloseAt, nil
}

func (r *appConfigRepository) GetQueueSyncState(ctx context.Context, venueID uuid.UUID) (int64, error) {
	var cfg AppConfig
	err := r.DB.WithContext(ctx).
		Where("venue_id = ?", venueID).
		First(&cfg).Error
	if err != nil {
		return 0, err
	}

	return cfg.QueueSyncRevision, nil
}

func (r *appConfigRepository) UpdateQueueSyncState(ctx context.Context, venueID uuid.UUID, revision int64) error {
	var cfg AppConfig
	err := r.DB.WithContext(ctx).
		Where("venue_id = ?", venueID).
		First(&cfg).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cfg = AppConfig{
				VenueID:           venueID,
				QueueSyncRevision: revision,
			}
			return r.DB.WithContext(ctx).Create(&cfg).Error
		}

		return err
	}

	cfg.QueueSyncRevision = revision
	return r.DB.WithContext(ctx).
		Model(&AppConfig{}).
		Where("venue_id = ?", venueID).
		Updates(map[string]any{
			"queue_sync_revision":   cfg.QueueSyncRevision,
			"queue_sync_updated_at": gorm.Expr("NOW()"),
		}).Error
}
