package customer

import (
	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type WalletBalance struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Type    int       `json:"type"`
	Balance float64   `json:"balance"`
}

type ResponseGet struct {
	ID             uuid.UUID        `json:"id"`
	Name           string           `json:"name"`
	Surname        string           `json:"surname"`
	MiddleName     string           `json:"middleName"`
	Birthday       helpers.IikoTime `json:"birthday"`
	Email          string           `json:"email"`
	Sex            int              `json:"sex"`
	WalletBalances []WalletBalance  `json:"walletBalances"`
	IsDeleted      bool             `json:"isDeleted"`
}
