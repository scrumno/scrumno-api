package modifier

import (
	"time"

	"github.com/google/uuid"
)

type CookingTimeModifierTable struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	ModifierID  string    `gorm:"not null;index:idx_cooking_time_modifier_table_modifier_id" json:"modifier_id"`
	CookingTime int       `gorm:"not null;default:0" json:"cooking_time"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()" json:"updated_at"`
}
