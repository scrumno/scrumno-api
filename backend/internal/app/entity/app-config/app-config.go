package app_config

import (
	"time"

	"github.com/google/uuid"
	themePreset "github.com/scrumno/scrumno-api/internal/app/entity/theme-preset"
)

type AppConfig struct {
	VenueID     uuid.UUID               `gorm:"type:uuid;primaryKey"              json:"venue_id"`
	ThemePreset themePreset.ThemePreset `gorm:"type:varchar(20);default:'light'"  json:"theme_preset"`
	// HEX цвет кнопок и акцентов. Валидировать: /^#[0-9A-Fa-f]{6}$/
	AccentColor string  `gorm:"type:varchar(7);default:'#000000'" json:"accent_color"`
	LogoURL     *string `gorm:"type:text"                         json:"logo_url,omitempty"`
	BannerURL   *string `gorm:"type:text"                         json:"banner_url,omitempty"`
	// Перекрывает venues.address только в UI PWA.
	// venues.address по-прежнему используется для геокодинга.
	AddressManual      *string    `gorm:"type:text" json:"address_manual,omitempty"`
	QueueSyncRevision  int64      `gorm:"not null;default:0" json:"queue_sync_revision"`
	QueueSyncUpdatedAt *time.Time `json:"queue_sync_updated_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`

	WorkingHours struct {
		OpenAt  string `gorm:"type:varchar(5);not null;default:'10:00'" json:"open_at"`
		CloseAt string `gorm:"type:varchar(5);not null;default:'22:00'" json:"close_at"`
	}
}
