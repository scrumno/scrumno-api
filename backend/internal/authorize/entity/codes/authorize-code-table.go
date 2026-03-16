package entity

import (
	"time"

	"github.com/google/uuid"
)

type CodesType string

const (
	RegisterType CodesType = "register"
	AuthType     CodesType = "authorize"
)

type AuthorizeCode struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Phone     string     `gorm:"type:varchar(20);not null;index:idx_authorize_codes_phone" json:"phone"`
	Code      string     `gorm:"type:varchar(10);not null" json:"-"`
	CodeType  CodesType  `gorm:"type:varchar(30);not null" json:"code_type"`
	ExpiresAt time.Time  `gorm:"not null;index:idx_authorize_codes_expires" json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func NewAuthorizeCode(phone string, code string, codeType CodesType) *AuthorizeCode {
	return &AuthorizeCode{
		ID: uuid.New(),
		Phone: phone,
		Code: code,
		CodeType: codeType,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}
}