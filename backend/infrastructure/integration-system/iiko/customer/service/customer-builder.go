package service

import (
	"context"

	"github.com/google/uuid"
	iikoConfig "github.com/scrumno/scrumno-api/infrastructure/integration-system/iiko/config"
	user "github.com/scrumno/scrumno-api/internal/authorize/entity"
)

type SearchableType string

const (
	Phone SearchableType = "phone"
)

type CustomerBodyBuilder struct {
	config *iikoConfig.Config
}

func NewCustomerBodyBuilder(config *iikoConfig.Config) *CustomerBodyBuilder {
	return &CustomerBodyBuilder{
		config: config,
	}
}

type BuilderCommand struct {
	SearchValue    string         `json:"phone"`
	Type           SearchableType `json:"type"`
	OrganizationId uuid.UUID      `json:"organizationId"`
}

type SexType int32

const (
	SexNotSpecified SexType = 0
	SexMale         SexType = 1
	SexFemale       SexType = 2
)

type ConsentStatusType int32

const (
	ConsentUnknown ConsentStatusType = 0
	ConsentGiven   ConsentStatusType = 1
	ConsentRevoked ConsentStatusType = 2
)

type BuilderSetCommand struct {
	ID                            *uuid.UUID         `json:"id,omitempty"`
	Phone                         *string            `json:"phone,omitempty"`
	CardTrack                     *string            `json:"cardTrack,omitempty"`
	CardNumber                    *string            `json:"cardNumber,omitempty"`
	Name                          *string            `json:"name,omitempty"`
	MiddleName                    *string            `json:"middleName,omitempty"`
	SurName                       *string            `json:"surName,omitempty"`
	Birthday                      *string            `json:"birthday,omitempty"`
	Email                         *string            `json:"email,omitempty"`
	Sex                           *SexType           `json:"sex,omitempty"`
	ConsentStatus                 *ConsentStatusType `json:"consentStatus,omitempty"`
	ShouldReceiveLoyaltyInfo      *bool              `json:"shouldReceiveLoyaltyInfo,omitempty"`
	ShouldReceivePromoActionsInfo *bool              `json:"shouldReceivePromoActionsInfo,omitempty"`
	ReferrerId                    *uuid.UUID         `json:"referrerId,omitempty"`
	UserData                      *string            `json:"userData,omitempty"`
	IsDeleted                     *bool              `json:"isDeleted,omitempty"`
	OrganizationId                uuid.UUID          `json:"organizationId"`
}

func (h *CustomerBodyBuilder) BuildGet(ctx context.Context, u *user.User) any {

	if u.Phone != "" {
		return &BuilderCommand{
			SearchValue:    u.Phone,
			Type:           Phone,
			OrganizationId: h.config.OrganizationID,
		}
	}

	return &BuilderCommand{}
}

func (h *CustomerBodyBuilder) BuildSetFromUser(ctx context.Context, u *user.User) any {
	if u == nil {
		return &BuilderSetCommand{}
	}

	name := u.FullName
	isDeleted := false
	loyaltyInfo := false
	promoInfo := false
	sex := SexMale
	consent := ConsentUnknown

	cmd := &BuilderSetCommand{
		OrganizationId:                h.config.OrganizationID,
		Phone:                         &u.Phone,
		Name:                          name,
		Email:                         u.Email,
		Sex:                           &sex,
		ConsentStatus:                 &consent,
		ShouldReceiveLoyaltyInfo:      &loyaltyInfo,
		ShouldReceivePromoActionsInfo: &promoInfo,
		IsDeleted:                     &isDeleted,
	}
	if u.BirthDate != nil {
		birthday := u.BirthDate.Format("2006-01-02")
		cmd.Birthday = &birthday
	}

	return cmd
}
