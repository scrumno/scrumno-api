package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/scrumno/scrumno-api/infrastructure/integration-system/shared/helpers"
)

type CustomerSetResponse struct {
	ID uuid.UUID `json:"id"`
}

type CustomerResponse struct {
	ID                            uuid.UUID         `json:"id"`
	ReferrerID                    *uuid.UUID        `json:"referrerId,omitempty"`
	Name                          string            `json:"name"`
	Surname                       string            `json:"surname"`
	MiddleName                    string            `json:"middleName"`
	Comment                       string            `json:"comment"`
	Phone                         string            `json:"phone"`
	CultureName                   string            `json:"cultureName"`
	Birthday                      helpers.IikoTime  `json:"birthday"`
	Email                         string            `json:"email"`
	Sex                           int               `json:"sex"`
	ConsentStatus                 int               `json:"consentStatus"`
	Anonymized                    bool              `json:"anonymized"`
	Cards                         []Card            `json:"cards"`
	Categories                    []Category        `json:"categories"`
	WalletBalances                []WalletBalance   `json:"walletBalances"`
	UserData                      *string           `json:"userData,omitempty"`
	ShouldReceivePromoActionsInfo bool              `json:"shouldReceivePromoActionsInfo"`
	ShouldReceiveLoyaltyInfo      bool              `json:"shouldReceiveLoyaltyInfo"`
	ShouldReceiveOrderStatusInfo  bool              `json:"shouldReceiveOrderStatusInfo"`
	PersonalDataConsentFrom       *helpers.IikoTime `json:"personalDataConsentFrom,omitempty"`
	PersonalDataConsentTo         *helpers.IikoTime `json:"personalDataConsentTo,omitempty"`
	PersonalDataProcessingFrom    *helpers.IikoTime `json:"personalDataProcessingFrom,omitempty"`
	PersonalDataProcessingTo      *helpers.IikoTime `json:"personalDataProcessingTo,omitempty"`
	IsDeleted                     bool              `json:"isDeleted"`
	WhenRegistered                helpers.IikoTime  `json:"whenRegistered"`
	LastProcessedOrderDate        *helpers.IikoTime `json:"lastProcessedOrderDate,omitempty"`
	FirstOrderDate                *helpers.IikoTime `json:"firstOrderDate,omitempty"`
	LastVisitedOrganizationID     *uuid.UUID        `json:"lastVisitedOrganizationId,omitempty"`
	RegistrationOrganizationID    uuid.UUID         `json:"registrationOrganizationId"`
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	Track       string    `json:"track"`
	Number      string    `json:"number"`
	ValidToDate time.Time `json:"validToDate"`
}

type Category struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	IsActive              bool      `json:"isActive"`
	IsDefaultForNewGuests bool      `json:"isDefaultForNewGuests"`
}

type WalletBalance struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Type    int       `json:"type"`
	Balance float64   `json:"balance"`
}
