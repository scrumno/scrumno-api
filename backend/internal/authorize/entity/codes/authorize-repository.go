package entity

import (
	"context"
	"errors"
	"time"

	factory "github.com/scrumno/scrumno-api/shared/factories/gorm"
	"github.com/scrumno/scrumno-api/shared/interfaces/base"
	"gorm.io/gorm"
)

var (
	ErrRateLimitTooFrequent = errors.New("Повторная отправка возможна через 1 минуту")
	ErrRateLimitHourly      = errors.New("Превышен лимит запросов в час")
)

const (
	minIntervalBetweenCodes = 1 * time.Minute
	maxCodesPerHour         = 5
)

type SmsCodesRepository interface {
	base.BaseRepository[AuthorizeCode]
	ValidateCode(ctx context.Context, phone string, code string, codeType CodesType) (bool, error)
	CountCreatedSince(ctx context.Context, phone string, since time.Time) (int64, error)
	ValidateCodeByCreatedAt(ctx context.Context, phone string) (bool, error)
}

type smsCodesGormRepository struct {
	*factory.GormRepository[AuthorizeCode]
}

func NewSmsCodesRepository(db *gorm.DB) SmsCodesRepository {
	return &smsCodesGormRepository{
		GormRepository: factory.NewGormRepository[AuthorizeCode](db),
	}
}

func (r *smsCodesGormRepository) ValidateCode(
	ctx context.Context,
	phone string,
	code string,
	codeType CodesType,
) (bool, error) {
	var ac AuthorizeCode
	err := r.DB.WithContext(ctx).
		Where("phone = ?", phone).
		Where("code = ?", code).
		Where("code_type = ?", codeType).
		Where("expires_at > ?", time.Now()).
		First(&ac).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	if err := r.DB.WithContext(ctx).Delete(&ac).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *smsCodesGormRepository) CountCreatedSince(ctx context.Context, phone string, since time.Time) (int64, error) {
	var count int64
	sinceUTC := since.UTC()
	err := r.DB.WithContext(ctx).Model(&AuthorizeCode{}).
		Where("phone = ?", phone).
		Where("created_at >= ?", sinceUTC).
		Count(&count).Error
	return count, err
}

func (r *smsCodesGormRepository) ValidateCodeByCreatedAt(ctx context.Context, phone string) (bool, error) {
	now := time.Now().UTC()

	countRecent, err := r.CountCreatedSince(ctx, phone, now.Add(-minIntervalBetweenCodes))
	if err != nil {
		return false, err
	}
	if countRecent >= 1 {
		return false, ErrRateLimitTooFrequent
	}

	countHour, err := r.CountCreatedSince(ctx, phone, now.Add(-time.Hour))
	if err != nil {
		return false, err
	}
	if countHour >= maxCodesPerHour {
		return false, ErrRateLimitHourly
	}

	return true, nil
}
