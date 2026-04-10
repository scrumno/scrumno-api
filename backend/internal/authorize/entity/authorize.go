package entity

import (
	"time"

	"github.com/google/uuid"
)

type SexType int32

const (
	SexNotSpecified SexType = 0
	SexMale         SexType = 1
	SexFemale       SexType = 2
)

type WalletType int32

const (
	Deposit     WalletType = 0
	Bonus       WalletType = 1
	Products    WalletType = 2
	Discount    WalletType = 3
	Certificate WalletType = 4
)

type WalletBalance struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	IntegrationWalletId *uuid.UUID `gorm:"type:uuid;index" json:"iiko_wallet_id,omitempty"`
	UserID              uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Name                string     `gorm:"type:varchar(255)" json:"name"`
	Type                WalletType `gorm:"type:smallint" json:"type"`
	Balance             float64    `gorm:"type:numeric(12,2);default:0.00" json:"balance"`
}

type User struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Phone         string          `gorm:"uniqueIndex:idx_users_phone;not null"            json:"phone"`
	FullName      *string         `gorm:"null"                                            json:"full_name,omitempty"`
	BirthDate     *time.Time      `gorm:"type:date;null"                                  json:"birth_date,omitempty"`
	IsActive      bool            `gorm:"default:true"                                    json:"is_active"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"                                  json:"created_at"`
	Email         *string         `gorm:"null"                                            json:"email,omitempty"`
	Sex           SexType         `gorm:"type:smallint;default:0" json:"sex"`
	WalletBalance []WalletBalance `gorm:"foreignKey:UserID" json:"walletBalances"`
}

func NewUser(phone string) *User {
	return &User{
		ID:            uuid.New(),
		FullName:      nil,
		BirthDate:     nil,
		IsActive:      true,
		CreatedAt:     time.Now(),
		Email:         nil,
		Phone:         phone,
		Sex:           0,
		WalletBalance: []WalletBalance{},
	}
}
