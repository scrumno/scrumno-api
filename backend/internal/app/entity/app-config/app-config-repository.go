package app_config

import (
	"context"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

type AppConfigRepository interface {
	base.BaseRepository[AppConfig]
	GetWorkingHours(ctx context.Context) (string, string, error)
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
	venueID := ctx.Value("venue_id").(string)
	err := r.DB.WithContext(ctx).
		Where("venue_id = ?", venueID).
		First(&workingHours).Error
	if err != nil {
		return "", "", err
	}
	return workingHours.OpenAt, workingHours.CloseAt, nil
}
