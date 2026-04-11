package app_config

import (
	"context"
	"errors"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"github.com/google/uuid"
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
	var workingHours struct {
		OpenAt  string `json:"open_at"`
		CloseAt string `json:"close_at"`
	}
	db := r.DB.WithContext(ctx)
	venueID, ok := ctx.Value("venue_id").(string)
	if ok && venueID != "" {
		db = db.Where("venue_id = ?", venueID)
	}

	err := db.First(&workingHours).Error
	if err != nil {
		return "", "", err
	}
	return workingHours.OpenAt, workingHours.CloseAt, nil
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
