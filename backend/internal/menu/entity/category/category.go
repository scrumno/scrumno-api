package category

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID uint `gorm:"primaryKey"`

	ExternalID string `gorm:"size:128;uniqueIndex;not null"`
	ParentID   *uint  `gorm:"index"`

	Parent    *Category  `gorm:"foreignKey:ParentID"`
	Children []Category `gorm:"foreignKey:ParentID"`

	Name         string `gorm:"size:255;not null"`
	Description  string `gorm:"type:text"`
	ButtonImage  string `gorm:"size:1024"`
	IikoGroupID  string `gorm:"size:128;index"`
	ScheduleID   string `gorm:"size:128;index"`
	ScheduleName string `gorm:"size:255"`
	Sort         int    `gorm:"default:0"`
	IsHidden     bool   `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewBuild() *Category {
	return &Category{}
}
