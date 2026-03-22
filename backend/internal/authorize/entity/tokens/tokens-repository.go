package entity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/shared/base"
	"github.com/scrumno/scrumno-api/shared/factory"
	"gorm.io/gorm"
)

type TokensRepository interface {
	base.BaseRepository[AuthorizeToken]
	FindTokenPairBySessionId(ctx context.Context, sessionID string) (*AuthorizeToken, error)
	RevokeTokensByUserSessionId(ctx context.Context, userID uuid.UUID) error
	IsSessionActive(ctx context.Context, sessionID string) (bool, error)
	RevokeTokenBySessionId(ctx context.Context, sessionID string) error
}

type tokensGormRepository struct {
	*factory.GormRepository[AuthorizeToken]
}

func NewTokensRepository(db *gorm.DB) TokensRepository {
	return &tokensGormRepository{
		GormRepository: factory.NewGormRepository[AuthorizeToken](db),
	}
}

func (r *tokensGormRepository) FindTokenPairBySessionId(ctx context.Context, sessionID string) (*AuthorizeToken, error) {
	var t AuthorizeToken
	nowUnix := time.Now().Unix()
	err := r.DB.WithContext(ctx).
		Where("id = ?", sessionID).
		Where("expires_at > ?", nowUnix).
		First(&t).Error

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *tokensGormRepository) RevokeTokensByUserSessionId(ctx context.Context, userID uuid.UUID) error {
	var ac AuthorizeToken

	err := r.DB.WithContext(ctx).
		Model(&ac).
		Where("user_id = ?", userID).
		Where("expires_at > ?", time.Now().Unix()).
		Update("expires_at", time.Now().Unix()).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *tokensGormRepository) RevokeTokenBySessionId(ctx context.Context, sessionID string) error {
	var ac AuthorizeToken

	err := r.DB.WithContext(ctx).
		Model(&ac).
		Where("id = ?", sessionID).
		Update("expires_at", time.Now().Unix()).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *tokensGormRepository) IsSessionActive(ctx context.Context, sessionID string) (bool, error) {
	_, err := r.FindTokenPairBySessionId(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
