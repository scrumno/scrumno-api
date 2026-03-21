package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthorizeToken struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index:idx_authorize_tokens_user" json:"user_id"`
	RefreshToken string    `gorm:"type:text;not null" json:"-"`
	ExpiresAt    int64     `gorm:"not null;index:idx_authorize_tokens_expires" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func NewAuthorizeToken(uuid uuid.UUID, userID uuid.UUID, refreshToken string, expiredDate int64, createdDate time.Time) *AuthorizeToken {
	return &AuthorizeToken{
		ID:           uuid,
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    expiredDate,
		CreatedAt:    createdDate,
	}
}
