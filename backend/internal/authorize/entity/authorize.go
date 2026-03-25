package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Phone             string     `gorm:"uniqueIndex:idx_users_phone;not null"            json:"phone"`
	FullName          *string    `gorm:"null"                                            json:"full_name,omitempty"`
	BirthDate         *time.Time `gorm:"type:date;null"                                  json:"birth_date,omitempty"`
	IsActive          bool       `gorm:"default:true"                                    json:"is_active"`
	CreatedAt         time.Time  `gorm:"autoCreateTime"                                  json:"created_at"`
	Email             *string    `gorm:"null"                                            json:"email,omitempty"`
}

func NewUser(phone string) *User {
	return &User{
		ID:        uuid.New(),
		FullName:  nil,
		BirthDate: nil,
		IsActive:  true,
		CreatedAt: time.Now(),
		Email:     nil,
		Phone:     phone,
	}
}

